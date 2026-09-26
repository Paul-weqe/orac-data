package iface

type FRRResponse struct {
	Lib FRRInterfaceLib `json:"frr-interface:lib"`
}

type FRRInterfaceLib struct {
	Interface []FRRInterface `json:"interface"`
}

type FRRInterface struct {
	Name  string            `json:"name"`
	VRF   string            `json:"vrf"`
	State FRRInterfaceState `json:"state"`
}

type FRRInterfaceState struct {
	IfIndex    int    `json:"if-index"`
	MTU        int    `json:"mtu"`
	MTU6       int    `json:"mtu6"`
	Speed      int    `json:"speed"`
	Metric     int    `json:"metric"`
	PhyAddress string `json:"phy-address"`
}
