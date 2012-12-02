package net

import (
	"net/url"
	"testing"
)

func TestPOST(t *testing.T) {
	Fetch("GET", "http://www.baidu.com/s", url.Values{"wd": {"java"}})
}
