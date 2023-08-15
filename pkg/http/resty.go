package http

import "github.com/go-resty/resty/v2"

const PostMethod = "POST"
const GetMethod = "GET"

type Method string

type Client struct {
	resty       *resty.Client
	url         string
	queryString string
	body        interface{}
	pathParams  map[string]string
	headers     map[string]string
}

func NewHttpClient() *Client {
	return &Client{
		headers:    make(map[string]string, 0),
		pathParams: make(map[string]string, 0),
	}
}

func (cl *Client) SetUrl(url string) *Client {
	cl.url = url
	return cl
}

func (cl *Client) SetBody(body interface{}) *Client {
	cl.body = body
	return cl
}

func (cl *Client) SetPathParams(pathParams map[string]string) *Client {
	cl.pathParams = pathParams
	return cl
}

func (cl *Client) SetQueryString(queryString string) *Client {
	cl.queryString = queryString
	return cl
}

func (cl *Client) SetHeader(key string, val string) *Client {
	cl.headers[key] = val
	return cl
}

func (cl *Client) Send(method string, res interface{}) (*resty.Response, error) {
	cl.resty = resty.New()
	var (
		req = cl.resty.R()
		err error
	)
	if len(cl.headers) > 0 {
		for k, val := range cl.headers {
			req.SetHeader(k, val)
		}
		if _, ok := cl.headers["Content-Type"]; !ok {
			req.SetHeader("Content-Type", "application/json")
		}
	} else {
		req.SetHeader("Content-Type", "application/json")
	}
	if len(cl.queryString) > 0 {
		req.SetQueryString(cl.queryString)
	}
	if len(cl.pathParams) > 0 {
		req.SetPathParams(cl.pathParams)
	}
	if cl.body != nil {
		req.SetBody(cl.body)
	}
	req.SetResult(res).ForceContentType("application/json")
	switch method {
	case GetMethod:
		return req.Get(cl.url)
	case PostMethod:
		return req.Post(cl.url)
	default:
		return nil, err
	}
}
