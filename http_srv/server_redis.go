package http_srv

import (
	"github.com/garyburd/redigo/redis"
)

func (s *Server) cnt_redis() (redis.Conn, error) {
	c, err := redis.Dial("tcp", s.conf.Redis.Addr)
	if err != nil {
		return nil, err
	}
	if _, err := c.Do("AUTH", s.conf.Redis.Passwd); err != nil {
		c.Close()
		return nil, err
	}
	return c, nil
}

func (s *Server) init_redis() {
	s.pool = redis.NewPool(s.cnt_redis, s.conf.Redis.MaxIdle)
}
