package proxy

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
)

func Decrypt() {
	caCert, err := ioutil.ReadFile(os.Getenv("ROOT_CERT_PATH"))
	if err != nil {
		log.Fatalf("failed to load cert: %s", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cert, err := tls.LoadX509KeyPair(os.Getenv("CERT_PATH"), os.Getenv("KEY_PATH"))

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},        // server certificate which is validated by the client
		ClientCAs:    caCertPool,                     // used to verify the client cert is signed by the CA and is therefore valid
		ClientAuth:   tls.RequireAndVerifyClientCert, // this requires a valid client certificate to be supplied during handshake
	}

	ln, err := tls.Listen("tcp", "0.0.0.0:"+os.Getenv("BIND_PORT"), tlsConfig)
	if err != nil {
		log.Fatalf("failed to create listener: %s", err)
	}

	log.Println("listen: ", ln.Addr())

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("failed to accept conn: %s", err)
			continue
		}

		tlsConn, ok := conn.(*tls.Conn)
		if !ok {
			log.Printf("failed to cast conn to tls.Conn")
			continue
		}

		go func() {
			tag := fmt.Sprintf("[%s -> %s]", tlsConn.LocalAddr(), tlsConn.RemoteAddr())
			log.Printf("%s accept", tag)

			defer tlsConn.Close()

			// this is required to complete the handshake and populate the connection state
			// we are doing this so we can print the peer certificates prior to reading / writing to the connection
			err := tlsConn.Handshake()
			if err != nil {
				log.Printf("failed to complete handshake: %s", err)
				return
			}

			if len(tlsConn.ConnectionState().PeerCertificates) > 0 {
				log.Printf("%s client common name: %+v", tag, tlsConn.ConnectionState().PeerCertificates[0].Subject.CommonName)
			}

			conn2, err := net.Dial("tcp", os.Getenv("REMOTE_ADDR_PAIR"))
			if err != nil {
				log.Println("error dialing remote addr", err)
				return
			}

			go io.Copy(conn2, conn)
			io.Copy(conn, conn2)
			conn2.Close()
			conn.Close()

		}()
	}
}
