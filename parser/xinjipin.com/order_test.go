package xinjipin_com

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestParseOrder(t *testing.T) {

	contents, err := ioutil.ReadFile("./order_test.html")
	if err != nil {
		panic(err)
	}
	list := ParseOrder(contents)
	if len(list.Items) != 5 {
		t.Errorf("proxy_item length not as expect")
	}
	fmt.Print(list)
}
