package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	// Request /hello over port 8080 via the GET method
	// r, err := http.Get("http://localhost:8080/hello")
	// Request /hello over HTTPS port 8443 via the GET method
	client := getHTTPSClientFromFile()
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

func getHTTPSClientFromFile() *http.Client {
	caCert, err := ioutil.ReadFile("cert.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Read the key pair to create certificate
	cert, err := tls.LoadX509KeyPair("client_intermediate_cert.pem", "client_intermediate_key.pem")
	if err != nil {
		log.Fatal(err)
	}
	// Create an HTTPS client and supply the created CA pool
	// CA pool is a group of certificates
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{cert},
			},
		},
	}

	return client
}
