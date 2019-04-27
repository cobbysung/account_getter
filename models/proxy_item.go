package models

import (
	"encoding/json"
)

type ProxyItem struct {
	Ip         string `json:"ip"`
	Port       string `json:"port"`
	Type       string `json:"type"`
	Status     bool   `json:"status"`
	UpdateTime string `json:"update_time"`
	Source     string `json:"source"`
}

func (p ProxyItem) String() string {
	data, _ := json.Marshal(p)
	return string(data)
}
