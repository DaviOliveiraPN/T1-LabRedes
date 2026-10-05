package http1

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"
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
	r.Headers["Content-Length"] = fmt.Sprint(len(r.Body))
	r.Headers["Date"] = httpDate()
	r.Headers["Server"] = "LabRedes-Go/1.0"
	if _, ok := r.Headers["Content-Type"]; !ok {
		r.Headers["Content-Type"] = "text/plain; charset=utf-8"
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "HTTP/1.1 %d %s\r\n", r.StatusCode, statusText[r.StatusCode])
	for k, v := range r.Headers {
		fmt.Fprintf(&sb, "%s: %s\r\n", k, v)
	}
	sb.WriteString("\r\n")

	out := []byte(sb.String())
	if !head {
		out = append(out, r.Body...)
	}
	_, err := w.Write(out)
	return err
}

// IMF-fixdate (RFC 9110), sempre em GMT.
func httpDate() string {
	return time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05") + " GMT"
}

var mimeTypes = map[string]string{
	".html": "text/html; charset=utf-8",
	".css":  "text/css",
	".js":   "application/javascript",
	".json": "application/json",
	".txt":  "text/plain; charset=utf-8",
	".png":  "image/png",
	".jpg":  "image/jpeg",
	".pdf":  "application/pdf",
}

func contentType(path string) string {
	if t, ok := mimeTypes[strings.ToLower(filepath.Ext(path))]; ok {
		return t
	}
	return "application/octet-stream"
}
