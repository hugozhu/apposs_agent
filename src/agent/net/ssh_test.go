package net

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestCommands(t *testing.T) {
	tmp, _ := ioutil.ReadFile(os.Getenv("HOME") + "/.ssh/test_id_rsa")
	privateKey := string(tmp)
	user := "hugo"
	host := "test_host:22"
	ssh_client := NewSSHClient(user, privateKey, host)

	output, err := ssh_client.Run("ls -l")
	t.Log(output)
	if err != nil {
		t.Logf("Error: %v", err)
		t.Error("should not return error")
	}

	output, err = ssh_client.Run("whoami")
	t.Log(output)
	if err != nil {
		t.Logf("Error: %v", err)
		t.Error("should not return error")
	}

	output, err = ssh_client.Run("not_exist_command")
	t.Log(output)
	if err != nil {
		t.Logf("Error: %v", err)
	} else {
		t.Error("not exist command should return error")
	}
}
