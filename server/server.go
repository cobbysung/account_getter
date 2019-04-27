package server

import (
	"account_getter/fetcher"
	"account_getter/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"
)

//扩展订单信息
type OrderExtend struct {
	models.OrderDetail
	From string
	Days int64
}

//排序算法
type OrderSlice []OrderExtend

func (s OrderSlice) Len() int      { return len(s) }
func (s OrderSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s OrderSlice) Less(i, j int) bool {
	iTime, _ := TimeStr2Stamp(s[i].Time)
	jTime, _ := TimeStr2Stamp(s[j].Time)
	return iTime > jTime
}

func RunHttpServ(port string) {

	http.HandleFunc("/home", func(writer http.ResponseWriter, request *http.Request) {
		html, _ := ioutil.ReadFile("./data/home.html")
		writer.Write(html)
	})
	http.HandleFunc("/get_list", func(writer http.ResponseWriter, request *http.Request) {
		orderList := getAllOrders()
		jsonBytes, _ := json.Marshal(orderList)
		writer.Write(jsonBytes)
	})

	http.HandleFunc("/get_stat", func(writer http.ResponseWriter, request *http.Request) {
		statInfoBytes, _ := ioutil.ReadFile("./data/stat.info.json")
		writer.Write(statInfoBytes)
	})

	http.HandleFunc("/get_order", func(writer http.ResponseWriter, request *http.Request) {
		request.ParseForm()
		orderId := request.Form.Get("order_id")
		contents := getOrderDetail(orderId)
		var infoMap map[string]interface{}
		json.Unmarshal(contents, &infoMap)
		output := fmt.Sprintf("%+v", infoMap)
		writer.Write([]byte(output))
	})
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("HTTP_SERVER_ERROR %v", err)
	}

}

func getOrderDetail(orderId string) []byte {
	api := fmt.Sprintf("http://faka.xinjipin.com/checkgoods.htm?orderid=%s&_=%d123", orderId, time.Now().Unix())
	contents, _ := fetcher.Fetch(api)
	return contents
}

func getAllOrders() OrderSlice {
	//获取所有文件
	dir := "./data/contact_json/"
	fileList, _ := getDirList(dir)

	//获取所有数据
	var orderFromMap = map[string]string{}
	var allOrderList OrderSlice
	var nowStamp = time.Now().Unix()
	for _, filename := range fileList {
		data, _ := ioutil.ReadFile(filename)
		//info, _ := os.Stat(filename)
		//fmt.Print(filename, info.ModTime().Unix())
		var orderList []models.OrderDetail
		json.Unmarshal(data, &orderList)
		for _, item := range orderList {
			if item.OrderMoney != item.PayMoney {
				continue
			}
			timeStamp, _ := TimeStr2Stamp(item.Time)
			tdays, from := GetOrderDaysType(item)
			//记录来源
			orderFromMap[item.Id] = from
			//过滤有效期
			if timeStamp+tdays*24*3600 < nowStamp {
				continue
			}
			newItem := OrderExtend{
				From:        from,
				Days:        tdays,
				OrderDetail: item,
			}
			allOrderList = append(allOrderList, newItem)
		}
	}
	//排序
	sort.Stable(allOrderList)

	return allOrderList

}

//获取文件列表
func getDirList(dirpath string) ([]string, error) {
	var dir_list []string
	dir_err := filepath.Walk(dirpath,
		func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() == false {
				dir_list = append(dir_list, path)
				return nil
			}

			return nil
		})
	return dir_list, dir_err
}

//转时间戳
func TimeStr2Stamp(timeStr string) (int64, error) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return 0, err
	}
	timeTmp, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, loc)
	if err != nil {
		return 0, err
	}
	return timeTmp.Unix(), nil
}

//获取账号有效天数
func GetOrderDaysType(item models.OrderDetail) (int64, string) {
	type GoodsList []map[string]float64

	allInfoMap := map[string]GoodsList{
		"dupan": {
			{"money": 1.80, "days": 3},
			{"money": 2.50, "days": 6},
			{"money": 3.50, "days": 6},
			{"money": 6.50, "days": 30},
			{"money": 27, "days": 180},
			{"money": 45, "days": 365},
		},
		"thunder": {
			{"money": 2.88, "days": 2},
			{"money": 4.5, "days": 18},
			{"money": 4.9, "days": 20},
			{"money": 6.9, "days": 21},
			{"money": 7.5, "days": 0}, //激活码
			{"money": 7.9, "days": 30},
		},
		"youku": {
			{"money": 6.80, "days": 52},
			{"money": 7.80, "days": 60},
			{"money": 55, "days": 365},
		},
		"ximalaya": {
			{"money": 3.9, "days": 30},
			{"money": 9.8, "days": 90},
		},
		"wenku": {
			{"money": 2.9, "days": 0}, //文库下载器
		},
		"others": {
			{"money": 999, "days": -1},
			{"money": 11.80, "days": 0}, //激活码
			{"money": 19.80, "days": 0}, //激活码
		},
	}

	var days float64 = -100
	var typeStr string = "unkown"
	for siteName, goodsList := range allInfoMap {
		for _, goods := range goodsList {
			if goods["money"] == item.PayMoney {
				days = math.Max(days, goods["days"])
				typeStr = siteName
				break
			}
		}
	}
	if days == -100 {
		days = 180
	}
	return int64(days), typeStr
}
