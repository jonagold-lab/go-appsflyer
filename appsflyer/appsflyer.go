package appsflyer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/google/go-querystring/query"
)

const (
	defaultBaseURL = "https://hq.appsflyer.com/export/"
	userAgent      = "go-appsflyer"
)

// https://hq.appsflyer.com/export/com.greatapp/installs_report/v5?api_token=xxxx&from=2018-04-09&to=2018-05-09&category=standard&media_source=googleadwords_int
// New client init with session

// A Client manages communication with the Adjust API.
type Client struct {
	client    *http.Client
	BaseURL   *url.URL
	APIToken  string
	AppID     string
	UserAgent string
	common    service
	Report    *ReportService
}

type service struct {
	client *Client
}

func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// NewClient with either http.Client or with user and pw
func NewClient(httpClient *http.Client, apiToken, appID string) (*Client, error) {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent, APIToken: apiToken, AppID: appID}
	c.common.client = c
	c.Report = (*ReportService)(&c.common)
	return c, nil
}

// NewRequest to build request
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}

	u, err := c.BaseURL.Parse(c.AppID + "/" + urlStr)
	if err != nil {
		return nil, err
	}
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	return req, nil
}

// Do to execute request and handle repsonse of all services
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	req.WithContext(ctx)
	resp, err := c.client.Do(req)
	if err != nil {
		// If context is cancled then return that error
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		// If the error type is *url.Error
		if e, ok := err.(*url.Error); ok {
			if url, err := url.Parse(e.URL); err == nil {
				e.URL = url.String()
				return nil, e
			}
		}

		return nil, err
	}
	defer resp.Body.Close()
	response := newResponse(resp)
	err = CheckResponse(resp)
	if err != nil {
		return response, err
	}
	err = json.NewDecoder(resp.Body).Decode(&v)
	return response, err
}

// Response to hold adjust response
type Response struct {
	*http.Response
}

// newResponse creates a new Response for the provided http.Response.
// r must not be nil.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	return response
}

// CheckResponse to build an erro if code is not 2xx
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	return errorResponse
}

/*
An ErrorResponse reports one or more errors caused by an API request.
*/
type ErrorResponse struct {
	Response *http.Response // HTTP response that caused this error
	Errors   []string       `json:"error"` // more detail on individual errors
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %+v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Errors)
}
