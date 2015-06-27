package weixin

import (
	"encoding/xml"
	"github.com/i11cn/go_logger"
	"net/http"
	"time"
)

type WXRequestInfo struct {
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgId        int64
}

type WXTextRequest struct {
	Content string
}

type WXLocationRequest struct {
	Location_X float64
	Location_Y float64
	Scale      int
	Label      string
}

type WXRequest struct {
	WXRequestInfo
	MsgType string
	WXTextRequest
	WXLocationRequest
}

type WXResponseInfo struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   int64
}

type WXTextResponse struct {
	WXResponseInfo
	MsgType string
	Content string
}

type WXConfig struct {
	Token     string
	MsgKey    string
	AppID     string
	AppSecret string
}

type OnValidateFail func(w http.ResponseWriter, r *http.Request)
type UnsupportedRequest func(w http.ResponseWriter, info *WXRequestInfo)
type OnRequestError func(w http.ResponseWriter, r *http.Request)

type OnTextRequest func(w http.ResponseWriter, info *WXRequestInfo, text *WXTextRequest)
type OnLocationRequest func(w http.ResponseWriter, info *WXRequestInfo, pos *WXLocationRequest)

func (info *WXRequestInfo) Response(w http.ResponseWriter, v interface{}) {
	output, err := xml.MarshalIndent(v, "", "")
	if err != nil {
		go_logger.GetLogger("weixin").Error("创建响应xml失败:", err.Error())
		w.WriteHeader(500)
		return
	}
	w.Write(output)
	go_logger.GetLogger("weixin").Trace("给用户返回响应:", string(output))
}

func (info *WXRequestInfo) ResponseText(w http.ResponseWriter, content string) {
	resp := WXTextResponse{}
	resp.ToUserName = info.FromUserName
	resp.FromUserName = info.ToUserName
	resp.CreateTime = time.Now().Unix()
	resp.MsgType = "text"
	resp.Content = content
	info.Response(w, resp)
}

func NewWeixinServ(conf *WXConfig, uc interface{}) *Weixin {
	serv := &Weixin{WXConfig: *conf, UserCustom: uc}
	if serv.init() {
		return serv
	} else {
		return nil
	}
}
