package xinjipin_com

import (
	"account_getter/logger"
	"account_getter/models"
	"bytes"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ParseOrder(contents []byte, contact string) models.ParserResult {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(contents))
	result := models.ParserResult{}
	doc.Find("div.search_list tr").Each(func(i int, s *goquery.Selection) {
		tds := s.Find("td")
		if tds.Length() == 0 {
			return
		}
		var tdsArr []string
		tds.Each(func(i int, s *goquery.Selection) {
			tdsArr = append(tdsArr, strings.Trim(s.Text(), " \n\r\t"))
		})
		money1, _ := strconv.ParseFloat(tdsArr[3], 64)
		money2, _ := strconv.ParseFloat(tdsArr[4], 64)
		tmpOrder := models.OrderDetail{
			Id:         tdsArr[1],
			Contact:    contact,
			Time:       tdsArr[0],
			PayMethod:  tdsArr[2],
			OrderMoney: money1,
			PayMoney:   money2,
		}
		result.Items = append(result.Items, tmpOrder)
		logger.DebugLog.Println(tmpOrder)
	})
	return result
}
