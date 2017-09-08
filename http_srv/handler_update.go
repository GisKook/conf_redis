package http_srv

import (
	"fmt"
	"log"
	"net/http"
)

func (s *Server) handler_common(w http.ResponseWriter, r *http.Request, redis_set string, table string, column string) {

	defer func() {
		if x := recover(); x != nil {
			log.Println("crash")
			fmt.Fprint(w, EncodeResponse(HTTP_RESP_INTERNAL_ERR))
		}
	}()
	dump_requst(r)
	if !s.flag_ci.lock() {
		fmt.Fprint(w, EncodeResponse(HTTP_RESP_EXCLUSIVE))
		return
	}
	if err := s.update_core(redis_set, table, column); err != nil {
		log.Println(err.Error())
		fmt.Fprint(w, EncodeResponse(HTTP_RESP_EXCLUSIVE))
		return
	}
	s.flag_ci.unlock()

	fmt.Fprint(w, EncodeResponse(HTTP_RESP_SUCCESS))
}

func (s *Server) handler_update_ci(w http.ResponseWriter, r *http.Request) {
	s.handler_common(w, r, "ci", "ci", "ci")
}

func (s *Server) handler_update_black_list(w http.ResponseWriter, r *http.Request) {
	s.handler_common(w, r, "black_list", "black_list", "phone")
}

func (s *Server) handler_update_county_number(w http.ResponseWriter, r *http.Request) {
	s.handler_common(w, r, "county_number", "county_number", "county_number")
}

func (s *Server) handler_update_local_number(w http.ResponseWriter, r *http.Request) {
	s.handler_common(w, r, "local_number", "local_number", "local_number")
}

func (s *Server) handler_update_white_list(w http.ResponseWriter, r *http.Request) {
	s.handler_common(w, r, "white_list", "white_list", "phone")
}

func (s *Server) handler_update_unsub_number(w http.ResponseWriter, r *http.Request) {
	s.handler_common(w, r, "unsub_number", "unsub_number", "unsub_number")
}
