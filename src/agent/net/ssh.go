package net

import (
	"bytes"
	"code.google.com/p/go.crypto/ssh"
	"crypto"
	"crypto/dsa"
	"crypto/rsa"
	_ "crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
)

type keychain struct {
	keys []interface{}
}

func (k *keychain) Key(i int) (interface{}, error) {
	if i < 0 || i >= len(k.keys) {
		return nil, nil
	}
	switch key := k.keys[i].(type) {
	case *rsa.PrivateKey:
		return &key.PublicKey, nil
	case *dsa.PrivateKey:
		return &key.PublicKey, nil
	}
	panic("unknown key type")
}

func (k *keychain) Sign(i int, rand io.Reader, data []byte) (sig []byte, err error) {
	hashFunc := crypto.SHA1
	h := hashFunc.New()
	h.Write(data)
	digest := h.Sum(nil)
	switch key := k.keys[i].(type) {
	case *rsa.PrivateKey:
		return rsa.SignPKCS1v15(rand, key, hashFunc, digest)
	}
	return nil, errors.New("ssh: unknown key type")
}

func (k *keychain) loadPEM(file string) error {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	block, _ := pem.Decode(buf)
	if block == nil {
		return errors.New("ssh: no key found")
	}
	r, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	}
	k.keys = append(k.keys, r)
	return nil
}

type SSHClient struct {
	User           string
	ClientKeychain *keychain
	Host           string
	Session        ssh.Session
}

func NewSSHClient(user string, privateKey string, host string) *SSHClient {
	clientKeychain := new(keychain)
	block, _ := pem.Decode([]byte(privateKey))
	rsakey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	clientKeychain.keys = append(clientKeychain.keys, rsakey)
	return &SSHClient{
		User:           user,
		ClientKeychain: clientKeychain,
		Host:           host,
	}
}

func (c *SSHClient) init() {
	config := &ssh.ClientConfig{
		User: c.User,
		Auth: []ssh.ClientAuth{
			ssh.ClientAuthKeyring(c.ClientKeychain),
		},
	}
	client, err := ssh.Dial("tcp", c.Host, config)
	// defer client.Close()
	if err != nil {
		panic("Failed to dial:" + err.Error())
	}

	c.Session, _ = client.NewSession()

	// defer session.Close()
}

func (c *SSHClient) Run(command string) string {
	var b bytes.Buffer
	c.Session.Stdout = &b
	if err := c.Session.Run(command); err != nil {
		panic("Failed to run: " + err.Error())
	}
	return b.String()
}
