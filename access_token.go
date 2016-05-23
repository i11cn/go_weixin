package weixin

import (
	"time"
)

type (
	WXToken interface {
		Expired()
		GetToken() string
	}

	access_token_local struct {
		WXConfig
		token string
	}

	access_token_redis struct {
		WXConfig
	}

	access_token_etcd struct {
		WXConfig
		token string
		last  time.Time
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

func (wx *Weixin) SetAccessSource(t, s int) {
	switch t {
	case AccessToken:
	case JSApiToken:
	case WXCardToken:
	}
}

func (wx *access_token_local) Expired() {
}

func (wx *access_token_local) GetToken() string {
	return wx.token
}
