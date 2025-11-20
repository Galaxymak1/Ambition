package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	s "strings"
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
		//conn.Close()
	}
}

func setupFlags() {
	flag.StringVar(&fileDir, "directory", ".", "Directory to serve files from")
	flag.Parse()
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
	requestLine, headers, body := parseRequest(req)
	response := handleRoutes(requestLine, headers, body)
	return response
}

func handleRoutes(requestLine string, headers []string, body string) string {
	method, urlParts, _, err := parseRequestLine(requestLine)
	if err != nil {
		fmt.Println("Error parsing request line: ", err.Error())
		return BAD_REQUEST + "\r\n"
	}

	lenUrl := len(urlParts)
	if lenUrl == 0 {
		return OK + "\r\n"
	}
	baseRoute := urlParts[0]
	switch baseRoute {
	case "echo":
		var acceptEncoding string
		for _, header := range headers {
			if s.Contains(header, "Accept-Encoding") {
				acceptEncoding = s.TrimSpace(s.Split(header, ": ")[1])
			}
		}

		if lenUrl > 1 {
			if acceptEncoding == "gzip" {
				return OK + "Content-Type: text/plain\r\nContent-Encoding: gzip\r\nContent-Length: " + strconv.Itoa(len(urlParts[1])) + "\r\n\r\n" + urlParts[1]
			} else {
				return OK + "Content-Type: text/plain\r\nContent-Length: " + strconv.Itoa(len(urlParts[1])) + "\r\n\r\n" + urlParts[1]
			}
		} else {
			return BAD_REQUEST + "\r\n"
		}

	case "user-agent":
		var userAgent string
		for _, header := range headers {
			if s.Contains(header, "User-Agent") {
				userAgent = s.TrimSpace(s.Split(header, ": ")[1])
			}
		}
		if userAgent != "" {
			return OK + "Content-Type: text/plain\r\nContent-Length: " + strconv.Itoa(len(userAgent)) + "\r\n\r\n" + userAgent
		}
	case "files":
		if lenUrl != 2 {
			return BAD_REQUEST + "\r\n"
		}
		switch method {
		case "GET":
			filename := urlParts[1]
			path := filepath.Join(fileDir, "/", filename)
			file, err := os.ReadFile(path)
			if err != nil {
				fmt.Println("Error reading file: ", err.Error())
				return NOT_FOUND + "\r\n"
			}
			return OK + "Content-Type: application/octet-stream\r\nContent-Length: " + strconv.Itoa(len(file)) + "\r\n\r\n" + string(file)
		case "POST":
			filename := urlParts[1]
			file, err := os.Create(filepath.Join(fileDir, "/", filename))
			if err != nil {
				fmt.Println("Error creating file: ", err.Error())
				return NOT_FOUND + "\r\n"
			}
			l, err := file.WriteString(body)
			if err != nil {
				fmt.Println("Error writing to file: ", err.Error())
				return NOT_FOUND + "\r\n"
			}
			fmt.Println(l, "bytes written successfully")
			return CREATED + "\r\n"
		}

	default:
		return NOT_FOUND + "\r\n"

	}
	return NOT_FOUND + "\r\n"
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
	return requestLine, headers, s.TrimPrefix(body, "\n")
}

func parseRequestLine(requestLine string) (string, []string, string, error) {
	parts := s.Split(requestLine, " ")
	if len(parts) != 3 {
		return "", []string{""}, "", errors.New("invalid request line")
	}
	method, urlParts, protocol := parts[0], deleteEmpty(s.Split(parts[1], "/")), parts[2]
	return method, urlParts, protocol, nil
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
