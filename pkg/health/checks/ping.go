package checks

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type PingCheck struct {
	URL     string
	Method  string
	Timeout int
	client  http.Client
	Body    io.Reader
	Headers map[string]string
}

func (p PingCheck) Name() string {
	return p.Method + " " + p.URL
}

func NewPingCheck(url, method string, timeout int, body io.Reader, headers map[string]string) PingCheck {
	if method == "" {
		method = "GET"
	}

	if timeout == 0 {
		timeout = 500
	}

	pingCheck := PingCheck{
		URL:     url,
		Method:  method,
		Timeout: timeout,
		Body:    body,
		Headers: headers,
	}
	pingCheck.client = http.Client{
		Timeout: time.Duration(timeout) * time.Millisecond,
	}

	return pingCheck
}
func (p PingCheck) Check(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, p.Method, p.URL, p.Body)

	if err != nil {
		return err
	}

	for key, value := range p.Headers {
		req.Header.Add(key, value)
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode > http.StatusBadRequest {
		return fmt.Errorf("statusCode to high, code: %d", resp.StatusCode)
	}
	return nil
}
