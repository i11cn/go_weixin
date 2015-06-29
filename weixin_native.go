package weixin

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/i11cn/go_logger"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

func exist_all_values(values url.Values, keys []string) bool {
	for _, key := range keys {
		if _, ok := values[key]; !ok {
			return false
		}
	}
	return true
}

func (serv *Weixin) init() bool {
	if serv.UserCustom == nil {
		return false
	}
	values := reflect.ValueOf(serv.UserCustom)

	if fn := values.MethodByName("OnValidateFail"); fn.IsValid() && fn.Type().NumIn() == 2 {
		f := reflect.ValueOf(&serv.onValidateFail).Elem()
		f.Set(fn)
	} else {
		serv.onValidateFail = serv.validateFail
	}

	if fn := values.MethodByName("UnsupportedRequest"); fn.IsValid() && fn.Type().NumIn() == 2 {
		f := reflect.ValueOf(&serv.unsupported).Elem()
		f.Set(fn)
	} else {
		serv.unsupported = serv.unsupportedRequest
	}

	if fn := values.MethodByName("OnRequestError"); fn.IsValid() && fn.Type().NumIn() == 2 {
		f := reflect.ValueOf(&serv.onRequestError).Elem()
		f.Set(fn)
	} else {
		serv.onRequestError = serv.requestError
	}

	if fn := values.MethodByName("OnTextRequest"); fn.IsValid() && fn.Type().NumIn() == 3 {
		f := reflect.ValueOf(&serv.onTextRequest).Elem()
		f.Set(fn)
	} else {
		serv.onTextRequest = serv.postText
	}

	if fn := values.MethodByName("OnLocationRequest"); fn.IsValid() && fn.Type().NumIn() == 3 {
		f := reflect.ValueOf(&serv.onLocationRequest).Elem()
		f.Set(fn)
	} else {
		serv.onLocationRequest = serv.postLocation
	}

	if fn := values.MethodByName("OnSubscribeEvent"); fn.IsValid() && fn.Type().NumIn() == 3 {
		f := reflect.ValueOf(&serv.onSubscribeEvent).Elem()
		f.Set(fn)
	} else {
		serv.onSubscribeEvent = func(w http.ResponseWriter, info *WXRequestInfo, sub bool) {
			w.WriteHeader(200)
		}
	}

	if fn := values.MethodByName("OnQRScanEvent"); fn.IsValid() && fn.Type().NumIn() == 4 {
		f := reflect.ValueOf(&serv.onQRScanEvent).Elem()
		f.Set(fn)
	} else {
		serv.onQRScanEvent = func(w http.ResponseWriter, info *WXRequestInfo, key, ticket string) {
			info.ResponseText(w, "很抱歉，对您的请求暂时无法响应")
		}

	}

	if fn := values.MethodByName("OnLocationEvent"); fn.IsValid() && fn.Type().NumIn() == 3 {
		f := reflect.ValueOf(&serv.onLocationEvent).Elem()
		f.Set(fn)
	} else {
		serv.onLocationEvent = func(w http.ResponseWriter, info *WXRequestInfo, pos *WXLocationEvent) {
			info.ResponseText(w, "很抱歉，对您的请求暂时无法响应")
		}
	}

	if fn := values.MethodByName("OnMenuEvent"); fn.IsValid() && fn.Type().NumIn() == 3 {
		f := reflect.ValueOf(&serv.onMenuEvent).Elem()
		f.Set(fn)
	} else {
		serv.onMenuEvent = func(w http.ResponseWriter, info *WXRequestInfo, key string) {
			info.ResponseText(w, "很抱歉，对您的请求暂时无法响应")
		}
	}

	if fn := values.MethodByName("OnLinkEvent"); fn.IsValid() && fn.Type().NumIn() == 3 {
		f := reflect.ValueOf(&serv.onLinkEvent).Elem()
		f.Set(fn)
	} else {
		serv.onLinkEvent = func(w http.ResponseWriter, info *WXRequestInfo, url string) {
			info.ResponseText(w, "很抱歉，对您的请求暂时无法响应")
		}
	}
	go serv.get_access_token()
	return true
}

type access_token_json struct {
	AccessToken string `json:"access_token"`
	Expire      int    `json:"expires_in"`
}

func (serv *Weixin) get_access_token() {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", serv.AppID, serv.AppSecret)
	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		d := access_token_json{}
		if err = json.Unmarshal(body, &d); err == nil {
			serv.AccessToken = d.AccessToken
			time.Sleep(time.Duration(d.Expire-100) * time.Second)
			go serv.get_access_token()
			return
		}
	}
	time.Sleep(10 * time.Second)
	go serv.get_access_token()
}

func (serv *Weixin) onGet(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	if _, ok := v["echostr"]; ok {
		go_logger.GetLogger("weixin").Info("服务器来验证了")
		w.Write([]byte(v.Get("echostr")))
	} else {
		// 这是个什么请求？只有验证，没有echostr
		w.WriteHeader(200)
	}
}

func (serv *Weixin) onPost(w http.ResponseWriter, r *http.Request) {
	go_logger.GetLogger("weixin").Trace("收到一次POST请求")
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	req := WXRequest{}
	go_logger.GetLogger("weixin").Trace("内容是:", body)
	if err := xml.Unmarshal(body, &req); err != nil {
		go_logger.GetLogger("weixin").Error("解析xml出错:", err.Error())
		serv.onRequestError(w, r)
		return
	}
	if strings.ToLower(req.MsgType) == "text" {
		serv.onTextRequest(w, &req.WXRequestInfo, &req.Content)
	} else if strings.ToLower(req.MsgType) == "location" {
		serv.onLocationRequest(w, &req.WXRequestInfo, &req.WXLocationRequest)
	} else if strings.ToLower(req.MsgType) == "image" {
		go_logger.GetLogger("weixin").Trace("用户发上来一张图片: ", req.MediaId)
		go_logger.GetLogger("weixin").Trace(req.PicUrl)
	} else if strings.ToLower(req.MsgType) == "voice" {
		go_logger.GetLogger("weixin").Trace("用户发上来一段语音: ", req.MediaId)
		go_logger.GetLogger("weixin").Trace(req.Format)
	} else if strings.ToLower(req.MsgType) == "event" {
		serv.onEvent(w, &req.WXRequestInfo, &req.WXEvent)
	} else {
		go_logger.GetLogger("weixin").Error("不支持的POST请求:", req.MsgType)
		serv.unsupported(w, &req.WXRequestInfo)
		fmt.Println(body)
	}
}

func (serv *Weixin) onEvent(w http.ResponseWriter, info *WXRequestInfo, e *WXEvent) {
	if strings.ToLower(e.Event) == "subscribe" {
		serv.onSubscribeEvent(w, info, true)
		if len(e.EventKey) > 0 {
			serv.onQRScanEvent(w, info, e.EventKey, e.Ticket)
		}
	} else if strings.ToLower(e.Event) == "unsubscribe" {
		serv.onSubscribeEvent(w, info, false)
	} else if strings.ToLower(e.Event) == "scan" {
		serv.onQRScanEvent(w, info, e.EventKey, e.Ticket)
	} else if strings.ToLower(e.Event) == "location" {
		serv.onLocationEvent(w, info, &e.WXLocationEvent)
	} else if strings.ToLower(e.Event) == "click" {
		serv.onMenuEvent(w, info, e.EventKey)
	} else if strings.ToLower(e.Event) == "view" {
		serv.onLinkEvent(w, info, e.EventKey)
	} else {
		go_logger.GetLogger("weixin").Trace("不支持的Event: ", e.Event)
	}
}

func (serv *Weixin) validateFail(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
}

func (serv *Weixin) unsupportedRequest(w http.ResponseWriter, info *WXRequestInfo) {
	info.ResponseText(w, "很抱歉，对您的请求暂时无法响应")
}

func (serv *Weixin) requestError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)
}

func (serv *Weixin) postText(w http.ResponseWriter, info *WXRequestInfo, content *string) {
	go_logger.GetLogger("weixin").Trace("用户发送了一个文本:", content)
	serv.unsupported(w, info)
}

func (serv *Weixin) postLocation(w http.ResponseWriter, info *WXRequestInfo, pos *WXLocationRequest) {
	go_logger.GetLogger("weixin").Trace("用户发送了一个坐标:", pos.Location_Y, ",", pos.Location_X)
	serv.unsupported(w, info)
}
