package www_89ip_cn

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestParseProxyItem(t *testing.T) {

	contents, err := ioutil.ReadFile("./proxy_item_test.html")
	if err != nil {
		panic(err)
	}
	list := ParseProxyItem(contents)
	if len(list.Items) != 15 {
		t.Errorf("proxy_item length not as expect")
	}
	fmt.Print(list)
}

func TestCheckStaus(t *testing.T) {
	status, _ := CheckStaus("127.0.1.1", "8001", "http")
	if status != false {
		t.Errorf("check status not as expect")
	}
}
