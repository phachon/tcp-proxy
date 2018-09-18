package main

import (
	"net"
	"log"
	"github.com/spf13/pflag"
	"io"
	"os"
	"fmt"
)

// a tcp proxy

// version
const (

	Version = "0.0.1"

	Author = "phachon@163.com"

	Link = "https://github.com/phachon/tcp-proxy"
)

var (
	local  = pflag.String("local", ":40001", "local listen host:port")
	remote  = pflag.String("remote", "", "remote listen host:port")
)

func init() {
	initFlag()
	initPoster()
}

// init poster
func initPoster() {
	logo := `
 _____ ____ ____    ____  ____   _____  ____   __
|_   _/ ___|  _ \  |  _ \|  _ \ / _ \ \/ /\ \ / /
  | || |   | |_) | | |_) | |_) | | | \  /  \ V / 
  | || |___|  __/  |  __/|  _ <| |_| /  \   | |  
  |_| \____|_|     |_|   |_| \_\\___/_/\_\  |_|

` +
		" Author: phachon@163.com \r\n" +
		" Version: " + Version + "\r\n" +
		" Link: " + Link + "\r\n" +
		"-----------------------------"
	fmt.Println(logo)
}

// init flag
func initFlag()  {
	pflag.Parse()
	if *remote == "" {
		log.Println("args remote is not empty!")
		os.Exit(100)
	}
}

func main()  {

	listen, err := net.Listen("tcp", *local)
	if err != nil {
		log.Printf("tcp proxy listen error, %s", err.Error())
		os.Exit(100)
	}
	defer listen.Close()

	log.Printf("tcp proxy start listen %s", *local)

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("tcp proxy listener accept error, %s", err.Error())
			break
		}
		go func(conn *net.Conn) {
			defer func() {
				(*conn).Close()
				e := recover()
				if e != nil {
					log.Printf("tcp proxy handle conn crash: %v", e)
				}
			}()
			err := handleConn(conn)
			if err != nil {
				log.Printf("tcp proxy handle conn error: %s", err.Error())
			}
		}(&conn)
	}
}

// handle conn
func handleConn(localConn *net.Conn) (err error) {

	log.Printf("tcp proxy receive conn %s", (*localConn).LocalAddr())

	remoteConn, err := net.Dial("tcp", *remote)
	if err != nil {
		return
	}
	defer remoteConn.Close()

	// I/O copy
	go func() {
		io.Copy(*localConn, remoteConn)
	}()
	io.Copy(remoteConn, *localConn)
	return err
}
