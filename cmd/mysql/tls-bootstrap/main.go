package main

import (
	"context"
	"log"

	"github.com/rturner3/spire-mysql-demo/pkg/common"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
)

const (
	// SPIRE Agent socket path
	socketPath = "unix:///run/spire/sockets/agent.sock"
)

func main() {
	ctx := context.Background()
	// Creates a new Workload API client, connecting to provided socket path
	// Environment variable `SPIFFE_ENDPOINT_SOCKET` is used as default
	client, err := workloadapi.New(ctx, workloadapi.WithAddr(socketPath))
	if err != nil {
		log.Fatalf("Unable to create workload API client: %v", err)
	}
	defer client.Close()

	x509Context, err := client.FetchX509Context(ctx)
	if err != nil {
		log.Fatalf("Unable to fetch x509Context %v", err)
	}

	if err := common.WriteMySQLServerSVIDFiles(x509Context); err != nil {
		log.Printf("Failed to write SVID/Bundle to disk: %v", err)
		return
	}

	log.Printf("SVID/Bundle files written successfully")
}
