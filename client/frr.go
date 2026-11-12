package client

import (
	"context"
	"fmt"

	internaladapter "github.com/Paul-weqe/orac-data/internal/adapter"
	ifacemodel "github.com/Paul-weqe/orac-data/internal/models/iface"
	ifacenorm "github.com/Paul-weqe/orac-data/internal/normalizer/iface"
	"github.com/Paul-weqe/orac-data/models"
)

// FRRConfig holds the connection parameters for an FRR device.
type FRRConfig struct {
	Host string `json:"host"`
	Port int    `json:"port,omitempty"`
}

// frrClient wraps the internal FRR adapter and implements NetworkAdapter
// by pairing each raw fetch with the corresponding normalizer.
type frrClient struct {
	raw *internaladapter.FRRAdapter
}

// connectFRR validates the config, dials the device, and returns a
// NetworkAdapter that yields normalized data.
func connectFRR(ctx context.Context, cfg *FRRConfig) (NetworkAdapter, error) {
	if cfg.Host == "" {
		return nil, fmt.Errorf("missing 'Host' in FRR config")
	}
	port := cfg.Port
	if port == 0 {
		port = 50051
	}

	raw := internaladapter.NewFRRAdapter(cfg.Host, port)
	if err := raw.Connect(ctx); err != nil {
		return nil, err
	}

	return &frrClient{raw: raw}, nil
}

func (c *frrClient) FetchInterfaces(ctx context.Context) (*models.InterfaceData, error) {
	raw, err := c.raw.FetchInterfaces(ctx)
	if err != nil {
		return nil, err
	}

	resp, ok := raw.(*ifacemodel.FRRInterfaceResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected FRR interface response type: %T", raw)
	}

	norm := ifacenorm.NewFRRInterfaceNormalizer()
	return norm.Normalize(resp)
}

func (c *frrClient) FetchRoutes(ctx context.Context) (*models.BGPData, error) {
	return nil, fmt.Errorf("FetchRoutes not yet implemented")
}

func (c *frrClient) FetchNeighbors(ctx context.Context) (interface{}, error) {
	return nil, fmt.Errorf("FetchNeighbors not yet implemented")
}

func (c *frrClient) Vendor() string {
	return string(VendorFRR)
}

func (c *frrClient) Close() error {
	return c.raw.Close()
}
