package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	internaladapter "github.com/Paul-weqe/orac-data/internal/adapter"
)

func main() {
	ctx := context.Background()

	a := internaladapter.NewFRRAdapter("localhost", 50051)
	defer a.Close()

	if err := a.Connect(ctx); err != nil {
		log.Fatalf("connect: %v", err)
	}

	raw, err := a.FetchInterfaces(ctx)
	if err != nil {
		log.Fatalf("fetch: %v", err)
	}

	pretty, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		log.Fatalf("marshal: %v", err)
	}
	fmt.Println(string(pretty))
}
