package weixin

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"
)

type (
	WXHandler struct {
		WXComponent
	}

	Common struct {
		ToUserName   string
		FromUserName string
		CreateTime   int64
		MsgType      string
	}
)

func NewHandler(info *WXGlobalInfo, wx *Weixin) (*WXHandler, error) {
	return &WXHandler{WXComponent{info, wx}}, nil
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

func (wh *WXHandler) doGet(w http.ResponseWriter, r *http.Request) {
	if _, ok := r.URL.Query()["echostr"]; ok {
		w.Write([]byte(r.URL.Query().Get("echostr")))
	} else {
		wh.Logger.Error("错误的请求，没有echostr")
		w.WriteHeader(500)
	}
}

func (wh *WXHandler) doPost(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		wh.Logger.Error(err.Error())
	}
	wh.Logger.Trace(string(body))
	req := Common{}
	if err = xml.Unmarshal(body, &req); err != nil {
		wh.Logger.Error("解析xml出错:", err.Error())
		w.WriteHeader(500)
		return
	}
	wh.Logger.Trace(req)
}
