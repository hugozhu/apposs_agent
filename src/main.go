package main

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
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World\n")
}

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

var (
	privateKey = `-----BEGIN RSA PRIVATE KEY-----
-----END RSA PRIVATE KEY-----`

	clientKeychain = new(keychain)
)

func init() {
	block, _ := pem.Decode([]byte(privateKey))
	rsakey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	clientKeychain.keys = append(clientKeychain.keys, rsakey)
}

func main() {
	config := &ssh.ClientConfig{
		User: "hugo",
		Auth: []ssh.ClientAuth{
			ssh.ClientAuthKeyring(clientKeychain),
		},
	}
	client, err := ssh.Dial("tcp", "", config)
	defer client.Close()
	if err != nil {
		panic("Failed to dial:" + err.Error())
	}
	session, err := client.NewSession()
	defer session.Close()
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("ls -l"); err != nil {
		panic("Failed to run: " + err.Error())
	}
	fmt.Println(b.String())
}

func main2() {
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
