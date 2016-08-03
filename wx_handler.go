package weixin

import (
	"encoding/xml"
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

	Common struct {
		ToUserName   string
		FromUserName string
		CreateTime   int64
		MsgType      string
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
	v := r.URL.Query()
	if err := exist_all_values(v, "signature", "timestamp", "nonce"); err != nil {
		wh.log.Error("没有校验的签名，请求不是来自微信")
		w.WriteHeader(500)
		return
	}
	if !check_sign(wh.cfg.Token, v.Get("signature"), v.Get("timestamp"), v.Get("nonce"), wh.log) {
		wh.log.Error("签名验证失败")
		w.WriteHeader(301)
		return
	}
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
	if _, ok := r.URL.Query()["echostr"]; ok {
		w.Write([]byte(r.URL.Query().Get("echostr")))
	} else {
		wh.log.Error("错误的请求，没有echostr")
		w.WriteHeader(500)
	}
}

func (wh *WXHandler) doPost(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		wh.log.Error(err.Error())
	}
	wh.log.Trace(string(body))
	req := Common{}
	if err = xml.Unmarshal(body, &req); err != nil {
		wh.log.Error("解析xml出错:", err.Error())
		w.WriteHeader(500)
		return
	}
	wh.log.Trace(req)
}
