package http_srv

import (
	"log"
	"net/http"
	"net/http/httputil"
)

func dump_requst(r *http.Request) {
	v, e := httputil.DumpRequest(r, true)
	if e != nil {
		log.Println(e.Error())
		return
	}
	log.Println(string(v))
}
