package http_srv

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
)

type ResponseCode int

const (
	HTTP_RESP_SUCCESS      ResponseCode = 0
	HTTP_RESP_ERR          ResponseCode = 1
	HTTP_RESP_EXCLUSIVE    ResponseCode = 2
	HTTP_RESP_INTERNAL_ERR ResponseCode = 3
)

var HTTP_RESP_DESC []string = []string{
	"成功",
	"失败",
	"其他用户正在操作",
	"内部错误",
}

func (c ResponseCode) Desc() string {
	return HTTP_RESP_DESC[c]
}

type Response struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
}

func EncodeResponse(code ResponseCode) string {
	response := &Response{
		Code: int(code),
		Desc: code.Desc(),
	}

	resp, _ := json.Marshal(response)

	return string(resp)
}

func dump_requst(r *http.Request) {
	v, e := httputil.DumpRequest(r, true)
	if e != nil {
		log.Println(e.Error())
		return
	}
	log.Println(string(v))
}
