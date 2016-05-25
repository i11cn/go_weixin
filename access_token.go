package weixin

import (
	"errors"
	"github.com/i11cn/go_logger"
	rest "github.com/i11cn/go_rest_client"
	"math/rand"
	"time"
)

type (
	WXToken interface {
		Expired()
		GetToken() string
		Close()
	}

	Tokens struct {
		AccessToken WXToken
		JSApiToken  WXToken
		WXCardToken WXToken
	}

	token_base struct {
		WXConfig
		fn   func(rc *rest.RestClient, app_id, app_secret string, log *logger.Logger) (string, int, error)
		rc   *rest.RestClient
		log  *logger.Logger
		tch  chan string
		flag chan int
	}
	token_in_local struct {
		token_base
		token string
	}
	token_in_redis struct {
		token_base
		key   string
		token string // 临时用的，将来需要以key存入redis中
	}
	token_in_etcd struct {
		token_base
		key   string
		token string
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

func (wx *Weixin) SetTokenSource(t, s int, main bool) (ret WXToken) {
	switch t {
	case AccessToken, JSApiToken, WXCardToken:
	default:
		return
	}
	switch s {
	case Local:
		ret = (&token_in_local{token_base{wx.cfg, get_access_token, wx.rc, wx.log, make(chan string), make(chan int)}, ""}).start()
	case Redis:
		ret = (&token_in_redis{token_base{wx.cfg, get_access_token, wx.rc, wx.log, make(chan string), make(chan int)}, "access_token", ""}).start(main)
	case Etcd:
		ret = (&token_in_etcd{token_base{wx.cfg, get_access_token, wx.rc, wx.log, make(chan string), make(chan int)}, "access_token", ""}).start(main)
	default:
		return nil
	}
	return
}

func get_access_token(rc *rest.RestClient, app_id, app_secret string, log *logger.Logger) (string, int, error) {
	r := &access_token_response{}
	if _, err := rc.Get("/token?grant_type=client_credential&appid=%s&secret=%s", r, app_id, app_secret); err != nil {
		log.Error(err.Error())
		return "", 0, err
	}
	if r.ErrCode > 0 {
		log.Error(r.ErrMsg)
		return "", 0, errors.New(r.ErrMsg)
	}
	log.Trace("获取到了Token:", r.Token, " , 有效时长为:", r.Expire, "秒")
	return r.Token, r.Expire, nil
}

func (t *token_base) fetch_routine() {
	delay := 0
	fn := func() {
		token, exp, err := t.fn(t.rc, t.AppID, t.AppSecret, t.log)
		if err != nil {
			t.log.Error(err.Error())
			delay = 10
			t.log.Log("延时", delay, "秒钟重新获取Token")
			return
		}
		delay = exp*9/10 + rand.Intn(exp/20)
		t.log.Trace("获取到的Token是:", token, "，过期时长为:", exp)
		t.log.Trace("延时", delay, "秒执行下一次获取Token的操作")
		t.tch <- token
	}
	for {
		select {
		case i := <-t.flag:
			if i == 0 {
				t.log.Log("从微信服务器获取token的routing退出")
				return
			}
			t.log.Log("强制更新token")
			fn()
		case <-time.After(time.Duration(delay) * time.Second):
			fn()
		}
	}
}

func (t *token_base) Expired() {
	t.flag <- 1
}

func (t *token_base) Close() {
	t.flag <- 0
}

func (t *token_in_local) start() *token_in_local {
	t.log.Trace("启动token的维护工作")
	go t.fetch_routine()
	go func(t *token_in_local) {
		for {
			tk, ok := <-t.tch
			if !ok {
				return
			}
			t.token = tk
		}
	}(t)
	return t
}

func (t *token_in_local) GetToken() string {
	return t.token
}

func (t *token_in_redis) start(main bool) *token_in_redis {
	if main {
		t.log.Trace("启动token的维护工作")
		go t.fetch_routine()
		go func(t *token_in_redis) {
			for {
				tk, ok := <-t.tch
				if !ok {
					return
				}
				t.token = tk
			}
		}(t)
	}
	return t
}

func (t *token_in_redis) GetToken() string {
	return t.token
}

func (t *token_in_etcd) start(main bool) *token_in_etcd {
	if main {
		t.log.Trace("启动token的维护工作")
		go t.fetch_routine()
		go func(t *token_in_etcd) {
			for {
				tk, ok := <-t.tch
				if !ok {
					return
				}
				t.token = tk
				// 此处再将token以key为键值存入etcd
			}
		}(t)
	}
	go func() {
		// 此处监视etcd中key的变化，如有变化，去除保存到t.token
	}()
	return t
}

func (t *token_in_etcd) GetToken() string {
	return t.token
}
