package net

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func Fetch(method string, base string, query url.Values) (string, error) {
	url1 := base
	var resp *http.Response
	var err error
	if method == "POST" {
		resp, err = http.PostForm(url1, query)
	} else {
		url1 += query.Encode()
		resp, err = http.Get(url1)
	}
	defer resp.Body.Close()

	if err == nil {
		bytes, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode == 200 {
			return string(bytes), nil
		} else {
			log.Printf("Fetch url %s error: %s", url1, string(bytes))
			return string(bytes), errors.New(resp.Status)
		}
	} else {
		log.Printf("fetch url %s %s", url1, err)
	}
	return "", err
}
