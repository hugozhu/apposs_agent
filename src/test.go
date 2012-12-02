package main

import (
	"agent"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var (
	privateKey string
	user       string
	host       string
)

func init() {
	tmp, _ := ioutil.ReadFile(os.Getenv("HOME") + "/.ssh/test_id_rsa")
	privateKey = string(tmp)
	user = "hugo"
	host = "test_host:22"
}

func main() {
	my_agent := agent.New(user, privateKey, host)
	my_agent.Run(func(output string, err error) {
		log.Println(output, err)
	})
	my_agent.Send("ls -l")
	my_agent.Send("whoami")
	time.Sleep(5 * time.Second)
}
