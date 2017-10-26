package http_srv

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
)

var (
	ErrNoData = errors.New("table is nil")
)

const (
	SQL_FMT string = "select t2.%s ,t1.version from %s t2, rules_version t1 where t2.rules_id=t1.id and t1.valid=1"
	VERSION string = "version"
)

func (s *Server) update_core(redis_set string, table string, column string) error {
	//1 .select ci from mysql
	//2a.set ci to redis ci_new
	//2b.set mysql 2 ci
	//3 .find diff element(ci_new ci)
	//4 .delete diff from ci

	redis_set_new := redis_set + "_newversion"
	conn := s.pool.Get()
	defer conn.Close()
	conn.Do("DEL", redis_set_new)

	sql := fmt.Sprintf(SQL_FMT, column, table)
	log.Println(sql)
	rows, err := s.db.Query(sql)

	if err != nil {
		return err
	}
	defer rows.Close()

	row_count := 0
	var version string
	for rows.Next() {
		row_count++
		var value string
		if err := rows.Scan(&value, &version); err != nil {
			return err
		}
		conn.Send("SADD", redis_set, value)
		conn.Send("SADD", redis_set_new, value)
	}
	if row_count == 0 {
		return ErrNoData
	}
	conn.Send("SET", table+VERSION, version)
	conn.Do("")
	byte_values, _ := conn.Do("SDIFF", redis_set, redis_set_new)
	values, _ := redis.Strings(byte_values, nil)
	for _, diff := range values {
		err := conn.Send("SREM", redis_set, diff)
		if err != nil {
			log.Println(err.Error())
			return err
		}
	}
	conn.Do("")

	return nil
}
