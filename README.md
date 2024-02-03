### Mutual TLS Pod to Pod 
Implementation of mutual transport layer security (TLS) between two docker containers running Go binaries

## Getting Started

1. Generate root key and certificates with CN (Common Name) and SAN (Subject Alternative Name) as `localhost`:

_Note: CN has been deprecated and most modern TLS libraries require SAN to to be set instead, including Golang's `net/http`.
Here, we set both, because some libraries still require CN to be set or use it as a fallback._

First, change directories to the mTLS directory (and be sure to run all commands out of this directory):
`cd 02-client-server-mtls`

Then let's generate a key and certificate.
```
openssl req -newkey rsa:2048 \
  -nodes -x509 \
  -days 3650 \
  -keyout key.pem \
  -out cert.pem \
  -subj "/C=US/ST=Montana/L=Bozeman/O=Organization/OU=Unit/CN=localhost" \
  -addext "subjectAltName = DNS:localhost"

```

2. `go run -v server-mtls.go`  
3. `go run -v client-mtls.go`
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

