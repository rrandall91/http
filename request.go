package http

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Request represents an HTTP request.
type Request struct {
	Method  string
	URL     string
	Body io.Reader
	Headers []Param
	Query   []Param
}

// NewRequest creates a new request.
func NewRequest(method, url string) *Request {
	return &Request{
		Method: method,
		URL:    url,
	}
}

// AddHeader adds a header to the request.
func (r *Request) AddHeader(key, value string) {
	r.Headers = append(r.Headers, newParam(key, value))
}

// AddQuery adds a query to the request.
func (r *Request) AddQuery(key, value string) {
	r.Query = append(r.Query, newParam(key, value))
}

// AddBody adds a body to the request.
func (r *Request) AddBody(body io.Reader) {
	r.Body = body
}

// AddBodyString adds a string body to the request.
func (r *Request) AddBodyString(body string) {
	r.Body = io.NopCloser(strings.NewReader(body))
}

// AddBodyJSON adds a JSON body to the request.
func (r *Request) AddBodyJSON(body interface{}) {
	r.AddHeader("Content-Type", "application/json")

	b, err := json.Marshal(body)
	if err != nil {
		return
	}
	
	r.AddBodyString(string(b))
}

// AddBodyXML adds an XML body to the request.
func (r *Request) AddBodyXML(body interface{}) {
	r.AddHeader("Content-Type", "application/xml")

	b, err := xml.Marshal(body)
	if err != nil {
		return
	}

	r.AddBodyString(string(b))
}

// AddBodyForm adds a form body to the request.
func (r *Request) AddBodyForm(body map[string]string) {
	r.AddHeader("Content-Type", "application/x-www-form-urlencoded")

	form := url.Values{}
	for key, value := range body {
		form.Add(key, value)
	}

	r.AddBodyString(form.Encode())
}

// AddBodyMultipartForm adds a multipart form body to the request.
func (r *Request) AddBodyMultipartForm(body map[string]string) {
	r.AddHeader("Content-Type", "multipart/form-data")
	
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	for key, value := range body {
		w.WriteField(key, value)
	}
	w.Close()

	r.AddBody(&b)
}

// AddBodyMultipartFormFile adds a multipart form body with a file to the request.
func (r *Request) AddBodyMultipartFormFile(body map[string]string, fileKey, fileName, filePath string) {
	r.AddHeader("Content-Type", "multipart/form-data")

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	for key, value := range body {
		w.WriteField(key, value)
	}

	fw, err := w.CreateFormFile(fileKey, fileName)
	if err != nil {
		return
	}

	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer f.Close()

	_, err = io.Copy(fw, f)

	w.Close()

	r.AddBody(&b)

}

// GetHeader returns the value of the header with the given key.
func (r *Request) GetHeader(key string) string {
	for _, header := range r.Headers {
		if header.Key == key {
			return header.Value
		}
	}
	return ""
}

// GetQuery returns the value of the query with the given key.
func (r *Request) GetQuery(key string) string {
	for _, query := range r.Query {
		if query.Key == key {
			return query.Value
		}
	}
	return ""
}

// make creates a new http.Request from the request.
func (r *Request) make() *http.Request {
	req, err := http.NewRequest(r.Method, r.URL, r.Body)
	if err != nil {
		return nil
	}

	for _, header := range r.Headers {
		req.Header.Add(header.Key, header.Value)
	}

	q := req.URL.Query()
	for _, query := range r.Query {
		q.Add(query.Key, query.Value)
	}

	req.URL.RawQuery = q.Encode()

	req.Body = http.NoBody

	return req
}

// Send executes the request and returns a response.
func (r *Request) Send() (*Response, error) {
	start := time.Now()

	req := r.make()
	if req == nil {
		return nil, nil
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	end := time.Now()

	return &Response{
		StatusCode: resp.StatusCode,
		Body:       resp.Body,
		Duration:  end.Sub(start),
	}, nil
}
