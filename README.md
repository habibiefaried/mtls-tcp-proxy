# mtls-tcp-proxy
Mutual Authentication TLS encryption TCP proxy with golang

# Why?

I created this because of sometimes, it is not possible for us to establish secure connection and authentication between client and server for some reason (e.g no budget for VPNs). Forcing both parties to connect those services over TCP network, that is plaintext by design

![Alt text](/screenshot/unsecuredlink.png?raw=true "Unsecured link")

If somehow we manage to create secure proxy link, that stands between those client and server, then I think it's sufficient enough. 

# Certificate Setup
Navigate to provided CSR files provided.

```
cd certs
```

Generate the CA certificate and private key.

```
cfssl gencert -initca ca-csr.json | cfssljson -bare ca
```

Generate a server cert using the CSR provided. You can change hostname as you want the client connect to (in this case, `localhost`)

```
cfssl gencert  \
    -ca=ca.pem \
    -ca-key=ca-key.pem \
    -config=ca-config.json \
    -hostname=localhost,127.0.0.1 \
    -profile=mtlstcp server-csr.json | cfssljson -bare server
```

Generate a client cert using the CSR provided.

```
cfssl gencert \
  -ca=ca.pem \
  -ca-key=ca-key.pem \
  -config=ca-config.json \
  -profile=mtlstcp \
  client-csr.json | cfssljson -bare client
```

# Run

This is for testing purpose on localhost

Encryptor

```
CERT_PATH=./certs/client.pem KEY_PATH=./certs/client-key.pem ROOT_CERT_PATH=./certs/ca.pem BIND_PORT=10000 REMOTE_ADDR_PAIR=localhost:10001 ./main encryptor
```

Decryptor

```
CERT_PATH=./certs/server.pem KEY_PATH=./certs/server-key.pem ROOT_CERT_PATH=./certs/ca.pem BIND_PORT=10001 REMOTE_ADDR_PAIR=localhost:10002 ./main decryptor
```

TCP Server (netcat)

```
nc -nlvp 10002
```

Client (netcat)

```
nc -vvv localhost 10000
```

# Testing

```
Client ---> encryptor (port 10000) -> decryptor (port 10001) -> Server (port 10002)
```

Like this diagram, represent the real-world use case for this program

![Alt text](/screenshot/securedlink.png?raw=true "Secured link")

# Screenshots

This is the picture, testing successful with netcat, representing client and server

![Alt text](/screenshot/success.png?raw=true "Successful test")

When you try to connect directly to the server (decryptor). It is not valid TLS handshake

![Alt text](/screenshot/rejected.png?raw=true "Any non-tls connection will be rejected")