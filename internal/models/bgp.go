package models

type FRRRoutesResponse struct {
	Routes      map[string][]FRRPath `json:"routes"`
	TotalRoutes int                  `json:"totalRoutes"`
	TotalPaths  int                  `json:"totalPaths"`
}

type FRRPath struct {
	Valid           bool   `json:"valid"`
	Bestpath        bool   `json:"bestpath"`
	SelectionReason string `json:"selectionReason,omitempty"`
	PathFrom        string `json:"pathFrom,omitempty"`
	Prefix          string `json:"prefix"`
	PrefixLen       int    `json:"prefixLen"`
	Network         string `json:"network"`
	Version         int    `json:"version"`
	Metric          int    `json:"metric"`
	Weight          int    `json:"weight"`
	PeerID          string `json:"peerId"`
	Path            string `json:"path"`
	Origin          string `json:"origin"`
	Nexthops        []struct {
		IP       string `json:"ip"`
		Hostname string `json:"hostname,omitempty"`
		AFI      string `json:"afi"`
		Used     bool   `json:"used"`
	} `json:"nexthops"`
}

// What our Normalizer should produce.
// Regardless of the vendor being used.
type UnifiedRoute struct {
	Network         string                 `json:"network"`
	PrefixLength    int                    `json:"prefixLength"`
	CIDR            string                 `json:"cidr"`
	NextHop         string                 `json:"nextHop"`
	NextHopHostname string                 `json:"nextHopHostname"`
	PeerID          string                 `json:"peerId"`
	Metric          int                    `json:"metric"`
	Weight          int                    `json:"weight"`
	ASPath          []int                  `json:"asPath"`
	ASPathStr       string                 `json:"asPathStr"`
	ASPathLen       int                    `json:"asPathLen"`
	Origin          string                 `json:"origin"`
	Best            bool                   `json:"best"`
	Valid           bool                   `json:"valid"`
	SelectionReason string                 `json:"selectionReason"`
	AddressFamily   string                 `json:"addressFamily"`
	VendorMetadata  map[string]interface{} `json:"vendorMetadata,omitempty"`
}

// Represents all BGP data from a device.
type BGPData struct {
	Vendor         string                    `json:"vendor"`
	Hostname       string                    `json:"hostname"`
	TotalRoutes    int                       `json:"totalRoutes"`
	TotalPaths     int                       `json:"totalPaths"`
	Routes         []UnifiedRoute            `json:"Routes"`
	RoutesByPrefix map[string][]UnifiedRoute `json:"routesByPrefix"`
	BestRoutes     map[string]UnifiedRoute   `json:"bestRoutes"`
	Stats          BGPStats                  `json:"stats"`
}

// BGPStats holds statistical information about BGP data
type BGPStats struct {
	// Basic counts
	UniquePrefixes int `json:"uniquePrefixes"`
	TotalPaths     int `json:"totalPaths"`

	// AS Path statistics
	AverageASPathLen float64 `json:"averageASPathLen"`
	MaxASPathLen     int     `json:"maxASPathLen"`

	// Distribution
	PrefixDistribution map[int]int `json:"prefixDistribution"`

	UniqueASNs []int `json:"uniqueASNs"`
}
