package http_srv

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func (s *Server) init_db() error {
	conn_string := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?autocommit=true", s.conf.DB.User, s.conf.DB.Passwd, s.conf.DB.Host, s.conf.DB.Port, s.conf.DB.DbName)

	var err error
	s.db, err = sql.Open("mysql", conn_string)
	if err != nil {
		return err
	}
	s.db.SetMaxOpenConns(100)

	return nil
}
