package weixin

import (
	"errors"
	"github.com/i11cn/go_logger"
	rest "github.com/i11cn/go_rest_client"
	"math/rand"
	"time"
)

type (
	WXTokenMgr interface {
		GetAccessToken() WXToken
		GetJSApiToken() WXToken
		GetWXCardToken() WXToken
		SetConfig(WXConfig)
		SetLogger(*logger.Logger)
	}

	WXToken interface {
		Expired()
		GetToken() string
		Close()
		SetSource()
		SetConfig(WXConfig)
		SetLogger(*logger.Logger)
	}

	TokenStorage interface {
		SetToken(string)
		GetToken() string
		Start()
	}
)

type (
	default_token_mgr struct {
		rc          *rest.RestClient
		log         *logger.Logger
		AccessToken WXToken
		JSApiToken  WXToken
		WXCardToken WXToken
	}
)

func (this *default_token_mgr) GetAccessToken() WXToken {
	return this.AccessToken
}

func (this *default_token_mgr) GetJSApiToken() WXToken {
	return this.JSApiToken
}

func (this *default_token_mgr) GetWXCardToken() WXToken {
	return this.WXCardToken
}

func (this *default_token_mgr) SetConfig(cfg WXConfig) {
	this.AccessToken.SetConfig(cfg)
	this.JSApiToken.SetConfig(cfg)
	this.WXCardToken.SetConfig(cfg)
}

func (this *default_token_mgr) SetLogger(log *logger.Logger) {
	this.AccessToken.SetLogger(log)
	this.JSApiToken.SetLogger(log)
	this.WXCardToken.SetLogger(log)
}

type (
	token_base struct {
		app_id     string
		app_secret string
		fn         func(rc *rest.RestClient, app_id, app_secret string, log *logger.Logger) (string, int, error)
		rc         *rest.RestClient
		log        *logger.Logger
		store      TokenStorage
		flag       chan int
	}
)

func (t *token_base) fetch_routine() {
	delay := 0
	fn := func() {
		token, exp, err := t.fn(t.rc, t.app_id, t.app_secret, t.log)
		if err != nil {
			t.log.Error(err.Error())
			delay = 10
			t.log.Log("延时", delay, "秒钟重新获取Token")
			return
		}
		delay = exp*9/10 + rand.Intn(exp/20)
		t.log.Trace("获取到的Token是:", token, "，过期时长为:", exp)
		t.log.Trace("延时", delay, "秒执行下一次获取Token的操作")
		t.store.SetToken(token)
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

func (t *token_base) GetToken() string {
	return t.store.GetToken()
}

func (t *token_base) SetConfig(cfg WXConfig) {
	t.app_id = cfg.AppID
	t.app_secret = cfg.AppSecret
}

func (t *token_base) SetLogger(log *logger.Logger) {
	t.log = log
}

func (t *token_base) SetSource(ts TokenStorage) {
	t.store = ts
	t.store.Start()
}

type (
	access_token_response struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
		Token   string `json:"access_token"`
		Expire  int    `json:"expires_in"`
	}
)

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

type (
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
)

func (t *token_in_local) Start() {
}

func (t *token_in_local) SetToken(token string) {
	t.token = token
}

func (t *token_in_local) GetToken() string {
	return t.token
}

func (t *token_in_redis) Start() {
}

func (t *token_in_redis) SetToken(token string) {
	t.token = token
}

func (t *token_in_redis) GetToken() string {
	return t.token
}

func (t *token_in_etcd) start() {
	go func() {
		// 此处监视etcd中key的变化，如有变化，去保存到t.token
	}()
}

func (t *token_in_etcd) SetToken(token string) {
	t.token = token
}

func (t *token_in_etcd) GetToken() string {
	return t.token
}
