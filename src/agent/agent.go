package agent

import (
	"agent/net"
)

type SSHAgent struct {
	Commands chan string
	Client   *net.SSHClient
	quit     chan bool
}

func New(user string, privateKey string, host string) *SSHAgent {
	agent := new(SSHAgent)
	agent.quit = make(chan bool)
	agent.Client = net.NewSSHClient(user, privateKey, host)
	agent.init()
	return agent
}

func (a *SSHAgent) init() {
	agent.Commands = make(string, 100)
}

func (a *SSHAgent) Send(cmd string) {
	if cap(a.Commands) == len(a.Commands) {
		//buffer full, the remote server has something wrong
		a.Shutdown()
	} else {
		a.Commands <- cmd
	}
}

func (a *SSHAgent) Shutdown() {
	quit <- true
	a.Client.Close()
	agent.init()
}

func (a *SSHAgent) Run() {
	go func() {
		for {
			select {
			case <-quit:
				return
			case cmd := <-commands:
				result := a.Client.Run(cmd)
			}
		}
	}()
}
