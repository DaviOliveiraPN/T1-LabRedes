package http1

import (
	"io"
)

type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

var statusText = map[int]string{
	200: "OK",
	400: "Bad Request",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Allowed",
}

func newResponse(status int, body []byte) *Response {
	return &Response{StatusCode: status, Headers: map[string]string{}, Body: body}
}

func errorResponse(status int) *Response {
	return newResponse(status, []byte(statusText[status]+"\n"))
}

func (r *Response) Write(w io.Writer, head bool) error {
	return nil
}

func httpDate() string {
	return ""
}

func contentType(path string) string {
	return "application/octet-stream"
}
