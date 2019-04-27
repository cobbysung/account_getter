package scheduler

import (
	"account_getter/logger"
	"account_getter/models"
	"time"
)

type QueueScheduler struct {
	requestChan      chan models.Request
	workerChan       chan chan models.Request
	proxyRequestChan chan models.Request
	proxyWorkerChan  chan chan models.Request
}

func (s *QueueScheduler) Submit(r models.Request) {
	go func() {
		if r.NeedProxy {
			s.proxyRequestChan <- r
		} else {
			s.requestChan <- r
		}
	}()
}

func (s *QueueScheduler) WorkerReady(wChan chan models.Request) {
	s.workerChan <- wChan
}

func (s *QueueScheduler) ProxyWorkerReady(wChan chan models.Request) {
	s.proxyWorkerChan <- wChan
}

func (s *QueueScheduler) Run() {
	s.runType(true)
	s.runType(false)
}

func (s *QueueScheduler) runType(needProxy bool) {
	//构造chan
	var requestChan chan models.Request
	var workerChan chan chan models.Request
	var keyStr string
	if needProxy {
		s.requestChan = make(chan models.Request)
		s.workerChan = make(chan chan models.Request)
		requestChan = s.requestChan
		workerChan = s.workerChan
		keyStr = "PROXY"
	} else {
		s.proxyRequestChan = make(chan models.Request)
		s.proxyWorkerChan = make(chan chan models.Request)
		requestChan = s.proxyRequestChan
		workerChan = s.proxyWorkerChan
		keyStr = "NO_PROXY"
	}

	go func() {
		var requestQueue []models.Request
		var workerQueue []chan models.Request

		tick := time.Tick(time.Second * 2)
		tick2 := time.Tick(time.Millisecond * 100)

		for {
			var curReq models.Request
			var curWorker chan models.Request

			if len(requestQueue) > 0 && len(workerQueue) > 0 {
				curReq = requestQueue[0]
				curWorker = workerQueue[0]
			}

			select {
			case r := <-requestChan:
				<-tick2 //控制放的速度
				requestQueue = append(requestQueue, r)
			case w := <-workerChan:
				workerQueue = append(workerQueue, w)
			case curWorker <- curReq:
				requestQueue = requestQueue[1:]
				workerQueue = workerQueue[1:]
			case <-tick:
				logger.DebugLog.Printf("%s reqlen:%d \t worklen:%d\n",
					keyStr, len(requestQueue), len(workerQueue))

			}
		}

	}()
}
