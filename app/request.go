package main

import (
	"fmt"
	s "strings"
)

type Request struct {
	Method  string
	Path    string
	Version string
	Headers map[string]string
	Body    string
}

func parseRequestHeader(headers []string, req *Request) error {
	for _, header := range headers {
		parts := s.SplitN(header, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid header: %s", header)
		}

		key := s.TrimSpace(parts[0])
		value := s.TrimSpace(parts[1])

		req.Headers[key] = value
	}
	return nil
}

func parseRequestRequestLine(requestLine string, req *Request) error {
	fields := s.Fields(requestLine)
	if len(fields) != 3 {
		return fmt.Errorf("invalid request line: %s", requestLine)
	}
	req.Method = fields[0]
	req.Path = fields[1]
	req.Version = fields[2]
	return nil
}

func parseRequest(rawRequest string) (Request, error) {
	req := Request{
		Headers: make(map[string]string),
	}

	parts := s.SplitN(rawRequest, "\r\n\r\n", 2)

	head := parts[0]
	var body string
	if len(parts) == 2 {
		body = parts[1]
	}

	req.Body = body

	lines := s.Split(head, "\r\n")
	if len(lines) == 0 {
		return req, fmt.Errorf("empty request")
	}

	err := parseRequestRequestLine(lines[0], &req)
	if err != nil {
		return req, err
	}

	err = parseRequestHeader(lines[1:], &req)
	if err != nil {
		return req, err
	}

	return req, nil
}
