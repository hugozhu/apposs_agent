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
	"log"
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
	Connection     *ssh.ClientConn
	Connected      bool
}

/**
 * A SSHClient keeps long connection to remote server and run a command at a time
 */
func NewSSHClient(user string, privateKey string, host string) *SSHClient {
	clientKeychain := new(keychain)
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		panic("Failed to decode ssh private key")
	}
	rsakey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	clientKeychain.keys = append(clientKeychain.keys, rsakey)
	c := &SSHClient{
		User:           user,
		ClientKeychain: clientKeychain,
		Host:           host,
	}
	return c
}

func (c *SSHClient) connect() error {
	if !c.Connected {
		config := &ssh.ClientConfig{
			User: c.User,
			Auth: []ssh.ClientAuth{
				ssh.ClientAuthKeyring(c.ClientKeychain),
			},
		}
		conn, err := ssh.Dial("tcp", c.Host, config)
		// defer conn.Close()
		if err != nil {
			return err
		}
		log.Println("[Info] connected to " + c.Host)

		c.Connection = conn
		c.Connected = true
	}
	return nil
}

func (c *SSHClient) Run(command string) (string, error) {
	var err error
	defer func() {
		/**
		 *  Reconnect if there is problems
		 */
		if x := recover(); x != nil {
			c.Close()
			c.Connected = false
			err = errors.New(fmt.Sprintf("[Error] failed to connect %s: %v", c.Host, x))
		}
	}()

	err = c.connect()
	if err == nil {
		var b bytes.Buffer
		session, _ := c.Connection.NewSession()
		defer session.Close()

		session.Stdout = &b
		if err = session.Run(command); err != nil {
			log.Println("Failed to run: " + err.Error())
		}
		return b.String(), err
	}

	return "error", err
}

func (c *SSHClient) Close() {
	if c.Connection != nil {
		c.Connection.Close()
	}
}
