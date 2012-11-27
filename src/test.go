package main

import (
	"agent/net"
	"log"
)

var (
	privateKey = `-----BEGIN RSA PRIVATE KEY-----
-----END RSA PRIVATE KEY-----`
)

func main() {
	log.Println("Hello")
	ssh_client := net.NewSSHClient("hugo", privateKey, "localhost:22")
	ssh_client.Connect()
}
