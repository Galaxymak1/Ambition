package main

// 200 Codes
const OK = "HTTP/1.1 200 OK\r\n"
const CREATED = "HTTP/1.1 201 Created\r\n"

// Error Codes
const BAD_REQUEST = "HTTP/1.1 400 Bad Request\r\n"
const UNAUTHORIZED = "HTTP/1.1 401 Unauthorized\r\n"
const FORBIDDEN = "HTTP/1.1 403 Forbidden\r\n"
const NOT_FOUND = "HTTP/1.1 404 Not Found\r\n"
const METHOD_NOT_ALLOWED = "HTTP/1.1 405 Method Not Allowed\r\n"
const UNPROCESSABLE_CONTENT = "HTTP/1.1 407 Unprocessable Content\r\n"
const INTERNAL_ERROR = "HTTP/1.1 500 Internal Server Error\r\n"
