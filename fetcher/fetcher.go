package fetcher

import (
	"account_getter/models"
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

//获取url中的数据
func Fetch(url string) ([]byte, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	//request.Header.Add("Referer", "http://album.zhenai.com/u/1764131916")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36")
	return FetchByRequest(request)
}

//使用proxy进行访问
func FetchByProxyItem(item models.ProxyItem, reqUrl string) ([]byte, error) {
	//创建client
	proxyAddr := fmt.Sprintf("%s://%s:%s", item.Type, item.Ip, item.Port)
	proxy, err := url.Parse(proxyAddr)
	if err != nil {
		return nil, err
	}
	netTransport := &http.Transport{
		Proxy: http.ProxyURL(proxy),
		//MaxIdleConnsPerHost:   10,                             //每个host最大空闲连接
		ResponseHeaderTimeout: time.Second * time.Duration(3), //数据收发5秒超时
	}
	client := &http.Client{
		Timeout:   time.Second * 2,
		Transport: netTransport,
	}

	//创建request
	request, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36")

	//获取数据
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	//judge status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP_STATUS_CODE_ERROR:%d", resp.StatusCode)
	}
	bodyReader := bufio.NewReader(resp.Body)
	//judge charset
	e := determineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	return ioutil.ReadAll(utf8Reader)
}

//获取request数据
func FetchByRequest(r *http.Request) ([]byte, error) {
	//fetch
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//judge status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP_STATUS_CODE_ERROR:%d", resp.StatusCode)
	}
	bodyReader := bufio.NewReader(resp.Body)
	//judge charset
	e := determineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	return ioutil.ReadAll(utf8Reader)
}

//判断编码
func determineEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)
	if err != nil {
		return unicode.UTF8
	}
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}

//使用proxy进行访问
func FetchRequestByProxyItem(curRequest *models.Request, item *models.ProxyItem) ([]byte, error) {

	//创建client
	proxyAddr := fmt.Sprintf("%s://%s:%s", item.Type, item.Ip, item.Port)
	proxy, err := url.Parse(proxyAddr)
	if err != nil {
		return nil, err
	}
	netTransport := &http.Transport{
		Proxy: http.ProxyURL(proxy),
		//MaxIdleConnsPerHost:   10,                             //每个host最大空闲连接
		ResponseHeaderTimeout: time.Second * time.Duration(5), //数据收发5秒超时
	}
	client := &http.Client{
		Timeout:   time.Second * 5,
		Transport: netTransport,
	}

	//创建request
	request, err := http.NewRequest(curRequest.Method, curRequest.Url, strings.NewReader(curRequest.Data.Encode()))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.86 Safari/537.36")
	//获取数据
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	//judge status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP_STATUS_CODE_ERROR:%d", resp.StatusCode)
	}
	bodyReader := bufio.NewReader(resp.Body)
	//judge charset
	e := determineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	return ioutil.ReadAll(utf8Reader)
}

func BuildHttpBody(params map[string]string) string {
	return ""
}
