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

func Encrypt() {
	localAddr := "0.0.0.0:" + os.Getenv("BIND_PORT")
	fmt.Printf("Listening: %v\nProxying & Encrypting: %v\n\n", localAddr, os.Getenv("REMOTE_ADDR_PAIR"))

	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("error accepting connection", err)
			continue
		}
		go func() {
			fmt.Printf("Reading %v as certificate, %v as key and %v as root certificate\n", os.Getenv("CERT_PATH"), os.Getenv("KEY_PATH"), os.Getenv("ROOT_CERT_PATH"))
			cert, err := tls.LoadX509KeyPair(os.Getenv("CERT_PATH"), os.Getenv("KEY_PATH"))

			caCert, err := ioutil.ReadFile(os.Getenv("ROOT_CERT_PATH"))
			if err != nil {
				log.Fatalf("failed to load cert: %s", err)
			}

			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)

			tlsConfig := &tls.Config{
				Certificates: []tls.Certificate{cert}, // this certificate is used to sign the handshake
				RootCAs:      caCertPool,              // this is used to validate the server certificate
			}
			tlsConfig.BuildNameToCertificate()

			// Setup TLS done

			conn2, err := tls.Dial("tcp", os.Getenv("REMOTE_ADDR_PAIR"), tlsConfig)
			if err != nil {
				log.Println("error dialing remote addr", err)
				return
			}

			// this is required to complete the handshake and populate the connection state
			// we are doing this so we can print the peer certificates prior to reading / writing to the connection
			err = conn2.Handshake()
			if err != nil {
				log.Printf("failed to complete handshake: %s\n", err)
				return
			}

			tag := fmt.Sprintf("[%s -> %s]", conn2.LocalAddr(), conn2.RemoteAddr())
			log.Printf("%s connect", tag)

			if len(conn2.ConnectionState().PeerCertificates) > 0 {
				log.Printf("%s client common name: %+v", tag, conn2.ConnectionState().PeerCertificates[0].Subject.CommonName)
			}

			go io.Copy(conn2, conn)
			io.Copy(conn, conn2)
			conn2.Close()
			conn.Close()
		}()
	}
}
