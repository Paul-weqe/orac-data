package models

import (
	"encoding/json"
	"testing"
)

// Sample FRR JSON data
const frrSample = `{
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

func TestParseFRRJSON(t *testing.T) {
	var response FRRRoutesResponse

	err := json.Unmarshal([]byte(frrSample), &response)
	if err != nil {
		t.Fatalf("Failed to parse FRR JSON: %v", err)
	}

	// Check total counts
	if response.TotalRoutes != 2 {
		t.Errorf("Expected 2 routes, got %d", response.TotalRoutes)
	}
	if response.TotalPaths != 3 {
		t.Errorf("Expected 3 paths, got %d", response.TotalPaths)
	}

	// Check routes count
	if len(response.Routes) != 2 {
		t.Errorf("Expected 2 route entries, got %d", len(response.Routes))
	}

	// Check specific prefix
	paths, ok := response.Routes["2.2.2.2/32"]
	if !ok {
		t.Fatal("Route 2.2.2.2/32 not found")
	}

	if len(paths) != 2 {
		t.Errorf("Expected 2 paths for 2.2.2.2/32, got %d", len(paths))
	}

	// Check best path
	if !paths[0].Bestpath {
		t.Error("First path should be bestpath")
	}
	if paths[0].PeerID != "10.0.1.2" {
		t.Errorf("Expected peer 10.0.1.2, got %s", paths[0].PeerID)
	}

	// Check second path (non-best)
	if paths[1].Bestpath {
		t.Error("Second path should NOT be bestpath")
	}
	if paths[1].Path != "65000 65101" {
		t.Errorf("Expected path '65000 65101', got '%s'", paths[1].Path)
	}
}

func TestUnifiedRouteStructure(t *testing.T) {
	route := UnifiedRoute{
		Network:       "2.2.2.2",
		PrefixLength:  32,
		CIDR:          "2.2.2.2/32",
		NextHop:       "10.0.1.2",
		PeerID:        "10.0.1.2",
		ASPath:        []int{65101},
		ASPathStr:     "65101",
		ASPathLen:     1,
		Best:          true,
		Valid:         true,
		Origin:        "IGP",
		AddressFamily: "ipv4",
		VendorMetadata: map[string]interface{}{
			"selectionReason": "AS Path",
			"version":         2,
		},
	}

	// Test JSON marshaling
	data, err := json.Marshal(route)
	if err != nil {
		t.Fatalf("Failed to marshal unified route: %v", err)
	}

	// Test JSON unmarshaling
	var route2 UnifiedRoute
	err = json.Unmarshal(data, &route2)
	if err != nil {
		t.Fatalf("Failed to unmarshal unified route: %v", err)
	}

	if route2.Network != route.Network {
		t.Errorf("Expected network %s, got %s", route.Network, route2.Network)
	}
	if route2.ASPathLen != route.ASPathLen {
		t.Errorf("Expected ASPathLen %d, got %d", route.ASPathLen, route2.ASPathLen)
	}
}
