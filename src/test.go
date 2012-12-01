package main

import (
	"agent/net"
	"io/ioutil"
	"log"
)

func main() {
	privateKey, _ := ioutil.ReadFile("/Users/hugozhu/.ssh/id_rsa")
	ssh_client := net.NewSSHClient("hugo", string(privateKey), "us.myalert.info:22")
	log.Println(ssh_client.Run("ls -l"))
	log.Println(ssh_client.Run("whoami"))

}
