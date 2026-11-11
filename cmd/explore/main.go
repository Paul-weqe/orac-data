package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Paul-weqe/orac-data/internal/adapter"

	ifacemodel "github.com/Paul-weqe/orac-data/internal/models/iface"
	ifacenorm "github.com/Paul-weqe/orac-data/internal/normalizer/iface"
)

func main() {
	factory := adapter.NewAdapterFactory()
	a, err := factory.Create(adapter.VendorFRR, &adapter.FRRConfig{
		Host: "localhost",
		Port: 50051,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer a.Close()

	ctx := context.Background()

	// Fetch (typed)
	rawIface, err := a.FetchInterfaces(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Type-assert to what FetchInterfaces returns
	resp, ok := rawIface.(*ifacemodel.FRRInterfaceResponse)
	if !ok {
		log.Fatalf("unexpected type: %T", rawIface)
	}

	// Normalize
	norm := ifacenorm.NewFRRInterfaceNormalizer()
	data, err := norm.Normalize(resp)
	if err != nil {
		log.Fatal(err)
	}

	// Print unified output
	pretty, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(pretty))
}
