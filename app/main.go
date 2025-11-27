package main

import (
	"bytes"
	"compress/gzip"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
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
		if err != nil {
			break
		}
		response, closeConnection := handleRequest(req)

		data := []byte(response)
		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("Error writing ", err.Error())
			os.Exit(1)
		}
		if closeConnection {
			break
		}
	}

}
func handleRequest(req Request) (string, bool) {
	closeConnection := false
	for _, header := range req.headers {
		h := s.ToLower(s.TrimSpace(header))
		if s.HasPrefix(h, "connection:") && s.Contains(h, "close") {
			closeConnection = true
		}
	}

	response := handleRoutes(req, closeConnection)
	return response, closeConnection
}

func handleRoutes(req Request, closeConnection bool) string {
	res := Response{"", []string{}, ""}
	if closeConnection {
		res.addHeader("Connection: close")
	}

	lenUrl := len(req.url)
	if lenUrl == 0 {
		res.addStatus(OK)
		return res.constructResponse()
	}
	baseRoute := req.url[0]
	switch baseRoute {
	case "echo":
		var acceptEncoding string
		for _, header := range res.headers {
			if s.Contains(header, "Accept-Encoding") {
				acceptEncoding = s.TrimSpace(s.Split(header, ": ")[1])
			}
		}
		if lenUrl > 1 {
			for _, encoding := range deleteEmpty(s.Split(acceptEncoding, ", ")) {
				if encoding == "gzip" {
					var b bytes.Buffer
					gz := gzip.NewWriter(&b)
					_, err := gz.Write([]byte(req.url[1]))
					if err != nil {
						res.addStatus(BAD_REQUEST)
						res.addBody("text/plain", "Impossible to compress the body")
						return res.constructResponse()
					}
					if err := gz.Close(); err != nil {
						res.addStatus(BAD_REQUEST)
						res.addBody("text/plain", "Impossible to compress the body")
						return res.constructResponse()
					}
					res.addStatus(OK)
					res.addHeader("Content-Encoding: gzip")
					res.addBody("text/plain", b.String())
					return res.constructResponse()
				}
			}
			res.addStatus(OK)
			res.addBody("text/plain", req.url[1])
			return res.constructResponse()
		} else {
			return BAD_REQUEST + "\r\n"
		}

	case "user-agent":
		var userAgent string
		for _, header := range req.headers {
			if s.Contains(header, "User-Agent") {
				userAgent = s.TrimSpace(s.Split(header, ": ")[1])
			}
		}
		if userAgent != "" {
			res.addStatus(OK)
			res.addBody("text/plain", userAgent)
			return res.constructResponse()
		}
	case "files":
		if lenUrl != 2 {
			res.addStatus(BAD_REQUEST)
			return res.constructResponse()
		}
		switch req.method {
		case "GET":
			filename := req.url[1]
			path := filepath.Join(fileDir, "/", filename)
			file, err := os.ReadFile(path)
			if err != nil {
				fmt.Println("Error reading file: ", err.Error())
				res.addStatus(NOT_FOUND)
				return res.constructResponse()
			}
			res.addStatus(OK)
			res.addBody("application/octet-stream", string(file))
			return res.constructResponse()
		case "POST":
			filename := req.url[1]
			file, err := os.Create(filepath.Join(fileDir, "/", filename))
			if err != nil {
				fmt.Println("Error creating file: ", err.Error())
				res.addStatus(NOT_FOUND)
				return res.constructResponse()
			}
			l, err := file.WriteString(req.body)
			if err != nil {
				fmt.Println("Error writing to file: ", err.Error())
				res.addStatus(NOT_FOUND)
				return res.constructResponse()
			}
			fmt.Println(l, "bytes written successfully")
			res.addStatus(CREATED)
			return res.constructResponse()
		}
	}
	res.addStatus(NOT_FOUND)
	return res.constructResponse()
}

func parseRequest(req string) (Request, error) {
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
	method, urlParts, _, err := parseRequestLine(requestLine)
	if err != nil {
		return Request{}, err
	}
	request := Request{method, urlParts, headers, s.TrimPrefix(body, "\n")}
	return request, nil
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
