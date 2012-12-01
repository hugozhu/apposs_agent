package net

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func Fetch(method string, base string, query url.Values) string {
	url1 := base
	var resp *http.Response
	var err error
	if method == "POST" {
		resp, err = http.PostForm(url1, query)
	} else {
		url1 += query.Encode()
		resp, err = http.Get(url1)
	}
	if err != nil {
		log.Panic("fetch url %s %s", url1, err)
	}
	defer resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == 200 {
		return string(bytes)
	} else {
		log.Printf("Fetch url %s error: %s", url1, string(bytes))
	}
	return ""
}
