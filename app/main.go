package main

import (
	//"bytes"
	//"compress/gzip"
	"flag"
	"fmt"
	"net"
	"os"
	//"path/filepath"
	//s "strings"
)

var fileDir string

func main() {
	setupFlags()
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	fmt.Println("Listening on port 4221")
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn)
	}
}

func setupFlags() {
	flag.StringVar(&fileDir, "directory", ".", "Directory to serve files from")
	flag.Parse()
}
func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Accepting connection from ", conn.RemoteAddr(), "\n")
	for {
		buffer := make([]byte, 1024)
		c, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from connection: ", err.Error())
			break
		}
		request := string(buffer[:c])
		req, err := parseRequest(request)
		fmt.Printf("Received request: %+v \n", req)
		if err != nil {
			break
		}
		res := Response{
			headers: make(map[string]string),
		}
		handleRequest(&req, &res)
		response := serializeResponse(&res)

		_, err = conn.Write(response)
		if err != nil {
			fmt.Println("Error writing ", err.Error())
			os.Exit(1)
		}
		//if closeConnection {
		//	break
		//}
	}

}

func handleRequest(req *Request, res *Response) {
	Router(req, res)
	finalizeResponse(req, res)
}
