package agent

import (
	"agent/net"
)

type SSHAgent struct {
	Commands  chan string
	Client    *net.SSHClient
	quit_chan chan bool
}

func New(user string, privateKey string, host string) *SSHAgent {
	agent := new(SSHAgent)
	agent.quit_chan = make(chan bool)
	agent.Client = net.NewSSHClient(user, privateKey, host)
	agent.init()
	return agent
}

func (a *SSHAgent) init() {
	a.Commands = make(chan string, 100)
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
	a.quit_chan <- true
	a.Client.Close()
	a.init()
}

func (a *SSHAgent) Run(callback func(output string, err error)) {
	go func() {
		for {
			select {
			case <-a.quit_chan:
				return
			case cmd := <-a.Commands:
				result, err := a.Client.Run(cmd)
				callback(result, err)
			}
		}
	}()
}
