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
MIIEoAIBAAKCAQEA7+MgETklRNn5Z03txh0OMX6B+NXHBmegT3Fbsihl1rx6IiUq
J5AOH/38892PR/TkxjvIAJVTwp9nznRgVYtIBpgi6BbgzhPBNfMQQSjGc7QRpSr+
NGWDyav2XH+h4EocAK5/Fo97koBoG13HnNF1iAUFJuf8SvEwGT4g2fWgq6gl521g
hlv72CeMW+LlWcfcV3UoU2XL6359OrUlo6L7Qy4MzKTXjjqMGwejvflvPuXb63Pa
47CYRF5dAJXYVG93I67qaeyE7j7sr0rvmV8kIrY+o7dreXYMz5i+g+iv43kPvQIF
X6pUs/vZ6pYh20xYr01mI0xNYvqjn/fQuGEnqwIBIwKCAQB0hELGiXiAhyFeD+iE
zEi49E3CAW9MQPYX9TsqpehSW4vHcSMah8xYrpDOOGoqQ7+TfB9QvY8VY15OVchk
EHNw6s8gRaBkGDlGFvlEOGBkFaIrpysDgcOrGQKh9NmY04nxs9dUGc5OeOIb67i5
4hSD5S5jWrxe6i1OFtzBozDIZqygPwpF4wBPNLS9QXrfwOfIY9T43ZLkpoi8Ab6W
ao3qiBaw0ZiZHxsLBhtM9ozoasPi7Bn76/3ukDcw4Erpblf3I9oIg+V5et+eoES9
AxpPEKMA2ob7fSmLIgEvUCxYa8s+yByATuZO/raf1dcobvTLECX1ZwscMOKmbKgE
a0UbAoGBAPtP5LUBk44/ih0xZ/5bJ6X/lO9SxHcVTMyOhPaS6Tl4w7bTh5cCBVRS
b+pv8lSk1EV0PxmbYFSijcaYDn47kkpnyS0b3pMqaMSaLnMe5wm7YptRP+yAQRdA
tscmwr/M4+3NBcNEr9kcQ9753GHcq78SJZZ2oxkiGvR06S4iPlxZAoGBAPRcrIdP
QF3mfDdKe565fHdC1nN4uwhKvLVHO6at3mFgBtv2fMojyAAbTuqCZdphNcbvTqnn
TEfj1Nl8XzhWY2iR5zil71WNftQo4ZSa+T6XXo/2WzvkYHLZ2DEo6kl+qCUjb1K4
dij1FPtraubnngsKJK6za6Tot0XyiUirJ9OjAoGAQJ+D9Ae4OoVeB4G7p85MBh0m
TCs51XqBdm3Ka0pZOqoGcNdOwG474nuZHwbX6eE95flSDeYgFcNmSQKHYknVKRNf
neKYUbpytjZGdV+3wKyHEf7zLi+btYWrV8gjcyYOuX3c6RGpj5mNy59V7UdfXQSq
lGemSEqZN4vG7psXWZMCgYAN9qrF50WBtWYvC5IJEena7eBeXquaEuY2PpWxwMrg
/ixHFWY3b8JJJiHEQfcxDN6AZXI1mDA+n039zJe6E5Cm48QR3ZilzZI4AlYIfeJb
WRtYrv3mKk6neimcaLzuQcB3JpilpCQCVyXFOVaQ2gkH8fN3snst/qtjFSxqjW/9
dwKBgFrJHUwwC2cw6xLQ6x6YIgW60sp65pbPcAb6AtATBszZ0Z4aSSKZKBIHE7/E
kwItLvM74koiBu48eqe31MauN25rjl8gaTmwftYSoD4BuHRzi5qINgU+vj1YVGYj
gWIcSL9GCmH/ukGyPuJ6EcqMUazHx1E2Tz5J+Bx44MM32mRj
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
	client, err := ssh.Dial("tcp", "us.myalert.info:22", config)
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
