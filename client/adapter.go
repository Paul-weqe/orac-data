package client

import (
	"context"

	"github.com/Paul-weqe/orac-data/models"
)

type Vendor string

const (
	VendorFRR Vendor = "frr"

	// Future Vendors:
)

type NetworkAdapter interface {
	FetchInterfaces(ctx context.Context) (*models.InterfaceData, error)
	FetchRoutes(ctx context.Context) (*models.BGPData, error)
	FetchNeighbors(ctx context.Context) (interface{}, error)
	Vendor() string
	Close() error
}
