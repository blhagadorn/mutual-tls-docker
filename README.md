### Mutual TLS Pod to Pod 

## Getting Started

1. Generate root key and certificates with CN as localhost:

```
openssl req -newkey rsa:2048 \
  -new -nodes -x509 \
  -days 3650 \
  -out cert.pem \
  -keyout key.pem \
  -subj "/C=US/ST=Minnesota/L=Minneapolis/O=Organization/OU=Unit/CN=localhost"
  ```  

2. `go run -v server.go`  
3. `go run -v client.go`  
4. Optionally, export cert and key as PKSC12 bundle
```
openssl pkcs12 -export -in cert.pem -inkey key.pem -out server.p12
```


