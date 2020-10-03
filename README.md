### Mutual TLS Pod to Pod 
Implementation of mutual transport layer security (TLS) between two docker containers running Go binaries

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

## Using Docker

1. Build and start Docker containers
```
docker build -t mtls-server -f Dockerfile.server . && docker run --rm  -p 8433:8433 -it --network host mtls-server:latest
docker build -t mtls-client -f Dockerfile.client . && docker run --rm -it --network host mtls-client:latest
```

2. Analyze network traffic

Non-mutual TLS port traffic analysis:

```
docker run -it --network host --rm dockersec/tcpdump tcpdump -i any port 8080 -c 100 -A
```
Mutual TLS port traffic analysis:

```
docker run -it --network host --rm dockersec/tcpdump tcpdump -i any port 8443 -c 100 -A
```

