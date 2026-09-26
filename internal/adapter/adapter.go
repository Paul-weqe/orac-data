package adapter

import (
	"context"
)

type BGPAdapter interface {
	FetchRoutes(ctx context.Context) (interface{}, error)
	FetchNeighbors(ctx context.Context) (interface{}, error)
	Vendor() string
	Close() error
}
