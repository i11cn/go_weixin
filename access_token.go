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

	token_in_local struct {
		WXConfig
		token_type int
		fn         func(rc *rest.RestClient, app_id, app_secret string, log *logger.Logger) (string, int, error)
		token      string
		rc         *rest.RestClient
		log        *logger.Logger
		flag       chan int
	}

	token_in_redis struct {
		WXConfig
		token_type int
		rc         *rest.RestClient
		log        *logger.Logger
	}

	token_in_etcd struct {
		WXConfig
		token_type int
		token      string
		last       time.Time
		rc         *rest.RestClient
		log        *logger.Logger
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
	switch t {
	case AccessToken, JSApiToken, WXCardToken:
	default:
		return
	}
	switch s {
	case Local:
		ret = (&token_in_local{wx.cfg, t, get_access_token, "", wx.rc, wx.log, make(chan int)}).start()
	case Redis:
	case Etcd:
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

func (t *token_in_local) start() *token_in_local {
	t.log.Trace("启动token的维护工作")
	fn := func(t *token_in_local) {
		delay := 0
		for {
			select {
			case i := <-t.flag:
				if i == 0 {
					t.log.Log("token维护routing退出")
					return
				}
				t.log.Log("强制更新token")
				token, exp, err := t.fn(t.rc, t.AppID, t.AppSecret, t.log)
				if err != nil {
					t.log.Error(err.Error())
					t.log.Log("延时10秒钟重新获取Token")
					delay = 10
					continue
				}
				t.token = token
				delay = exp*9/10 + rand.Intn(exp/20)
				t.log.Trace("获取到的Token是:", token, "，过期时长为:", exp)
				t.log.Trace("延时", delay, "秒执行下一次获取Token的操作")
			case <-time.After(time.Duration(delay) * time.Second):
				token, exp, err := t.fn(t.rc, t.AppID, t.AppSecret, t.log)
				if err != nil {
					t.log.Log("延时10秒钟重新获取Token")
					delay = 10
					continue
				}
				t.token = token
				delay = exp*9/10 + rand.Intn(exp/20)
				t.log.Trace("获取到的Token是:", token, "，过期时长为:", exp)
				t.log.Trace("延时", delay, "秒执行下一次获取Token的操作")
			}
		}
	}
	go fn(t)
	return t
}

func (t *token_in_local) Expired() {
	t.flag <- 1
}

func (t *token_in_local) GetToken() string {
	return t.token
}

func (t *token_in_local) Close() {
	t.flag <- 0
}
