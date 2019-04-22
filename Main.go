package main

import (
	"fmt"
	"io/ioutil"
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

	fmt.Printf("Got from browser: \"\n%s\"\n", string(recvData))
	httpRequest := CreateHTTPRequest(string(recvData)) // TODO handle bad requests
	result := sendRequest(httpRequest)
	conn.Write(result)
}

func sendRequest(request HTTPRequest) []byte {
	tcpAddr, err := net.ResolveTCPAddr("tcp", request.Host+":"+request.Port)
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	connectionString := createConnectionString(request)
	conn.Write([]byte(connectionString))

	result, err := ioutil.ReadAll(conn)
	checkError(err)
	return result
}

func createConnectionString(request HTTPRequest) string {
	connectionString := ""
	connectionString += request.Method + " " + request.Route + " " + request.Version + "\n"
	connectionString += "Host: " + request.Host + "\n"
	connectionString += request.Connection + "\n"
	for _, header := range request.Headers {
		connectionString += header + "\n"
	}
	return connectionString
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
