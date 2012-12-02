package main

import (
	"agent"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	cmd := r.FormValue("cmd")
	host := r.FormValue("host")
	agent := findAgent(host)
	if agent != nil {
		agent.Send(cmd)
		io.WriteString(w, "OK")
	} else {
		io.WriteString(w, "Error")
	}
}

/**
 * todo: call center to get connection info
 * use connection info as key to cache agent
 */
func findAgent(host string) *agent.SSHAgent {
	return nil
}

var (
	user       string
	privateKey string
	hosts      map[string]*agent.SSHAgent
)

func init() {
	tmp, _ := ioutil.ReadFile(os.Getenv("HOME") + "/.ssh/test_id_rsa")
	privateKey = string(tmp)
	user = "hugo"
	hosts = make(map[string]*agent.SSHAgent)
}

func main() {
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
