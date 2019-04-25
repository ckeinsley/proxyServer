package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

var debug = false

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Usage ./proxy <PORT>")
		panic(nil)
	}

	if len(args) > 1 {
		if strings.Compare("-d", args[1]) == 0 {
			debug = true
		}
	}

	service := ":" + args[0]
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	checkError(err)
	fmt.Printf("listening on port %s\n~\n", service)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn) // This is where we use multi-processing in golang
	}
}

// When the client makes a request the handleClient function processes it
func handleClient(conn net.Conn) {
	defer conn.Close()
	var recvData = make([]byte, 1024*4)
	conn.Read(recvData)

	if debug {
		fmt.Printf("Got request: %s", string(recvData))
	}

	request := parseHTTPRequest(conn, string(recvData))

	sendRequest(conn, request)
}

// parseHTTPRequest takes the client request and constructs
// our HttpRequest object from it
func parseHTTPRequest(conn net.Conn, data string) HTTPRequest {
	defer func() {
		if recover() != nil {
			conn.Write([]byte("400 Bad Request\n"))
			conn.Close()
		}
	}()

	httpRequest := CreateHTTPRequest(data)

	return httpRequest
}

// Once the client request has been processed our sendRequest function
// handles external communication
func sendRequest(conn net.Conn, request HTTPRequest) {
	defer func() {
		if recover() != nil {
			conn.Write([]byte("500: Unable to connect to remote server\n"))
			conn.Close()
		}
	}()

	tcpAddr, err := net.ResolveTCPAddr("tcp", request.Host+":"+request.Port)
	checkError(err)
	remoteConn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	connectionString := createConnectionString(request)

	remoteConn.Write([]byte(connectionString))

	// Not getting EOF?
	result, err := ioutil.ReadAll(remoteConn)
	checkError(err)

	conn.Write(result)
	conn.Close()
}

// createConnectionString takes our HTTPRequest object and generates
// the connection string used in sendRequest
func createConnectionString(request HTTPRequest) string {
	connectionString := ""
	connectionString += request.Method + " " + request.Route + " " + request.Version + "\n"
	connectionString += "Host: " + request.Host + "\n"
	connectionString += request.Connection + "\n"
	for _, header := range request.Headers {
		connectionString += header + "\n"
	}
	connectionString += "\n"
	return connectionString
}

// checkError handles panic escalation
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		panic(err)
	}
}
