package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"orac-data/internal/models/iface"
	pb "orac-data/proto/frr"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var _ NetworkAdapter = (*FRRAdapter)(nil)

type FRRAdapter struct {
	// Connection details.
	host string
	port int

	// gRPC connection
	conn   *grpc.ClientConn
	client pb.NorthboundClient

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
	a.client = pb.NewNorthboundClient(conn)
	a.isConnected = true

	// TODO: initiate FRR client once we have the proto definitions
	return nil
}

// Retrieves all the interface information.
func (a *FRRAdapter) FetchInterfaces(ctx context.Context) (interface{}, error) {
	// 1. Ensure Connected.
	if err := a.connect(ctx); err != nil {
		return nil, err
	}

	// 2. Build requests matching the known-working grpcurl command.
	req := &pb.GetRequest{
		Type:         pb.GetRequest_ALL,
		Encoding:     pb.Encoding_JSON,
		WithDefaults: false,
		Path:         []string{"/frr-interface:lib"},
	}

	// 3. Call Get() Fetching server stream.
	stream, err := a.client.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("FRR Get() call failed %w", err)
	}

	// 4. Drain the stream, concatenating all JSON chunks.
	var buf strings.Builder
	chunkCount := 0
	for {
		reply, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("stream recv failed after %d chunks: %w", chunkCount, err)
		}
		buf.WriteString(reply.GetData().GetData())
		chunkCount++
	}

	if chunkCount == 0 {
		return nil, fmt.Errorf("FRR returned no interface data")
	}

	var resp iface.FRRInterfaceResponse
	if err := json.Unmarshal([]byte(buf.String()), &resp); err != nil {
		return nil, fmt.Errorf("failed to parse FRR interface JSON: %w", err)
	}

	return &resp, nil
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
