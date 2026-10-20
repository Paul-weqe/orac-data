package adapter

import (
	"context"
)

type NetworkAdapter interface {
	FetchInterfaces(ctx context.Context) (interface{}, error)
	FetchRoutes(ctx context.Context) (interface{}, error)
	FetchNeighbors(ctx context.Context) (interface{}, error)
	Vendor() string
	Close() error
}
