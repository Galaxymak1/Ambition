package main

import (
	"fmt"
	"net"
	"os"
	s "strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
//var _ = net.Listen
//var _ = os.Exit

func main() {

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
		//conn.Close()
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Accepting connection from ", conn.RemoteAddr(), "\n")
	buffer := make([]byte, 1024)
	c, err := conn.Read(buffer)
	request := string(buffer[:c])
	response := handleRequest(request)

	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
	}
	data := []byte(response)
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("Error writing ", err.Error())
		os.Exit(1)
	}
}
func handleRequest(req string) string {
	requestLine, _, _ := parseRequest(req)
	passStage := s.Split(requestLine, " ")[1]
	test := s.Split(passStage, "/")
	test = delete_empty(test)
	fmt.Printf("Test: %#v\n", passStage)
	fmt.Printf("Test2: %#v\n", test)
	fmt.Println("len", len(test))
	if len(test) > 0 {
		return "HTTP/1.1 404 Not Found\r\n\r\n"
	} else {
		return "HTTP/1.1 200 OK\r\n\r\n"
	}
}

func parseRequest(req string) (string, []string, string) {
	//fmt.Println("Received request :", req)
	reqParts := s.Split(req, "\n\r")
	head := s.Split(reqParts[0], "\n")
	var body string
	if len(reqParts) > 1 {
		body = reqParts[1]
	}
	requestLine := head[0]
	headers := head[1:]
	fmt.Printf("REQUEST LINE : %#v \n", requestLine)
	fmt.Printf("HEADERS LINE : %#v \n", headers)
	if len(body) > 0 {
		fmt.Printf("BODY  LINE : %#v \n", body)
	}
	return requestLine, headers, body
}

func delete_empty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
