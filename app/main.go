package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
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
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		return
	}

	request := string(buffer[:c])
	response := handleRequest(request)

	data := []byte(response)
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("Error writing ", err.Error())
		os.Exit(1)
	}
}
func handleRequest(req string) string {
	requestLine, headers, _ := parseRequest(req)
	response := handleRoutes(requestLine, headers)
	return response
}

func handleRoutes(requestLine string, headers []string) string {
	_, urlParts, _, err := parseRequestLine(requestLine)
	if err != nil {
		fmt.Println("Error parsing request line: ", err.Error())
		return "HTTP/1.1 400 Bad Request\r\n\r\n"
	}
	var userAgent string
	for _, header := range headers {
		if s.Contains(header, "User-Agent") {
			userAgent = s.Split(header, ": ")[1]
		}
	}
	lenUrl := len(urlParts)
	if lenUrl == 0 {
		return "HTTP/1.1 200 OK\r\n\r\n"
	} else if urlParts[0] == "echo" && lenUrl == 2 {
		return "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + strconv.Itoa(len(urlParts[1])) + "\r\n\r\n" + urlParts[1]
	} else if urlParts[0] == "user-agent" {
		if userAgent != "" {
			println(userAgent)
			return "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + strconv.Itoa(len(userAgent)) + "\r\n\r\n" + userAgent
		}
	} else {
		return "HTTP/1.1 404 Not Found\r\n\r\n"
	}
	return ""
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

func parseRequestLine(requestLine string) (string, []string, string, error) {
	parts := s.Split(requestLine, " ")
	if len(parts) != 3 {
		return "", []string{""}, "", errors.New("invalid request line")
	}
	method, urlParts, protocol := parts[0], delete_empty(s.Split(parts[1], "/")), parts[2]
	return method, urlParts, protocol, nil
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
