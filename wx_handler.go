package weixin

import (
	"net/http"
	"strings"
	"time"
)

type (
	VerifyHandle        func(token, sign, timestamp, nonce string) bool
	TextReqHandle       func(user string, t time.Time, id int64, msg string) (interface{}, error)
	ImageReqHandle      func(user string, t time.Time, id int64, url, mid string) (interface{}, error)
	VoiceReqHandle      func(user string, t time.Time, id int64, mid, format, recog string) (interface{}, error)
	VideoReqHandle      func(user string, t time.Time, id int64, mid, thumbid string) (interface{}, error)
	ShortVideoReqHandle func(user string, t time.Time, id int64, mid, thumbid string) (interface{}, error)
	LocationReqHandle   func(user string, t time.Time, id int64, x, y float64, scale int, label string) (interface{}, error)
	LinkReqHandle       func(user string, t time.Time, id int64, url, title, desc string) (interface{}, error)

	SubscribeHandle   func(user string, t time.Time, key, ticket string) (interface{}, error)
	UnsubscribeHandle func(user string, t time.Time) (interface{}, error)
	ScanHandle        func(user string, t time.Time, key, ticker string) (interface{}, error)
	LocationHandle    func(user string, t time.Time, lat, long, precision float64) (interface{}, error)
	MenuClickHandle   func(user string, t time.Time, key string) (interface{}, error)
	MenuViewHandle    func(user string, t time.Time, url string) (interface{}, error)

	WXHandler struct {
		WXComponent
		verify_handle VerifyHandle
	}
)

func NewHandler(info *WXGlobalInfo, wx *Weixin) *WXHandler {
	ret := &WXHandler{WXComponent: WXComponent{info, wx}}
	ret.verify_handle = func(token, sign, ts, nonce string) bool {
		return check_sign(token, sign, ts, nonce, ret.Logger)
	}
	return ret
}

func (wh *WXHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wh.Logger.Trace("收到一次", r.Method, "请求 : ", r.URL)
	v := r.URL.Query()
	if err := exist_all_values(v, "signature", "timestamp", "nonce"); err != nil {
		wh.Logger.Error("没有校验的签名，请求不是来自微信")
		w.WriteHeader(500)
		return
	}
	if !check_sign(wh.Config.Token, v.Get("signature"), v.Get("timestamp"), v.Get("nonce"), wh.Logger) {
		wh.Logger.Error("签名验证失败")
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

func (wh *WXHandler) Handle(fn interface{}) error {
	return wh.handle(fn)
}
