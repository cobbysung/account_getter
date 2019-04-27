package engine

import (
	"account_getter/fetcher"
	"account_getter/logger"
	"account_getter/models"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"sync/atomic"
	"time"
)

type QueueEngine struct {
	Scheduler      Scheduler
	WorkerCount    int
	OrderDoneCount int64
}

type Scheduler interface {
	Submit(models.Request)
	WorkerReady(chan models.Request)
	ProxyWorkerReady(chan models.Request)
	Run()
}

func (e *QueueEngine) Run(seeds ...models.Request) {

	out := make(chan models.ParserResult)

	//scheduler run
	e.Scheduler.Run()

	//添加不使用proxy的worker
	for i := 0; i < e.WorkerCount; i++ {
		e.createWorker(out)
	}

	//添加任务
	for _, item := range seeds {
		e.Scheduler.Submit(item)
	}

	//处理结果
	for {
		parseResult := <-out
		for _, item := range parseResult.Items {
			logger.DebugLog.Printf("Got item:%s\n", item)
			//如果是proxyItem则加入proxyWorker
			if proxyItem, ok := item.(*models.ProxyItem); ok {
				if proxyItem.Status {
					e.createProxyWorker(proxyItem, out)
				} else {
					logger.DebugLog.Println("proxy status errr drop")
				}
			}
			//如果是order
			if orderItem, ok := item.(models.OrderDetail); ok {
				logger.DebugLog.Println("ORDER:", orderItem)
			}
		}

		for _, nR := range parseResult.Requests {
			e.Scheduler.Submit(nR)
		}
	}
}

func (e *QueueEngine) createWorker(out chan models.ParserResult) chan models.Request {
	in := make(chan models.Request)
	go func() {
		for {
			e.Scheduler.WorkerReady(in)
			r := <-in
			parserResult, err := e.doWork(r)
			if err != nil {
				logger.DebugLog.Println("doWork Error", err)
				continue
			}
			out <- parserResult
		}
	}()
	return in
}

func (e *QueueEngine) doWork(curRequest models.Request) (models.ParserResult, error) {
	body, err := fetcher.Fetch(curRequest.Url)
	logger.DebugLog.Printf("getting Url:%s\n", curRequest.Url)
	if err != nil {
		logger.DebugLog.Printf("Fetch error %s\n", curRequest.Url)
		return models.ParserResult{}, err
	}
	return curRequest.ParserFunc(body), nil
}

func (e *QueueEngine) createProxyWorker(item *models.ProxyItem, out chan models.ParserResult) chan models.Request {
	in := make(chan models.Request)
	errCnt := 0
	go func() {
		for {
			e.Scheduler.ProxyWorkerReady(in)
			r := <-in
			parserResult, err := e.doProxyWork(item, r)
			logger.DebugLog.Println("doProxyWork", parserResult.Items)

			if err != nil {
				logger.DebugLog.Println("doProxyWork Error", err)
				errCnt++
				continue
			}

			//保存数据
			contact := r.Data.Get("kw")
			if contact != "" {
				//计数
				atomic.AddInt64(&e.OrderDoneCount, 1)
				logger.DebugLog.Printf("SAVE_JSON:%s\n", contact)
				jsonstr, _ := json.Marshal(parserResult.Items)
				filename := fmt.Sprintf("./data/contact_json/%s.json", Md5Sum(contact))
				ioutil.WriteFile(filename, jsonstr, 0666)
			}
			out <- parserResult

			time.Sleep(time.Second * 5)
			if errCnt == 2 {
				//失败次数大于2，则退出
				return
			}
		}
	}()
	return in
}

func (e *QueueEngine) doProxyWork(item *models.ProxyItem, curRequest models.Request) (models.ParserResult, error) {
	//抓取数据
	body, err := fetcher.FetchRequestByProxyItem(&curRequest, item)

	//如果是有联系方式保存下载
	contact := curRequest.Data.Get("kw")
	if contact != "" {
		logger.DebugLog.Printf("SAVE_HTML:%s\n", contact)
		filename := fmt.Sprintf("./data/contact_html/%s.html", Md5Sum(contact))
		ioutil.WriteFile(filename, body, 0666)
	}

	logger.DebugLog.Printf("getting Proxy Url:%s\n", curRequest.Url)
	if err != nil {
		logger.DebugLog.Printf("Fetch Proxy error %s\n", curRequest.Url)
		//如果失败则重试
		e.Scheduler.Submit(curRequest)
		return models.ParserResult{}, err
	}
	return curRequest.ParserFunc(body), nil
}

//获取md5值
func Md5Sum(str string) string {
	w := md5.New()
	io.WriteString(w, str)                   //将str写入到w中
	md5str2 := fmt.Sprintf("%x", w.Sum(nil)) //w.Sum(nil)将w的hash转成[]byte格式
	return md5str2
}
