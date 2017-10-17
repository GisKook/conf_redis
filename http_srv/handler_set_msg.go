package http_srv

import (
	"fmt"
	"log"
	"net/http"
)

const (
	MESSAGE string = "msg"
)

func (s *Server) handler_set_msg(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if x := recover(); x != nil {
			log.Println("crash")
			fmt.Fprint(w, EncodeResponse(HTTP_RESP_INTERNAL_ERR))
		}
	}()
	dump_requst(r)

	r.ParseForm()
	message := r.Form.Get(MESSAGE)
	if message == "" {
		fmt.Fprint(w, EncodeResponse(HTTP_RESP_ERR))
		return
	}
	conn := s.pool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", MESSAGE_KEY, message)
	if err != nil {
		fmt.Fprint(w, EncodeResponse(HTTP_RESP_ERR))
		return

	}
	fmt.Fprint(w, EncodeResponse(HTTP_RESP_SUCCESS))
}
