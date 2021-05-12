# mtls-tcp-proxy
Mutual Authentication TLS encryption TCP proxy with golang

# Certificate Setup
Navigate to provided CSR files provided.

```
cd certs
```

Generate the CA certificate and private key.

```
cfssl gencert -initca ca-csr.json | cfssljson -bare ca
```

Generate a server cert using the CSR provided. You can change hostname as you want the client connect to.

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