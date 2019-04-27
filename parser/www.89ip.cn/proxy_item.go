package www_89ip_cn

import (
	"account_getter/fetcher"
	"account_getter/logger"
	"account_getter/models"
	"bytes"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func ParseProxyItem(contents []byte) models.ParserResult {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(contents))
	result := models.ParserResult{}
	doc.Find("table.layui-table tr").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		tds := s.Find("td")
		if tds.Length() == 0 {
			return
		}
		var tdsArr []string
		tds.Each(func(i int, s *goquery.Selection) {
			tdsArr = append(tdsArr, strings.Trim(s.Text(), " \n\r\t"))
		})
		//基本信息
		ip := tdsArr[0]
		port := tdsArr[1]
		proxyType := "http"
		//检查状态
		//status, _ := CheckStaus(ip, port, proxyType)
		status := false
		tmpProxyItem := &models.ProxyItem{
			Ip:         ip,
			Port:       port,
			Type:       proxyType,
			Status:     status,
			UpdateTime: time.Now().Format("2006-01-02 15:04:05"),
		}
		result.Items = append(result.Items, tmpProxyItem)

	})

	//go routine
	outChan := make(chan bool)
	for _, item := range result.Items {
		//fmt.Println(item.(*models.ProxyItem))
		if item, ok := item.(*models.ProxyItem); ok {
			go func(item *models.ProxyItem) {
				status, _ := CheckProxyStaus(item)
				logger.DebugLog.Printf("Checking Status %t:%s", status, item)
				item.Status = status
				item.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
				outChan <- true
			}(item)
		}
	}

	//back
	for i := 0; i < len(result.Items); i++ {
		<-outChan
	}

	return result
}

//检查代理状态
func CheckStaus(ip string, port string, proxyType string) (bool, error) {
	//使用sug接口进行检测
	checkUrl := "http://m.baidu.com/sugrec?prod=wise&wd=%E5%8C%97%E4%BA%AC"

	item := models.ProxyItem{
		Ip:   ip,
		Port: port,
		Type: proxyType,
	}

	_, err := fetcher.FetchByProxyItem(item, checkUrl)
	if err != nil {
		return false, err
	}
	return true, nil
}

//检查proxy状态
func CheckProxyStaus(item *models.ProxyItem) (bool, error) {
	//使用sug接口进行检测
	checkUrl := "http://m.baidu.com/sugrec?prod=wise&wd=%E5%8C%97%E4%BA%AC"
	newItem := models.ProxyItem{
		Ip:   item.Ip,
		Port: item.Port,
		Type: item.Type,
	}
	_, err := fetcher.FetchByProxyItem(newItem, checkUrl)
	if err != nil {
		return false, err
	}
	return true, nil
}
