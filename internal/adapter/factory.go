package adapter

import "fmt"

type Vendor string

const (
	VendorFRR Vendor = "frr"

	// Future vensors:
)

type FRRConfig struct {
	Host string `json:"host"`
	Port int    `json:"port,omitempty"`
}

type AdapterFactory struct{}

func NewAdapterFactory() *AdapterFactory {
	return &AdapterFactory{}
}

func (f *AdapterFactory) Create(vendor Vendor, config interface{}) (NetworkAdapter, error) {
	switch vendor {
	case VendorFRR:
		return f.createFRRAdapter(config)

	default:
		return nil, fmt.Errorf("unsupported vendor: %s", vendor)
	}
}

func (f *AdapterFactory) createFRRAdapter(config interface{}) (NetworkAdapter, error) {
	cfg, ok := config.(*FRRConfig)

	if !ok {
		return nil, fmt.Errorf("invalid config type for FRR: expected *FRRConfig got %T", config)
	}

	if cfg.Host == "" {
		return nil, fmt.Errorf("missing 'Host' in FRR config")
	}

	if cfg.Port == 0 {
		cfg.Port = 50051
	}

	return NewFRRAdapter(cfg.Host, cfg.Port), nil
}
