package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Usage ./proxy <PORT>")
		panic(nil)
	}

	service := ":" + args[0]
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	checkError(err)
	fmt.Printf("listening on port %s\n", service)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	var recvData = make([]byte, 1024)
	conn.Read(recvData)
	fmt.Printf("\"\n%s\"\n", string(recvData))
	conn.Write([]byte("welcome to proxy"))
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
