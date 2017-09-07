package conf

import (
	"encoding/json"
	"os"
)

type Http struct {
	Addr    string
	TimeOut int
}

type DB struct {
	Host   string
	Port   string
	User   string
	Passwd string
	DbName string
}

type Redis struct {
	Addr        string
	MaxIdle     int
	IdleTimeOut int
	Passwd      string
}

type Conf struct {
	Http  *Http
	DB    *DB
	Redis *Redis
}

func ReadConfig(confpath string) (*Conf, error) {
	file, _ := os.Open(confpath)
	decoder := json.NewDecoder(file)
	config := Conf{}
	err := decoder.Decode(&config)

	return &config, err
}
