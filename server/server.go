package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Write "Hello, world!" to the response body
	io.WriteString(os.Stdout, "ANSWERING REQUEST\n")
	io.WriteString(w, "Hello, world!\n")
}

func main() {
	// Set up a /hello resource handler
	http.HandleFunc("/hello", helloHandler)

	// Listen to port 8081 and wait
	log.Fatal(http.ListenAndServeTLS(":8443", "../cert/server.crt", "../cert/server.key", nil))
}