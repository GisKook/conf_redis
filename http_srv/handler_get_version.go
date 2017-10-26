package http_srv

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"net/http"
)

const (
	TYPE string = "type"
)

type GetVersionResponse struct {
	Code    int    `json:"code"`
	Desc    string `json:"desc"`
	Version string `json:"version"`
	Count   string `json:"count"`
}

func EncodeGetVersionResponse(code int, errmsg string, version string, count string) string {
	response := &GetVersionResponse{
		Code:    code,
		Desc:    errmsg,
		Version: version,
		Count:   count,
	}

	resp, _ := json.Marshal(response)

	return string(resp)
}

func (s *Server) handler_get_version(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if x := recover(); x != nil {
			log.Println("crash")
			fmt.Fprint(w, EncodeResponse(HTTP_RESP_INTERNAL_ERR))
		}
	}()
	dump_requst(r)

	t := r.Form.Get(TYPE)
	if t == "" {
		fmt.Fprint(w, EncodeGetVersionResponse(int(HTTP_RESP_ERR), "缺少参数type type可选值为 ci,black_list,local_number,county_number", "", ""))
		return
	}

	conn := s.pool.Get()
	defer conn.Close()
	v, err := conn.Do("GET", t+VERSION)
	if err != nil {
		fmt.Fprint(w, EncodeGetVersionResponse(int(HTTP_RESP_ERR), err.Error(), "", ""))
		return
	}
	version, _ := redis.String(v, nil)

	v, err = conn.Do("SCARD", t)
	if err != nil {
		fmt.Fprint(w, EncodeGetVersionResponse(int(HTTP_RESP_ERR), err.Error(), version, ""))
		return
	}
	count, _ := redis.String(v, nil)

	fmt.Fprint(w, EncodeGetVersionResponse(int(HTTP_RESP_SUCCESS), HTTP_RESP_SUCCESS.Desc(), version, count))

}
