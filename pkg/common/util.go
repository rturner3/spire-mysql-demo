package common

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	"github.com/spiffe/go-spiffe/v2/svid/x509svid"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
)

const (
	// SPIRE SVID related constants
	svidDir    = "/spire/certs"
	bundleFile = "bundle.0.pem"
	certFile   = "svid.0.pem"
	keyFile    = "svid.0.key"

	// MySQL related constants
	mysqlServerSVIDHint = "mysql-server"
	mysqlHost           = "mysql.mysql.svc.cluster.local"
	mysqlPort           = "3306"
	mysqlTLSConfigName  = "spire-ssl"
)

var (
	certFilePath   = fmt.Sprintf("%s/%s", svidDir, certFile)
	keyFilePath    = fmt.Sprintf("%s/%s", svidDir, keyFile)
	bundleFilePath = fmt.Sprintf("%s/%s", svidDir, bundleFile)

	// SPIFFE ID for MySQL Server
	mysqlServerSPIFFEID = spiffeid.RequireFromString("spiffe://example.org/mysql/server")
)

func WriteMySQLServerSVIDFiles(c *workloadapi.X509Context) error {
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

	if err = os.WriteFile(certFilePath, certBytes, 0o644); err != nil {
		return err
	}

	if err = os.WriteFile(keyFilePath, keyBytes, 0o644); err != nil {
		return err
	}

	if err = os.WriteFile(bundleFilePath, bundleBytes, 0o644); err != nil {
		return err
	}

	return nil
}

func NewMySQLDBWithSPIRETLSConfig(c *workloadapi.X509Context, mysqlUser string, dbName string, svidHint string) (*sql.DB, error) {
	// Create TLS config with client certificates
	tlsConf, err := createTLSConf(c, svidHint)
	if err != nil {
		log.Printf("Failed to create MySQL TLS config: %v", err)
		return nil, err
	}

	if err = mysql.RegisterTLSConfig(mysqlTLSConfigName, tlsConf); err != nil {
		log.Printf("Failed to register MySQL TLS config: %v", err)
		return nil, err
	}

	// Format is specified https://github.com/go-sql-driver/mysql#dsn-data-source-name
	dbConnectionString := fmt.Sprintf("%s@tcp(%s:%s)/%s?tls=%s", mysqlUser, mysqlHost, mysqlPort, dbName, mysqlTLSConfigName)

	db, err := sql.Open("mysql", dbConnectionString)
	if err != nil {
		log.Printf("Failed to open MySQL database: %v", err)
		return nil, err
	}
	return db, nil
}

func createTLSConf(c *workloadapi.X509Context, svidHint string) (*tls.Config, error) {
	var err error
	svid := c.DefaultSVID()
	if svidHint != "" {
		svid, err = getSVIDByHint(c, svidHint)
		if err != nil {
			return nil, err
		}
	}
	return tlsconfig.MTLSClientConfig(svid, c.Bundles, tlsconfig.AuthorizeID(mysqlServerSPIFFEID)), nil
}

func getSVIDByHint(c *workloadapi.X509Context, hint string) (*x509svid.SVID, error) {
	for _, svid := range c.SVIDs {
		if svid.Hint == hint {
			return svid, nil
		}
	}
	return nil, fmt.Errorf("svid not found for hint: %s", hint)
}
