package weixin

import (
	"github.com/i11cn/go_rest_client"
	"net/http"
)

type (
	WXClient struct {
		rest_client *rc.RestClient
		wx          *Weixin
		token_name  string
	}
)

func NewClient(wx *Weixin, host string, port int, uri, token_name string) *WXClient {
	ret := &WXClient{rc.NewClient(host, port, uri), wx, token_name}
	return ret
}

func (c *WXClient) Get(uri string, obj interface{}, args ...interface{}) (*http.Response, error) {
	q := make(map[string]interface{})
	q[c.token_name] = c.wx.GetAccessToken().GetToken()
	return c.rest_client.GetCaller("GET", uri)(&rc.Request{q, nil, obj}, args...)
}

func (c *WXClient) Post(uri string, body interface{}, obj interface{}, args ...interface{}) (*http.Response, error) {
	q := make(map[string]interface{})
	q[c.token_name] = c.wx.GetAccessToken().GetToken()
	return c.rest_client.GetCaller("POST", uri)(&rc.Request{q, body, obj}, args...)
}
