package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/go-sql-driver/mysql"
	"github.com/spiffe/go-spiffe/v2/svid/x509svid"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// SPIRE related constants
	socketPath = "unix:///run/spire/sockets/agent.sock"
	svidDir    = "/spire/certs"
	bundleFile = "bundle.0.pem"
	certFile   = "svid.0.pem"
	keyFile    = "svid.0.key"

	// MySQL related constants
	mysqlUser           = "mysql-tls-reloader"
	mysqlHost           = "mysql.mysql.svc.cluster.local"
	mysqlPort           = "3306"
	mysqlTLSConfigName  = "spire-ssl"
	mysqlServerSVIDHint = "mysql-server"
	mysqlClientSVIDHint = "mysql-client"
	reloadTLSQuery      = "ALTER INSTANCE RELOAD TLS"
)

var (
	certFilePath   = fmt.Sprintf("%s/%s", svidDir, certFile)
	keyFilePath    = fmt.Sprintf("%s/%s", svidDir, keyFile)
	bundleFilePath = fmt.Sprintf("%s/%s", svidDir, bundleFile)

	// Format is specified https://github.com/go-sql-driver/mysql#dsn-data-source-name
	dbConnectionString = fmt.Sprintf("%s@tcp(%s:%s)/?tls=%s", mysqlUser, mysqlHost, mysqlPort, mysqlTLSConfigName)
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// Wait for an os.Interrupt signal
	go waitForCtrlC(cancel)

	// Start X.509 watcher
	startWatcher(ctx)
}

func startWatcher(ctx context.Context) {
	var wg sync.WaitGroup

	// Creates a new Workload API client, connecting to provided socket path
	// Environment variable `SPIFFE_ENDPOINT_SOCKET` is used as default
	client, err := workloadapi.New(ctx, workloadapi.WithAddr(socketPath))
	if err != nil {
		log.Fatalf("Unable to create workload API client: %v", err)
	}
	defer client.Close()

	wg.Add(1)
	// Start a watcher for X.509 SVID updates
	go func() {
		defer wg.Done()
		err := client.WatchX509Context(ctx, &x509Watcher{})
		if err != nil && status.Code(err) != codes.Canceled {
			log.Fatalf("Error watching X.509 context: %v", err)
		}
	}()

	wg.Wait()
}

// x509Watcher is a sample implementation of the workloadapi.X509ContextWatcher interface
type x509Watcher struct{}

// OnX509ContextUpdate is run every time an SVID is updated
func (w *x509Watcher) OnX509ContextUpdate(c *workloadapi.X509Context) {
	// write SVID to disk
	err := writeX509Context(c)
	if err != nil {
		log.Printf("Failed to write SVID/Bundle to disk: %v", err)
		return
	}

	log.Printf("Successfully written SVID/Bundle to disk")

	// Create TLS config with client certificates
	tlsConf, err := createTLSConf(c)
	if err != nil {
		log.Printf("Failed to create MySQL TLS config: %v", err)
		return
	}

	err = mysql.RegisterTLSConfig("spire-ssl", tlsConf)
	if err != nil {
		log.Printf("Failed to register MySQL TLS config: %v", err)
		return
	}

	db, err := sql.Open("mysql", dbConnectionString)
	if err != nil {
		log.Printf("Failed to open MySQL database: %v", err)
		return
	}
	defer db.Close()

	_, err = db.ExecContext(context.Background(), reloadTLSQuery)
	if err != nil {
		log.Printf("Failed to run reload TLS query: %v", err)
		return
	}

	log.Printf("Successfully reloaded MySQL TLS config")
}

func writeX509Context(c *workloadapi.X509Context) error {
	svid, err := getSVIDByHint(c, mysqlServerSVIDHint)
	if err != nil {
		return err
	}

	certBytes, keyBytes, err := svid.Marshal()
	if err != nil {
		return err
	}

	bundleBytes, err := c.Bundles.Bundles()[0].Marshal()
	if err != nil {
		return err
	}

	err = os.WriteFile(certFilePath, certBytes, 0o644)
	if err != nil {
		return err
	}

	err = os.WriteFile(keyFilePath, keyBytes, 0o644)
	if err != nil {
		return err
	}

	err = os.WriteFile(bundleFilePath, bundleBytes, 0o644)
	if err != nil {
		return err
	}

	return nil
}

func createTLSConf(c *workloadapi.X509Context) (*tls.Config, error) {
	svid, err := getSVIDByHint(c, mysqlClientSVIDHint)
	if err != nil {
		return nil, err
	}

	certBytes, keyBytes, err := svid.Marshal()
	if err != nil {
		return nil, err
	}

	bundleBytes, err := c.Bundles.Bundles()[0].Marshal()
	if err != nil {
		return nil, err
	}

	rootCertPool := x509.NewCertPool()
	if ok := rootCertPool.AppendCertsFromPEM(bundleBytes); !ok {
		return nil, errors.New("failed to append PEM")
	}
	clientCert := make([]tls.Certificate, 0, 1)

	certs, err := tls.X509KeyPair(certBytes, keyBytes)
	if err != nil {
		return nil, err
	}

	clientCert = append(clientCert, certs)

	return &tls.Config{
		RootCAs:      rootCertPool,
		Certificates: clientCert,
	}, nil
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

func getSVIDByHint(c *workloadapi.X509Context, hint string) (*x509svid.SVID, error) {
	for _, svid := range c.SVIDs {
		if svid.Hint == hint {
			return svid, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("SVID not found for hint: %s", hint))
}
