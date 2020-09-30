package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spiffe/go-spiffe/v2/bundle/x509bundle"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
)

func main() {
	// Create a CA certificate pool and add cert.pem to it
	// caCert, err := ioutil.ReadFile("../cert/ca.crt")
	// if err != nil {
	// 	log.Fatalf("could not open certificate file: %v", err)
	// }
	// caCertPool := x509.NewCertPool()
	// caCertPool.AppendCertsFromPEM(caCert)
	// x509bundle.Source
	// trustDomain, err := spiffeid.FromString("spiffe://localhost/")
	trustDomain, err := spiffeid.TrustDomainFromString("localhost")
	if err != nil {
		log.Fatalf("could not create trustdomain from string: %v", err)
	}
	authorizedSpiffeId, err := spiffeid.New("localhost", "server")
	if err != nil {
		log.Fatalf("could not create spiffe id from string: %v", err)
	}
	// source, err := workloadapi.NewX509Source(context.Background())
	// if err != nil {
	// 	log.Fatalf("could not create new x509 source: %v", err)
	// }
	// bundle, err := x509bundle.Source.GetX509BundleForTrustDomain(source, trustDomain)
	// if err != nil {
	// 	log.Fatalf("could not bundle from trust domain: %v", err)
	// }
	bundle, err := x509bundle.Load(trustDomain, "../cert/ca.crt")
	if err != nil {
		log.Fatalf("could not create new x509 bundle: %v", err)
	}

	// Create a HTTPS client and supply the created CA pool
	// client := &http.Client{
	// 	Transport: &http.Transport{
	// 		TLSClientConfig: &tls.Config{
	// 			RootCAs:               caCertPool,
	// 			InsecureSkipVerify:    true,
	// 			VerifyPeerCertificate: tlsconfig.WrapVerifyPeerCertificate(tls.Config.VerifyPeerCertificate, bundle, authorizer, opts...),
	// 		},
	// 	},
	// }
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsconfig.TLSClientConfig(bundle, tlsconfig.AuthorizeID(authorizedSpiffeId)),
		},
	}

	// Request /hello over port 8443 via the GET method
	r, err := client.Get("https://localhost:8443/hello")
	if err != nil {
		log.Fatal(err)
	}

	// Read the response body
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Print the response body to stdout
	fmt.Printf("%s\n", body)
}
