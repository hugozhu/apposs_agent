package main

import (
	"agent"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	cmd := r.FormValue("cmd")
	host := r.FormValue("host")
	agent := hosts[host]
	io.WriteString(w, ssh_client.Run(cmd))
}

var (
	user       string
	privateKey string
	hosts      map[string]*net.SSHClient
)

func init() {
	tmp, _ := ioutil.ReadFile("/Users/hugozhu/.ssh/id_rsa")
	privateKey = string(tmp)
	user = "hugo"
	hosts = make(map[string]*net.SSHClient)
}
func main() {
	hosts["us"] = agent.New(user, privateKey, "us.myalert.info:22")
	hosts["jp"] = agent.New(user, privateKey, "jp.myalert.info:22")

	http.HandleFunc("/", handler)
	s := &http.Server{
		Addr:           ":1234",
		Handler:        nil,
		ReadTimeout:    1000 * time.Millisecond,
		WriteTimeout:   10000 * time.Microsecond,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
