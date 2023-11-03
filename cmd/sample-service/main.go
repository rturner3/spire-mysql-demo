package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/rturner3/spire-mysql-demo/pkg/common"
	"github.com/rturner3/spire-mysql-demo/pkg/store"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	usersAPIPath = "/api/v1/users"

	// SPIRE Agent socket path
	socketPath = "unix:///run/spire/sockets/agent.sock"

	mysqlUser   = "spire-mysql-client"
	mysqlDBName = "spiredemo"
)

type handler struct {
	dbStore *store.Store
}

func (h *handler) list(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users, err := h.dbStore.ListUsers(r.Context())
	if err != nil {
		writeErr(w, err)
		return
	}

	data, err := json.Marshal(users)
	if err != nil {
		writeErr(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *handler) create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeErr(w, err)
		return
	}

	var user store.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		writeErr(w, err)
		return
	}

	err = h.dbStore.CreateUser(r.Context(), user)
	if err != nil {
		writeErr(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "user created"}`))
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// Wait for an os.Interrupt signal
	go waitForCtrlC(cancel)

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

	db, err := common.NewMySQLDBWithSPIRETLSConfig(x509Context, mysqlUser, mysqlDBName, "")
	if err != nil {
		log.Fatalf("Failed to create MySQL Client: %v", err)
	}
	defer db.Close()

	h := &handler{
		dbStore: store.New(db),
	}

	// Start X.509 watcher
	go startWatcher(ctx, client, h)

	log.Printf("Starting API handlers")
	// Add API handlers
	http.HandleFunc(usersAPIPath, func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			h.list(w, req)
		case http.MethodPost:
			h.create(w, req)
		}
	})
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func startWatcher(ctx context.Context, client *workloadapi.Client, h *handler) {
	// Start a watcher for X.509 SVID updates
	doneCh := make(chan struct{}, 1)
	go func() {
		err := client.WatchX509Context(ctx, &x509Watcher{
			h: h,
		})
		if err != nil && status.Code(err) != codes.Canceled {
			log.Fatalf("Error watching X.509 context: %v", err)
		}
		doneCh <- struct{}{}
	}()
	<-doneCh
}

type x509Watcher struct {
	h *handler
}

// OnX509ContextUpdate is run every time an SVID is updated
func (w *x509Watcher) OnX509ContextUpdate(c *workloadapi.X509Context) {
	if err := common.LogSVIDs(c); err != nil {
		log.Printf("Failed to log SVIDs: %v", err)
		return
	}

	// Create new DB instance with udpate TLS config
	db, err := common.NewMySQLDBWithSPIRETLSConfig(c, mysqlUser, mysqlDBName, "")
	if err != nil {
		log.Printf("Failed to create MySQL Client: %v", err)
		return
	}

	// Update DB instance in store
	w.h.dbStore.UpdateDB(db)
	log.Printf("Successfully updated DB client TLS config")
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

func writeErr(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, `{"error": "%s"}`, err.Error())
}
