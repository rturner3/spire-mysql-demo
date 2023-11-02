package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/rturner3/spire-mysql-demo/pkg/common"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// SPIRE Agent socket path
	socketPath = "unix:///run/spire/sockets/agent.sock"

	// MySQL related constants
	mysqlUser           = "mysql-tls-reloader"
	mysqlClientSVIDHint = "mysql-client"
	reloadTLSQuery      = "ALTER INSTANCE RELOAD TLS"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// Wait for an os.Interrupt signal
	go waitForCtrlC(cancel)

	// Start X.509 watcher
	startWatcher(ctx)
}

func startWatcher(ctx context.Context) {
	// Creates a new Workload API client, connecting to provided socket path
	// Environment variable `SPIFFE_ENDPOINT_SOCKET` is used as default
	client, err := workloadapi.New(ctx, workloadapi.WithAddr(socketPath))
	if err != nil {
		log.Fatalf("Unable to create workload API client: %v", err)
	}
	defer client.Close()

	// Start a watcher for X.509 SVID updates
	doneCh := make(chan struct{}, 1)
	go func() {
		err := client.WatchX509Context(ctx, &x509Watcher{})
		if err != nil && status.Code(err) != codes.Canceled {
			log.Fatalf("Error watching X.509 context: %v", err)
		}
		doneCh <- struct{}{}
	}()

	<-doneCh
}

// x509Watcher is a sample implementation of the workloadapi.X509ContextWatcher interface
type x509Watcher struct{}

// OnX509ContextUpdate is run every time an SVID is updated
func (w *x509Watcher) OnX509ContextUpdate(c *workloadapi.X509Context) {
	if err := common.LogSVIDs(c); err != nil {
		return
	}

	if err := common.WriteMySQLServerSVIDFiles(c); err != nil {
		log.Printf("Failed to write SVID/Bundle to disk: %v", err)
		return
	}

	log.Printf("Successfully written SVID/Bundle to disk")

	db, err := common.NewMySQLDBWithSPIRETLSConfig(c, mysqlUser, "", mysqlClientSVIDHint)
	if err != nil {
		log.Printf("Failed to create MySQL Client: %v", err)
		return
	}

	_, err = db.ExecContext(context.Background(), reloadTLSQuery)
	if err != nil {
		log.Printf("Failed to run reload TLS query: %v", err)
		return
	}

	log.Printf("Successfully reloaded MySQL TLS config")
}

// OnX509ContextWatchError is run when the client runs into an error
func (w *x509Watcher) OnX509ContextWatchError(err error) {
	if status.Code(err) != codes.Canceled {
		log.Printf("OnX509ContextWatchError error: %v", err)
	}
}

// waitForCtrlC waits until an os.Interrupt signal is sent (ctrl + c)
func waitForCtrlC(cancel context.CancelFunc) {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	<-signalCh

	cancel()
}
