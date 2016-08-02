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
		SetGlobalInfo(*WXGlobalInfo)
	}

	TokenStorage interface {
		SetToken(string)
		GetToken() string
		Start()
		Stop()
	}

	WXToken interface {
		Expired()
		GetToken() string
		Close()
		Primary(bool)
		SetSource(TokenStorage)
		SetGlobalInfo(*WXGlobalInfo)
	}
)

type (
	default_token_mgr struct {
		AccessToken WXToken
	}
)

func new_token(fn func(rc *rest.RestClient, app_id, app_secret string, log *logger.Logger) (string, int, error), main bool) WXToken {
	ret := &token_base{main: main}
	ret.fn = fn
	ret.SetSource(&token_in_local{})
	ret.Primary(main)
	return ret
}

func DefaultTokenMgr(info *WXGlobalInfo) WXTokenMgr {
	ret := &default_token_mgr{new_token(get_access_token, true)}
	return ret
}

func (this *default_token_mgr) GetAccessToken() WXToken {
	return this.AccessToken
}

func (this *default_token_mgr) SetGlobalInfo(info *WXGlobalInfo) {
	this.AccessToken.SetGlobalInfo(info)
}

type (
	token_base struct {
		info  *WXGlobalInfo
		fn    func(rc *rest.RestClient, app_id, app_secret string, log *logger.Logger) (string, int, error)
		store TokenStorage
		main  bool
		flag  chan int
	}
)

func (t *token_base) fetch_routine() {
	delay := 0
	info := t.info
	fn := func() {
		token, exp, err := t.fn(info.RestClient, info.Config.AppID, info.Config.AppSecret, info.Log)
		if err != nil {
			info.Log.Error(err.Error())
			delay = 10
			info.Log.Log("延时", delay, "秒钟重新获取Token")
			return
		}
		delay = exp*9/10 + rand.Intn(exp/20)
		info.Log.Trace("获取到的Token是:", token, "，过期时长为:", exp)
		info.Log.Trace("延时", delay, "秒执行下一次获取Token的操作")
		t.store.SetToken(token)
	}
	for {
		select {
		case i := <-t.flag:
			if i == 0 {
				info.Log.Log("从微信服务器获取token的routing退出")
				return
			}
			info.Log.Log("强制更新token")
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

func (t *token_base) Primary(main bool) {
	if t.main == main {
		return
	}
	if t.main {
		t.Close()
	} else {
		go t.fetch_routine()
	}
	t.main = main
}

func (t *token_base) GetToken() string {
	return t.store.GetToken()
}

func (t *token_base) SetGlobalInfo(info *WXGlobalInfo) {
	t.info = info
}

func (t *token_base) SetSource(ts TokenStorage) {
	if t.store != nil {
		t.store.Stop()
	}
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

func (t *token_in_local) Stop() {
}

func (t *token_in_local) SetToken(token string) {
	t.token = token
}

func (t *token_in_local) GetToken() string {
	return t.token
}

func (t *token_in_redis) Start() {
}

func (t *token_in_redis) Stop() {
}

func (t *token_in_redis) SetToken(token string) {
	t.token = token
}

func (t *token_in_redis) GetToken() string {
	return t.token
}

func (t *token_in_etcd) Start() {
	go func() {
		// 此处监视etcd中key的变化，如有变化，去保存到t.token
	}()
}

func (t *token_in_etcd) Stop() {
}

func (t *token_in_etcd) SetToken(token string) {
	t.token = token
}

func (t *token_in_etcd) GetToken() string {
	return t.token
}
