package http_srv

import (
	"database/sql"
	"github.com/garyburd/redigo/redis"
	"github.com/giskook/conf_redis/conf"
	"log"
	"net/http"
	"sync/atomic"
)

type lock struct {
	key int32
}

type Server struct {
	conf *conf.Conf
	db   *sql.DB
	pool *redis.Pool

	flag_ci            lock
	flag_local_number  lock
	flag_county_number lock
	flag_black_list    lock
}

func NewServer(conf *conf.Conf) *Server {
	return &Server{
		conf: conf,
	}
}

func (l *lock) lock() bool {
	if !l.islocked() {
		atomic.StoreInt32(&l.key, 1)
		return true
	}

	return false
}

func (l *lock) unlock() {
	atomic.StoreInt32(&l.key, 0)
}

func (l *lock) islocked() bool {
	return atomic.LoadInt32(&l.key) == 1
}

func (s *Server) Init() error {
	err := s.init_db()
	if err != nil {
		return err
	}
	log.Println("<INFO> db connect success")
	s.init_redis()
	log.Println("<INFO> redis connect success")
	return nil
}

func (s *Server) Handle() {
	http.HandleFunc("/knet2sp/update_ci", s.handler_update_ci)
	http.ListenAndServe(s.conf.Http.Addr, nil)
}

func (s *Server) Close() {

}
