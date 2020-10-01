package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/spiffe/go-spiffe/v2/bundle/x509bundle"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Write "Hello, world!" to the response body
	io.WriteString(os.Stdout, "ANSWERING REQUEST\n")
	io.WriteString(w, "Hello, world!\n")
}

func main() {
	// Set up a /hello resource handler
	http.HandleFunc("/hello", helloHandler)

	// Listen to port 8443 and wait
	// log.Fatal(http.ListenAndServeTLS(":8443", "../cert/server.crt", "../cert/server.key", nil))

	// Create a CA certificate pool and add cert.pem to it
	caCert, err := ioutil.ReadFile("../cert/ca.crt")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	trustDomain, err := spiffeid.TrustDomainFromString("localhost")
	if err != nil {
		log.Fatalf("could not create trustdomain from string: %v", err)
	}
	authorizedSpiffeId, err := spiffeid.New("localhost", "client")
	if err != nil {
		log.Fatalf("could not create spiffe id from string: %v", err)
	}
	bundle, err := x509bundle.Load(trustDomain, "../cert/ca.crt")
	if err != nil {
		log.Fatalf("could not create new x509 bundle: %v", err)
	}
	// Create the TLS Config with the CA pool and enable Client certificate validation
	tlsConfig := &tls.Config{
		ClientCAs:             caCertPool,
		ClientAuth:            tls.RequireAndVerifyClientCert,
		VerifyPeerCertificate: tlsconfig.VerifyPeerCertificate(bundle, tlsconfig.AuthorizeID(authorizedSpiffeId), nil),
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
