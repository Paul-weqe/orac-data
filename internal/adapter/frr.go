package adapter

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type FRRAdapter struct {
	// Connection details.
	host string
	port int

	// gRPC connection
	conn   *grpc.ClientConn
	client interface{}

	// State
	isConnected bool
}

func NewFRRAdapter(host string, port int) *FRRAdapter {
	return &FRRAdapter{
		host:        host,
		port:        port,
		isConnected: false,
	}
}

// Establishes the gRPC connection to FRR.
func (a *FRRAdapter) connect(ctx context.Context) error {
	if a.isConnected {
		return nil
	}

	// Build the target address.
	target := fmt.Sprintf("%s:%d", a.host, a.port)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(target, opts...)
	if err != nil {
		return fmt.Errorf("failed to connect to FRR at %s: %w", target, err)
	}

	a.conn = conn
	a.isConnected = true

	// TODO: initiate FRR client once we have the proto definitions
	return nil
}

// FetchRoutes retrieves all BGP routers from FRR.
func (a *FRRAdapter) FetchRoutes(ctx context.Context) (interface{}, error) {
	if err := a.connect(ctx); err != nil {
		return nil, err
	}

	// TODO: Implement actual gRPC call to FRR
	// For now, return a placeholder error
	return nil, fmt.Errorf("FetchRoutes not yet implemented - waiting for FRR proto definitions")
}

func (a *FRRAdapter) FetchNeighbors(ctx context.Context) (interface{}, error) {
	if err := a.connect(ctx); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("FetchNeighbors not yet implemented - waiting for FRR proto definitions")
}

func (a *FRRAdapter) Vendor() string {
	return "frr"
}

func (a *FRRAdapter) Close() error {
	if a.conn != nil {
		err := a.conn.Close()
		a.isConnected = false
		return err
	}
	return nil
}
