package main

import (
	"fmt"
	"github.com/habibiefaried/mtls-tcp-proxy/proxy"
	"os"
)

//isValid will return string if corresponding env is not set. will return empty if all set
func isValid() string {
	if os.Getenv("CERT_PATH") == "" {
		return "CERT_PATH"
	}

	if os.Getenv("KEY_PATH") == "" {
		return "KEY_PATH"
	}

	if os.Getenv("ROOT_CERT_PATH") == "" {
		return "ROOT_CERT_PATH"
	}

	if os.Getenv("BIND_PORT") == "" {
		return "BIND_PORT"
	}

	if os.Getenv("REMOTE_ADDR_PAIR") == "" {
		return "REMOTE_ADDR_PAIR"
	}

	return ""
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./main <encryptor|decryptor>. And set env var accordingly")
	} else {
		s := isValid()

		if s == "" {
			if os.Args[1] == "encryptor" {
				fmt.Println("encryptor")
				proxy.Encrypt()
			} else if os.Args[1] == "decryptor" {
				fmt.Println("decryptor")
			} else {
				fmt.Println(os.Args[1] + " is not recognized")
			}
		} else {
			fmt.Printf("Env var %v is not set\n", s)
		}
	}
}
