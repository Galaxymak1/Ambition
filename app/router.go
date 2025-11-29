package main

import (
	"fmt"
	s "strings"
)

func Router(req *Request, res *Response) {
	parts := s.Split(s.Trim(req.Path, "/"), "/")
	fmt.Println(parts)
	switch parts[0] {
	case "":
		handleRoot(req, res)
	case "echo":
		handleEcho(req, res, parts[1:])
	default:
		res.statusCode = 404
		res.body = []byte("Not Found")
	}
}

func handleRoot(req *Request, res *Response) {
	switch req.Method {
	case "GET":
		res.statusCode = 200
		res.body = []byte("Hello World")
	default:
		res.statusCode = 405
		res.body = []byte("Method Not Allowed")
	}
}

func handleEcho(req *Request, res *Response, params []string) {
	if req.Method != "GET" {
		res.statusCode = 405
		res.body = []byte("Method Not Allowed")
		return
	}
	if len(params) == 0 {
		res.statusCode = 400
		res.body = []byte("Missing parameter")
		return
	}
	echoValue := params[0]
	res.statusCode = 200
	res.body = []byte(echoValue)
}
