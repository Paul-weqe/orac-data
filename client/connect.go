package client

import (
	"context"
	"fmt"
)

func Connect(ctx context.Context, cfg interface{}) (NetworkAdapter, error) {
	switch c := cfg.(type) {

	case *FRRConfig:
		return connectFRR(ctx, c)
	default:
		return nil, fmt.Errorf("unsupported config type: %T", cfg)
	}
}
