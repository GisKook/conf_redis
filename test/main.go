package main

import (
	"fmt"
	"github.com/giskook/conf_redis/conf"
	"github.com/giskook/conf_redis/http_srv"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cfg, err := conf.ReadConfig("./conf.json")
	if err != nil {
		log.Println(err.Error())
	}

	s := http_srv.NewServer(cfg)
	err = s.Init()
	if err != nil {
		log.Println(err.Error())
		s.Close()
		return
	}

	s.Handle()
	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)
}
