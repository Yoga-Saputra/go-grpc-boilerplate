package httpclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// HTTPResponse define result of HTTP request response
type HTTPResponse struct {
	StatusCode  int                    `json:"statusCode"`
	BodyString  string                 `json:"bodyString"`
	BodyArray   map[string]interface{} `json:"bodyArray"`
	BodyBytes   []byte                 `json:"bodyBytes"`
	Latency     string                 `json:"latency"`
	Cookies     []*http.Cookie         `json:"cookies"`
	RawResponse *http.Response         `json:"rawResponse"`
}

// New creates new http client
func New(config ...Config) *Config {
	// Set default config
	cfg := configDefault(config...)

	return cfg
}

// Token will set bearer token and return modify config pointer
func (c *Config) Token(t string) *Config {
	c.token = t
	return c
}

// SetTimeout will set request timeout in seconds
func (c *Config) SetTimeout(t int) *Config {
	c.timeout = time.Duration(t) * time.Second
	return c
}

// Get do http request method GET and return HTTPResponse.
// otherwise return error
func (c *Config) Get(uri string, params ...map[string]interface{}) (*HTTPResponse, error) {
	uri = fmt.Sprintf("%v%v", c.BaseUrl, uri)

	// Create new request
	req, err := c.acquireRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	// Set query params
	if len(params) > 0 {
		q := req.URL.Query()
		for k, v := range params[0] {
			q.Add(k, fmt.Sprintf("%v", v))
		}

		req.URL.RawQuery = q.Encode()
	}

	// Do client request
	resp, err := c.acquireResponse(req)

	return resp, err
}

// Post do http request method POST and return HTTPResponse,
// otherwise return error
// It will put given params to JSON body payload
func (c *Config) Post(uri string, params ...map[string]interface{}) (*HTTPResponse, error) {
	uri = fmt.Sprintf("%v%v", c.BaseUrl, uri)

	// Set body payload
	var payloads io.Reader
	if len(params) <= 0 {
		payloads = nil
	} else {
		p, err := setupPostPayload(params[0], true)
		if err != nil {
			return nil, err
		}

		payloads = p
	}

	// Create new request
	req, err := c.acquireRequest("POST", uri, payloads)
	if err != nil {
		return nil, err
	}

	// Do client request
	resp, err := c.acquireResponse(req)

	return resp, err
}

// Post do http request method POST and return HTTPResponse,
// otherwise return error
// It will put given params to Form URL-Encoded as a body params
func (c *Config) PostAsForm(uri string, params ...map[string]interface{}) (*HTTPResponse, error) {
	uri = fmt.Sprintf("%v%v", c.BaseUrl, uri)

	// Set spesific header values
	for range c.Headers {
		c.Headers["Content-Type"] = "application/x-www-form-urlencoded"
	}

	// Set body payload
	var payloads io.Reader
	if len(params) <= 0 {
		payloads = nil
	} else {
		p, err := setupPostPayload(params[0], false)
		if err != nil {
			return nil, err
		}

		payloads = p
	}

	// Create new request
	req, err := c.acquireRequest("POST", uri, payloads)
	if err != nil {
		return nil, err
	}

	// Do client request
	resp, err := c.acquireResponse(req)

	return resp, err
}

// GetJSONString will marshaling HTTP Response to json string
func (h *HTTPResponse) GetJSONString() string {
	if h != nil {
		rj := map[string]interface{}{
			"statusCode": h.StatusCode,
			"bodyString": h.BodyString,
			"latency":    h.Latency,
		}
		b, _ := json.Marshal(rj)

		return string(b)
	}

	return ""
}

// Create HTTP client
func (c *Config) acquireClient() *http.Client {
	cln := &http.Client{}

	// Set timeout
	if c.timeout > 0 {
		cln.Timeout = c.timeout
	}

	// Create new client with skip TLS verify
	if c.SkipInsecure {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: c.SkipInsecure,
			},
		}

		cln.Transport = tr
	}

	cln.Jar = c.CookieJar
	return cln
}

// Create HTTP Client request
func (c *Config) acquireRequest(method, uri string, params io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, uri, params)
	if err != nil {
		return nil, err
	}

	// Set request headers
	for k, v := range c.Headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))

		if len(c.token) > 0 {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.token))
		}
	}

	return req, nil
}

// Create HTTP response
func (c *Config) acquireResponse(req *http.Request) (*HTTPResponse, error) {
	start := time.Now()

	client := c.acquireClient()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	bodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	bodyArr := make(map[string]interface{})

	json.Unmarshal([]byte(bodyByte), &bodyArr)

	bodyStr := string(bodyByte)
	lat := time.Since(start).Seconds()

	return &HTTPResponse{
		StatusCode:  resp.StatusCode,
		BodyString:  bodyStr,
		BodyArray:   bodyArr,
		BodyBytes:   bodyByte,
		Latency:     fmt.Sprintf("%.4fs", lat),
		Cookies:     c.CookieJar.Cookies(resp.Request.URL),
		RawResponse: resp,
	}, nil
}

// Set the body payload of POST request
func setupPostPayload(params map[string]interface{}, isJson bool) (io.Reader, error) {
	if isJson {
		jsonByte, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}

		return bytes.NewReader(jsonByte), nil
	}

	p := url.Values{}
	for k, v := range params {
		p.Set(k, fmt.Sprintf("%v", v))
	}

	return strings.NewReader(p.Encode()), nil
}
