package weixin

import (
	"errors"
	"github.com/i11cn/go_logger"
	"io/ioutil"
	"net/http"
	"strings"
)

type (
	WXHandler struct {
		cfg   WXConfig
		token WXToken
		log   *logger.Logger
	}
)

func NewHandler(cfg WXConfig, tk WXToken, log *logger.Logger) (*WXHandler, error) {
	if tk != nil {
		return &WXHandler{cfg, tk, log}, nil
	} else {
		return nil, errors.New("没有初始化AccessToken，搞个P啊")
	}
}

func (wh *WXHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wh.log.Trace("收到一次", r.Method, "请求 : ", r.URL)
	switch strings.ToUpper(r.Method) {
	case "GET":
		wh.doGet(w, r)
	case "POST":
		wh.doPost(w, r)
	default:
		w.WriteHeader(500)
	}
}

func (wh *WXHandler) doGet(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	if err := exist_all_values(v, "signature", "timestamp", "nonce", "echostr"); err != nil {
		wh.log.Log("目前好像只有验证服务器的请求是GET方法，怎么有个这么乱的请求？")
		w.WriteHeader(500)
		return
	}
	if check_sign(wh.cfg.Token, v.Get("signature"), v.Get("timestamp"), v.Get("nonce"), wh.log) {
		w.Write([]byte(v.Get("echostr")))
	} else {
		wh.log.Error("签名验证失败")
		w.WriteHeader(301)
	}
}

func (wh *WXHandler) doPost(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		wh.log.Error(err.Error())
	}
	wh.log.Trace(string(body))
}
