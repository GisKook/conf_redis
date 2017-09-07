package http_srv

import (
	//"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"net/http"
)

// select t2.ci from ci t2, rules_version t1 where t2.rules_id=t1.id and t1.valid=1;
func (s *Server) update_ci() error {
	//1 .select ci from mysql
	//2a.set ci to redis ci_new
	//2b.set mysql 2 ci
	//3 .find diff element(ci_new ci)
	//4 .delete diff from ci
	conn := s.pool.Get()
	defer conn.Close()
	conn.Do("DEL", "ci_new")

	rows, err := s.db.Query("select t2.ci from ci t2, rules_version t1 where t2.rules_id=t1.id and t1.valid=1")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			return err
		}
		conn.Send("SADD", "ci", value)
		conn.Send("SADD", "ci_new", value)
	}
	conn.Do("")
	byte_values, _ := conn.Do("SDIFF", "ci", "ci_new")
	values, _ := redis.Strings(byte_values, nil)
	log.Println(values)
	for _, diff := range values {
		err := conn.Send("SREM", "ci", diff)
		log.Println(diff)
		if err != nil {
			log.Println(err.Error())
		}
	}
	conn.Do("")

	return nil
}

func (s *Server) handler_update_ci(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if x := recover(); x != nil {
			log.Println("crash")
			//fmt.Fprint(w, "")
		}
	}()
	dump_requst(r)
	if !s.flag_ci.lock() {
		log.Println("lock")
		return
	}
	if err := s.update_ci(); err != nil {
		log.Println(err.Error())
		log.Println("update error")
	} else {
		log.Println("update")
	}

	s.flag_ci.unlock()
}
