package main

import (
	"agent"
	"io/ioutil"
	"log"
	"time"
)

var (
	privateKey string
	user       string
	host       string
)

func init() {
	tmp, _ := ioutil.ReadFile("/Users/hugozhu/.ssh/mefans_id_rsa")
	privateKey = string(tmp)
	user = "hugo"
	host = "us.myalert.info:22"
}

func main() {
	my_agent := agent.New(user, privateKey, host)
	my_agent.Run(func(output string, err error) {
		log.Println(output, err)
	})
	my_agent.Send("ls -l")
	time.Sleep(5 * time.Second)
}
