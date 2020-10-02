package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/spiffe/go-spiffe/v2/bundle/x509bundle"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	"github.com/spiffe/go-spiffe/v2/svid/x509svid"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {

	io.WriteString(os.Stdout, "ANSWERING REQUEST\n")

	// Extract SPIFFE ID from URI SAN
	clientSpiffeId, err := x509svid.IDFromCert(r.TLS.PeerCertificates[0])
	if err != nil {
		log.Fatalf("could not get spiffe id from client cert: %v", err)
	}

	// Create the authorized SPIFFE ID that can connect to the server
	authorizedSpiffeId, err := spiffeid.New("localhost", "client")
	if err != nil {
		log.Fatalf("could not create authorized spiffe id from string: %v", err)
	}

	// Compare received SPIFFE ID with expected one
	var matcher = spiffeid.MatchID(authorizedSpiffeId)
	err = matcher(clientSpiffeId)
	if err != nil {
		message := fmt.Sprintf("failed to match spiffe id: %v", err)
		log.Printf(message)
		http.Error(w, message, http.StatusUnauthorized)
		return
	}

	// Write "Hello, world!" to the response body
	io.WriteString(w, "Hello, world!\n")
}

func main() {
	// Set up a /hello resource handler
	http.HandleFunc("/hello", helloHandler)

	// Create a CA certificate pool and add ca.crt to it
	caCert, err := ioutil.ReadFile("../cert/ca.crt")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create trust domain
	trustDomain, err := spiffeid.TrustDomainFromString("localhost")
	if err != nil {
		log.Fatalf("could not create trustdomain from string: %v", err)
	}

	// Load keys for this trust domain
	bundle, err := x509bundle.Load(trustDomain, "../cert/ca.crt")
	if err != nil {
		log.Fatalf("could not create new x509 bundle: %v", err)
	}

	// Create the TLS Config enable Client certificate validation
	tlsConfig := &tls.Config{
		ClientCAs:             caCertPool,
		ClientAuth:            tls.RequireAndVerifyClientCert,
		VerifyPeerCertificate: tlsconfig.VerifyPeerCertificate(bundle, tlsconfig.AuthorizeAny(), nil),
	}
	tlsConfig.BuildNameToCertificate()

	// Create a Server instance to listen on port 8443 with the TLS config
	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig,
	}

	// Listen to HTTPS connections with the server certificate and wait
	log.Fatal(server.ListenAndServeTLS("../cert/server.crt", "../cert/server.key"))
}
