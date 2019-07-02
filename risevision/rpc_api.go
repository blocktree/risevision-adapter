package risevision

import (
	"fmt"
	"github.com/blocktree/openwallet/log"
	"github.com/imroc/req"
	"github.com/tidwall/gjson"
	"net/http"
)

type ClientInterface interface {
	Call(path string, request []interface{}) (*gjson.Result, error)
}

// A Client is a Bitcoin RPC client. It performs RPCs over HTTP using JSON
// request and responses. A Client must be configured with a secret token
// to authenticate with other Cores on the network.
type Client struct {
	BaseURL     string
	AccessToken string
	Debug       bool
	client      *req.Req
}

type Response struct {
	Code    int         `json:"code,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Message string      `json:"message,omitempty"`
	Id      string      `json:"id,omitempty"`
}

func NewClient(url string, debug bool) *Client {
	c := Client{
		BaseURL: url,
		Debug:   debug,
	}

	api := req.New()
	//trans, _ := api.Client().Transport.(*http.Transport)
	//trans.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	c.client = api

	return &c
}

// Call calls a remote procedure on another node, specified by the path.
func (c *Client) Call(path, method string, request interface{}) (*gjson.Result, error) {

	if c.client == nil {
		return nil, fmt.Errorf("API url is not setup. ")
	}

	if c.Debug {
		log.Std.Info("Start Request API...")
	}
	url := fmt.Sprintf("%s/v2%s", c.BaseURL, path)
	r, err := c.client.Do(method, url, request)

	if c.Debug {
		log.Std.Info("Request API Completed")
	}

	if c.Debug {
		log.Std.Info("%+v", r)
	}

	if err != nil {
		return nil, err
	}

	resp := gjson.ParseBytes(r.Bytes())

	if err = isError(r); err != nil {
		return nil, err
	}

	return &resp, nil
}

//isError 是否报错
func isError(r *req.Resp) error {

	if r.Response().StatusCode != http.StatusOK {
		message := gjson.GetBytes(r.Bytes(), "reason").String()
		message = fmt.Sprintf("[%s]%s", r.Response().Status, message)
		return fmt.Errorf(message)
	}

	return nil
}
