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

## mTLS with Intermediate Certificates

In the last example, we generated a root certificates and reused the same private key across the client and server, to make things simple.
Next, in  `03-client-server-mtls-intermediate`:

Generate the root certificate again:
```
openssl req -newkey rsa:2048 \
  -nodes -x509 \
  -days 3650 \
  -keyout key.pem \
  -out cert.pem \
  -subj "/C=US/ST=Montana/L=Bozeman/O=Organization/OU=Unit/CN=localhost" \
  -addext "subjectAltName = DNS:localhost"
```

Then, let's generate an intermediate certificate key for the client:  
```
openssl genrsa -out client_intermediate_key.pem 2048
```

Then the server:
```
openssl genrsa -out server_intermediate_key.pem 2048
```

Then let's generate some CSR's (Certificate Signing Requests). Because we are using SAN, we need some .cnf files. They are included in this directory for ease of use:

For the client:  
```
openssl req -new -key client_intermediate_key.pem -out client_intermediate_csr.pem -config client_csr.cnf
```

For the server:  
```
openssl req -new -key server_intermediate_key.pem -out server_intermediate_csr.pem -config server_csr.cnf
```

Next, let's sign the CSR's

For the client:

```
openssl x509 -req -days 3650 -in client_intermediate_csr.pem -CA cert.pem -CAkey key.pem -CAcreateserial -out client_intermediate_cert.pem -extfile client_ext.cnf
```

For the server:  

```
openssl x509 -req -days 3650 -in server_intermediate_csr.pem -CA cert.pem -CAkey key.pem -CAcreateserial -out server_intermediate_cert.pem -extfile server_ext.cnf
```

Lastly, let's veriify that all the signing actually worked:

For the client intermediate:
```
openssl verify -CAfile cert.pem client_intermediate_cert.pem
>client_intermediate_cert.pem: OK
```

For the server intermediate

```
openssl verify -CAfile cert.pem server_intermediate_cert.pem
>server_intermediate_cert.pem: OK
```
Perfect! The intermediate certificates have been verified by the root certificate `cert.pem`.

Alright, as you can see this is quite a few more steps, but now we have a better trust relationship.

Let's run the docker commands to verify it's all working:

```
docker build -t mtls-server -f Dockerfile.server . && docker run --rm  -p 8433:8433 -it --network host mtls-server:latest
```
Then in another terminal:  
```
docker build -t mtls-client -f Dockerfile.client . && docker run --rm -it --network host mtls-client:latest
> Hello, world WITH mutual TLS using intermediate certificates!
```

Success!
