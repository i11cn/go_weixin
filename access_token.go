package weixin

import (
	"errors"
	"time"
)

type (
	WXToken interface {
		Expired()
		GetToken() string
	}

	token_in_local struct {
		WXConfig
		token_type int
		token      string
	}

	token_in_redis struct {
		WXConfig
		token_type int
	}

	token_in_etcd struct {
		WXConfig
		token_type int
		token      string
		last       time.Time
	}

	access_token_response struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
		Token   string `json:"access_token"`
		Expire  int    `json:"expires_in"`
	}
)

const (
	AccessToken = iota
	JSApiToken
	WXCardToken
)

const (
	Local = iota
	Redis
	Etcd
)

func (wx *Weixin) SetTokenSource(t, s int) (ret WXToken) {
	switch s {
	case Local:
		ret = &token_in_local{wx.WXConfig, t, ""}
	case Redis:
	case Etcd:
	}
	return
}

func (wx *Weixin) get_access_token() (string, int, error) {
	rc := wx.GetRestClient()
	r := &access_token_response{}
	if _, err := rc.Get("/token?grant_type=client_credential&appid=%s&secret=%s", r, wx.AppID, wx.AppSecret); err != nil {
		return "", 0, err
	}
	if r.ErrCode > 0 {
		return "", 0, errors.New(r.ErrMsg)
	}
	return r.Token, r.Expire, nil
}

func (t *token_in_local) start() {
}

func (t *token_in_local) Expired() {
}

func (t *token_in_local) GetToken() string {
	return t.token
}
