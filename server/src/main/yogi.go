package main

import (
	"os"
	"log"
	"flag"
	"net"
	"runtime"
)

// some global service instances
var logger *log.Logger
var sysinfo SysChan

func main() {
	runtime.GOMAXPROCS(4)
	var logfile_name string
	var addr string
	flag.StringVar(&logfile_name, "l", "log.txt", "log file path")
	flag.StringVar(&addr, "addr", "127.0.0.1:8088", "bind address")
	flag.Parse()

	logfile, err := os.OpenFile(logfile_name, os.O_WRONLY | os.O_TRUNC | os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logger = log.New(logfile, "", log.Ldate|log.Ltime|log.Llongfile)
	logger.Println("run begin")

	go db_thread()

	sysinfo = make(SysChan)
	go sysinfo.Run()

    go theUserMgr.Run()

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Println(err.Error())
			continue
		}
		go NewConn(conn)
	}
}
