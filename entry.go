package main

import (
	"account_getter/engine"
	"account_getter/fetcher"
	"account_getter/logger"
	"account_getter/models"
	"account_getter/scheduler"
	"account_getter/server"
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"account_getter/parser/xinjipin.com"

	"account_getter/parser/www.89ip.cn"
)

func main() {

	updateFlag := flag.Bool("update", true, "是否需要重新更新数据")
	workerNum := flag.Int("worker-num", 5, "普通worker数")
	httpPort := flag.Int("http-port", 8000, "HTTP服务端口")
	flag.Parse()

	//开启httpserver
	go func() {
		port := strconv.Itoa(*httpPort)
		fmt.Printf("open http://localhost:%s/home\n", port)
		server.RunHttpServ(port)

	}()

	//如果需要update数据
	if *updateFlag {
		logger.DebugLog.Println("init engine")
		fmt.Println("crawler engine runing...")
		go runEngine(*workerNum)
	}

	for {
		select {
		case <-time.Tick(time.Second * 100):
			logger.DebugLog.Println("living...")
		}
	}

}

func runEngine(workerNum int) {

	ret := initEngine()
	if ret == false {
		fmt.Println("INIT_ERROR")
		return
	}

	//获取种子
	var proxySeeds []models.Request
	var contactSeeds []models.Request
	proxySeeds = getPorxySeeds()
	contactSeeds = getCotactSeeds()
	seeds := append(proxySeeds, contactSeeds...)

	//初始化engine
	e := engine.QueueEngine{
		Scheduler:   &scheduler.QueueScheduler{},
		WorkerCount: workerNum,
	}

	//爬取账号
	go func() {
		e.Run(seeds...)
	}()

	//记录
	go func() {
		for {
			statInfoMap := make(map[string]interface{})
			statInfoMap["allNum"] = len(contactSeeds)
			statInfoMap["doneNum"] = e.OrderDoneCount
			logger.DebugLog.Printf("PROCESS:%d/%d\n", statInfoMap["doneNum"], statInfoMap["allNum"])
			statInfoBytes, _ := json.Marshal(statInfoMap)
			ioutil.WriteFile("./data/stat.info.json", statInfoBytes, 0666)
			time.Sleep(time.Second)

		}
	}()

}

//初始化engine
func initEngine() bool {
	dirList := []string{
		"./data/contact_html",
		"./data/contact_json",
	}
	for _, dir := range dirList {
		_, err := os.Stat(dir)
		if err != nil && os.IsNotExist(err) {
			error := os.Mkdir(dir, os.ModePerm)
			if error != nil {
				fmt.Printf("mkdir err:%s\n", dir)
				return false
			} else {
				fmt.Printf("mkdir ok:%s\n", dir)
			}
		}
	}

	return true
}

func getPorxySeeds() []models.Request {

	var seedArr []models.Request
	pageNum := 10
	for i := 1; i <= pageNum; i++ {
		url := fmt.Sprintf("http://www.89ip.cn/index_%d.html", i)
		tempReq := models.Request{
			Url:    url,
			Method: "GET",
			ParserFunc: func(bytes []byte) models.ParserResult {
				return www_89ip_cn.ParseProxyItem(bytes)
			},
		}
		seedArr = append(seedArr, tempReq)
	}
	return seedArr
}

func getCotactSeeds() []models.Request {
	var seedArr []models.Request
	f, err := os.Open("./data/top_use_contact.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	rd := bufio.NewReader(f)

	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行
		if err != nil || io.EOF == err {
			break
		}
		contact := strings.Trim(line, "\n\r")
		tempParam := url.Values{}
		tempParam.Set("st", "contact")
		tempParam.Set("kw", contact)
		tempReq := models.Request{
			Url:       "http://faka.xinjipin.com/orderquery.htm",
			Method:    "POST",
			NeedProxy: true,
			Data:      tempParam,
			ParserFunc: func(bytes []byte) models.ParserResult {
				return xinjipin_com.ParseOrder(bytes, contact)
			},
		}
		seedArr = append(seedArr, tempReq)
	}
	return seedArr
}

func TestSingleProxy() {
	contact := "147258369"
	tempParam := url.Values{}
	tempParam.Set("st", "contact")
	tempParam.Set("kw", contact)
	curRequest := &models.Request{
		Url:    "http://faka.xinjipin.com/orderquery.htm",
		Method: "POST",
		Data:   tempParam,
		ParserFunc: func(bytes []byte) models.ParserResult {
			return xinjipin_com.ParseOrder(bytes, contact)
		},
	}

	item := &models.ProxyItem{
		Ip:   "127.0.0.1",
		Port: "8001",
		Type: "http",
	}
	body, _ := fetcher.FetchRequestByProxyItem(curRequest, item)
	//fmt.Println(string(body))
	ret := curRequest.ParserFunc(body)
	//fmt.Println(ret)
	jsonstr, _ := json.Marshal(ret.Items)
	fmt.Println(string(jsonstr))
}
