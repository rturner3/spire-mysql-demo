package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/spiffe/go-spiffe/v2/workloadapi"
)

const (
	socketPath = "unix:///run/spire/sockets/agent.sock"
	svidDir    = "/spire/certs"
	bundleFile = "bundle.0.pem"
	certFile   = "svid.0.pem"
	keyFile    = "svid.0.key"
)

var (
	certFilePath   = fmt.Sprintf("%s/%s", svidDir, certFile)
	keyFilePath    = fmt.Sprintf("%s/%s", svidDir, keyFile)
	bundleFilePath = fmt.Sprintf("%s/%s", svidDir, bundleFile)
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

	err = writeX509Context(x509Context)
	if err != nil {
		log.Fatalf("failed to write SVID/Bundle to disk: %v", err)
	}

	log.Printf("SVID/Bundle files written successfully in %s directory", svidDir)
}

func writeX509Context(c *workloadapi.X509Context) error {
	certBytes, keyBytes, err := c.SVIDs[0].Marshal()
	if err != nil {
		return err
	}

	bundleBytes, err := c.Bundles.Bundles()[0].Marshal()
	if err != nil {
		return err
	}

	err = os.WriteFile(certFilePath, certBytes, 0644)
	if err != nil {
		return err
	}

	err = os.WriteFile(keyFilePath, keyBytes, 0644)
	if err != nil {
		return err
	}

	err = os.WriteFile(bundleFilePath, bundleBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
