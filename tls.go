package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
)

func (cfg *Config) LoadTLS() error {
	caCert, err := ioutil.ReadFile(cfg.SSLCertFile)
	if err != nil {
		log.Println("[ERR] Loading certificate", err)
		return err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
		//ClientCAs: caCertPool,
		// NoClientCert
		// RequestClientCert
		// RequireAnyClientCert
		// VerifyClientCertIfGiven
		// RequireAndVerifyClientCert
		//ClientAuth: tls.RequireAndVerifyClientCert,
	}
	//tlsConfig.BuildNameToCertificate()
	cfg.TLSConfig = tlsConfig
	return nil
}
