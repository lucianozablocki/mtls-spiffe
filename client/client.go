package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/spiffe/go-spiffe/v2/bundle/x509bundle"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
)

func main() {
	// Create a CA certificate pool and add ca.crt to it
	caCert, err := ioutil.ReadFile("../cert/ca.crt")
	if err != nil {
		log.Fatalf("could not open certificate file: %v", err)
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
		log.Fatalf("could not load x509 bundle from cert: %v", err)
	}

	// Create the SPIFFE ID that's authorized to connect
	authorizedSpiffeId, err := spiffeid.New("localhost", "server")
	if err != nil {
		log.Fatalf("could not create authorized spiffe id from string: %v", err)
	}

	// Load client certificate to be presented
	certificate, err := tls.LoadX509KeyPair("../cert/client.crt", "../cert/client.key")
	if err != nil {
		log.Fatalf("could not load client certificate: %v", err)
	}

	// Create a HTTPS client, with the client certificate, and GO-SPIFFE'S VerifyPeerCertificate verifier to comply with SPIFFE standards
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:               caCertPool,
				InsecureSkipVerify:    true,
				VerifyPeerCertificate: tlsconfig.VerifyPeerCertificate(bundle, tlsconfig.AuthorizeID(authorizedSpiffeId), nil),
				Certificates:          []tls.Certificate{certificate},
			},
		},
	}

	// Request /hello over port 8443 via GET method
	r, err := client.Get("https://localhost:8443/hello")
	if err != nil {
		log.Fatalf("could not make GET request: %v", err)
	}

	// Read response body
	defer r.Body.Close()
	dump, err := httputil.DumpResponse(r, true)
	if err != nil {
		log.Fatalf("could not dump response: %v", err)
	}

	// Print response to stdout
	fmt.Printf("%s\n", dump)
}
