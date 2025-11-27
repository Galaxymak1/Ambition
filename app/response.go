package main

import "strconv"

type Response struct {
	statusLine string
	headers    []string
	body       string
}

type Request struct {
	method  string
	url     []string
	headers []string
	body    string
}

func (res *Response) addStatus(statusLine string) {
	res.statusLine = statusLine
}

func (res *Response) addHeader(header string) {
	for i, h := range res.headers {
		if h == header {
			res.headers[i] = header
			return
		}
	}
	res.headers = append(res.headers, header)
}
func (res *Response) addBody(contentType string, body string) {
	res.addHeader("Content-Type: " + contentType)
	res.addHeader("Content-Length: " + strconv.Itoa(len(body)))
	res.body = body
}

func (res *Response) constructResponse() string {
	fullResponse := res.statusLine
	for _, header := range res.headers {
		fullResponse += header + "\r\n"
	}
	fullResponse += "\r\n"
	fullResponse += res.body
	return fullResponse
}
