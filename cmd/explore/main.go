package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"orac-data/internal/adapter"

	ifacemodel "orac-data/internal/models/iface"
	ifacenorm "orac-data/internal/normalizer/iface"
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

/// Capabilities:
///
/// {
///   "frrVersion": "10.0-MyGrpcBuild-gUNKNOWN",
///   "supportedModules": [
///     {
///       "name": "frr-bgp-route-map",
///       "organization": "Free Range Routing",
///       "revision": "2020-01-02"
///     },
///     {
///       "name": "frr-filter",
///       "organization": "FRRouting",
///       "revision": "2019-07-04"
///     },
///     {
///       "name": "frr-interface",
///       "organization": "FRRouting",
///       "revision": "2020-02-05"
///     },
///     {
///       "name": "frr-route-map",
///       "organization": "FRRouting",
///       "revision": "2019-07-01"
///     },
///     {
///       "name": "frr-vrf",
///       "organization": "Free Range Routing",
///       "revision": "2019-12-06"
///     }
///   ],
///   "supportedEncodings": [
///     "JSON",
///     "XML"
///   ]
/// }
