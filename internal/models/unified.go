package models

import (
	"orac-data/internal/models/bgp"
	"orac-data/internal/models/iface"
)

type (
	// bgp
	FRRBGPPath           = bgp.FRRPath
	FRRBGPRoutesResponse = bgp.FRRRoutesResponse

	// interface
	FRRInterfaceResponse = iface.FRRResponse
	FRRInterface         = iface.FRRInterface
	FRRInterfaceLib      = iface.FRRInterfaceLib
	FRRInterfaceState    = iface.FRRInterfaceState
)

// ===== INTERFACE ======

type InterfaceData struct {
	Vendor           string                      `json:"vendor"`
	Interfaces       []UnifiedInterface          `json:"interfaces"`
	InterfacesByName map[string]UnifiedInterface `json:"interfacesByName"`
	Stats            InterfaceStats              `json:"stats"`
}

type UnifiedInterface struct {
	Name       string `json:"name"`
	VRF        string `json:"vrf"`
	Index      int    `json:"index"`
	MTU        int    `json:"mtu"`
	MTU6       int    `json:"mtu6"`
	Speed      int    `json:"speed"`
	Metric     int    `json:"metric"`
	MACAddress string `json:"macAddress"`

	// Any extra info from the vendor
	VendorMetadata map[string]interface{} `json:"vendorMetadata,omitempty"`
}

type InterfaceStats struct {
	Total int            `json:"total"`
	ByVRF map[string]int `json:"byVRF"`
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

// ==== BGP ====

type BGPData struct {
	Vendor         string                    `json:"vendor"`
	Hostname       string                    `json:"hostname"`
	TotalRoutes    int                       `json:"totalRoutes"`
	TotalPaths     int                       `json:"totalPaths"`
	Routes         []UnifiedRoute            `json:"routes"`
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
