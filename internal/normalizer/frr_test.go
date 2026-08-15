package normalizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const sampleFRRJSON = `{
  "routes": {
    "1.1.1.1/32": [{
      "valid": true,
      "bestpath": true,
      "selectionReason": "First path received",
      "pathFrom": "external",
      "prefix": "1.1.1.1",
      "prefixLen": 32,
      "network": "1.1.1.1/32",
      "version": 1,
      "metric": 0,
      "weight": 32768,
      "peerId": "(unspec)",
      "path": "",
      "origin": "IGP",
      "nexthops": [{
        "ip": "0.0.0.0",
        "hostname": "rt1",
        "afi": "ipv4",
        "used": true
      }]
    }],
    "2.2.2.2/32": [{
      "valid": true,
      "bestpath": true,
      "selectionReason": "AS Path",
      "pathFrom": "external",
      "prefix": "2.2.2.2",
      "prefixLen": 32,
      "network": "2.2.2.2/32",
      "version": 2,
      "metric": 0,
      "weight": 0,
      "peerId": "10.0.1.2",
      "path": "65101",
      "origin": "IGP",
      "nexthops": [{
        "ip": "10.0.1.2",
        "hostname": "rt2",
        "afi": "ipv4",
        "used": true
      }]
    }, {
      "valid": true,
      "bestpath": false,
      "selectionReason": "",
      "pathFrom": "external",
      "prefix": "2.2.2.2",
      "prefixLen": 32,
      "network": "2.2.2.2/32",
      "version": 2,
      "metric": 0,
      "weight": 0,
      "peerId": "10.0.2.3",
      "path": "65000 65101",
      "origin": "IGP",
      "nexthops": [{
        "ip": "10.0.2.3",
        "hostname": "rt3",
        "afi": "ipv4",
        "used": true
      }]
    }]
  },
  "totalRoutes": 2,
  "totalPaths": 3
}`

func TestFRRNormalizer(t *testing.T) {
	// Create normalizer
	normalizer := NewFRRNormalizer()

	// Normalize the sample data
	data, err := normalizer.Normalize([]byte(sampleFRRJSON))
	require.NoError(t, err)
	require.NotNil(t, data)

	// Check basic fields
	assert.Equal(t, "frr", data.Vendor)
	assert.Equal(t, 2, data.TotalRoutes)
	assert.Equal(t, 3, data.TotalPaths)

	// Check routes
	assert.Len(t, data.Routes, 3)
	assert.Len(t, data.RoutesByPrefix, 2)
	assert.Len(t, data.BestRoutes, 2)

	// Check specific prefix
	paths, ok := data.RoutesByPrefix["2.2.2.2/32"]
	assert.True(t, ok)
	assert.Len(t, paths, 2)

	// Check best route for 2.2.2.2/32
	best, ok := data.BestRoutes["2.2.2.2/32"]
	assert.True(t, ok)
	assert.Equal(t, "10.0.1.2", best.NextHop)
	assert.Equal(t, []int{65101}, best.ASPath)
	assert.True(t, best.Best)

	// Check non-best route
	nonBest := paths[1]
	assert.Equal(t, "10.0.2.3", nonBest.NextHop)
	assert.Equal(t, []int{65000, 65101}, nonBest.ASPath)
	assert.False(t, nonBest.Best)

	// Check statistics
	stats := data.Stats
	assert.Equal(t, 2, stats.UniquePrefixes)
	assert.Equal(t, 3, stats.TotalPaths)
	assert.Equal(t, 1.0, stats.AverageASPathLen) // (0 + 1 + 2) / 3 = 1.0
	assert.Equal(t, 2, stats.MaxASPathLen)
	assert.Contains(t, stats.UniqueASNs, 65000)
	assert.Contains(t, stats.UniqueASNs, 65101)
	assert.Equal(t, 2, stats.PrefixDistribution[32]) // Both routes are /32
}

func TestFRRNormalizer_EmptyData(t *testing.T) {
	normalizer := NewFRRNormalizer()

	emptyJSON := `{"routes": {}, "totalRoutes": 0, "totalPaths": 0}`
	data, err := normalizer.Normalize([]byte(emptyJSON))
	require.NoError(t, err)

	assert.Equal(t, 0, data.TotalRoutes)
	assert.Equal(t, 0, data.TotalPaths)
	assert.Empty(t, data.Routes)
	assert.Empty(t, data.RoutesByPrefix)
	assert.Empty(t, data.BestRoutes)

	// Check stats are zero
	stats := data.Stats
	assert.Equal(t, 0, stats.UniquePrefixes)
	assert.Equal(t, 0, stats.TotalPaths)
	assert.Equal(t, 0.0, stats.AverageASPathLen)
	assert.Equal(t, 0, stats.MaxASPathLen)
	assert.Empty(t, stats.UniqueASNs)
	assert.Empty(t, stats.PrefixDistribution)
}

func TestFRRNormalizer_InvalidJSON(t *testing.T) {
	normalizer := NewFRRNormalizer()

	invalidJSON := `{invalid json}`
	data, err := normalizer.Normalize([]byte(invalidJSON))
	assert.Error(t, err)
	assert.Nil(t, data)
}

func TestParseASPath(t *testing.T) {
	normalizer := NewFRRNormalizer()

	tests := []struct {
		name     string
		input    string
		expected []int
		str      string
		length   int
	}{
		{
			name:     "empty path",
			input:    "",
			expected: []int{},
			str:      "",
			length:   0,
		},
		{
			name:     "single ASN",
			input:    "65101",
			expected: []int{65101},
			str:      "65101",
			length:   1,
		},
		{
			name:     "multiple ASNs",
			input:    "65000 65101 65200",
			expected: []int{65000, 65101, 65200},
			str:      "65000 65101 65200",
			length:   3,
		},
		{
			name:     "with extra spaces",
			input:    "  65000   65101  ",
			expected: []int{65000, 65101},
			str:      "  65000   65101  ",
			length:   2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asPath, str, length := normalizer.parseASPath(tt.input)
			assert.Equal(t, tt.expected, asPath)
			assert.Equal(t, tt.str, str)
			assert.Equal(t, tt.length, length)
		})
	}
}

func TestVendorMetadata(t *testing.T) {
	normalizer := NewFRRNormalizer()

	// Use the sample data
	data, err := normalizer.Normalize([]byte(sampleFRRJSON))
	require.NoError(t, err)

	// Check metadata for first route (1.1.1.1/32)
	route := data.Routes[0]
	assert.NotNil(t, route.VendorMetadata)
	assert.Equal(t, "First path received", route.VendorMetadata["selectionReason"])
	assert.Equal(t, "external", route.VendorMetadata["pathFrom"])
	assert.Equal(t, float64(1), route.VendorMetadata["version"])
}
