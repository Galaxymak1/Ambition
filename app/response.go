package main

import (
	"strconv"
	s "strings"
)

type Response struct {
	statusCode int
	headers    map[string]string
	body       []byte
}

func finalizeResponse(req *Request, res *Response) {
	if supportsGzip(req) {
		compressedBody, _ := compressToGzip(res.body)
		res.body = compressedBody
		res.headers["Content-Encoding"] = "gzip"
	}
	res.headers["Content-Length"] = strconv.Itoa(len(res.body))

	if _, ok := res.headers["Content-Type"]; !ok {
		res.headers["Content-Type"] = "text/plain"
	}
}

func serializeResponse(res *Response) []byte {
	statusText := statusMap[res.statusCode]
	startLine := "HTTP/1.1 " + strconv.Itoa(res.statusCode) + " " + statusText + "\r\n"
	headerLines := ""
	for key, value := range res.headers {
		headerLines += key + ": " + value + "\r\n"
	}
	fullHeader := startLine + headerLines + "\r\n"
	headerBytes := []byte(fullHeader)
	response := append(headerBytes, res.body...)
	return response
}

func supportsGzip(req *Request) bool {
	if value, found := req.Headers["Content-Encoding"]; found {
		for _, encoding := range deleteEmpty(s.Split(value, ", ")) {
			if encoding == "gzip" {
				return true
			}
		}
	}
	return false
}

//
//func (res *Response) setStatues(statusCode string) {
//	res.statusCode = statusCode
//}
//
//func (res *Response) addHeader(header string, value string) {
//	res.headers[header] = value
//}
//func (res *Response) addBody(body string) {
//	//res.addHeader("Content-Type: " + contentType)
//	//res.addHeader("Content-Length: " + strconv.Itoa(len(body)))
//	res.body = body
//}
//
//func (res *Response) constructResponse() string {
//	fullResponse := res.statusLine
//	for _, header := range res.headers {
//		fullResponse += header + "\r\n"
//	}
//	fullResponse += "\r\n"
//	fullResponse += res.body
//	return fullResponse
//}
