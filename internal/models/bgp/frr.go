package bgp

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
