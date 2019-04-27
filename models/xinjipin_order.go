package models

import "encoding/json"

type OrderDetail struct {
	Id         string
	Contact    string
	Time       string
	PayMethod  string
	OrderMoney float64
	PayMoney   float64
}

func (p OrderDetail) String() string {
	data, _ := json.Marshal(p)
	return string(data)
}
