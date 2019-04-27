package models

import "net/url"

type Request struct {
	Url        string
	Method     string
	Data       url.Values
	NeedProxy  bool
	ParserFunc func([]byte) ParserResult
}
type ParserResult struct {
	Requests []Request
	Items    []interface{}
}

func NilParser([]byte) ParserResult {
	return ParserResult{}
}
