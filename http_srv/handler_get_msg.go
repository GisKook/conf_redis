package http_srv

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"net/http"
)

const (
	MESSAGE_KEY string = "message"
)

type ReadResponse struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
	Msg  string `json:"msg"`
}

func EncodeReadResponse(code int, errmsg string, msg string) string {
	response := &ReadResponse{
		Code: code,
		Desc: errmsg,
		Msg:  msg,
	}

	resp, _ := json.Marshal(response)

	return string(resp)
}

func (s *Server) handler_get_msg(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if x := recover(); x != nil {
			log.Println("crash")
			fmt.Fprint(w, EncodeResponse(HTTP_RESP_INTERNAL_ERR))
		}
	}()
	dump_requst(r)

	conn := s.pool.Get()
	defer conn.Close()
	v, err := conn.Do("GET", MESSAGE_KEY)
	if err != nil {
		fmt.Fprint(w, EncodeReadResponse(int(HTTP_RESP_ERR), err.Error(), ""))
		return

	}
	msg, _ := redis.String(v, nil)
	fmt.Fprint(w, EncodeReadResponse(int(HTTP_RESP_SUCCESS), HTTP_RESP_SUCCESS.Desc(), msg))

}
