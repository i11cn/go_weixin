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

func (serv *Weixin) resp_unsupport(in []reflect.Value) []reflect.Value {
	if w, ok := in[0].Interface().(http.ResponseWriter); ok {
		if info, ok := in[1].Interface().(*WXRequestInfo); ok {
			info.ResponseText(w, "不支持的功能")
		}
	}
	return []reflect.Value{}
}

func (serv *Weixin) resp_error(in []reflect.Value) []reflect.Value {
	if w, ok := in[0].Interface().(http.ResponseWriter); ok {
		w.WriteHeader(500)
	}
	return []reflect.Value{}
}

func param_match(fn1 reflect.Value, fn2 reflect.Value) bool {
	t1, t2 := fn1.Type(), fn2.Type()
	if (t1.NumIn() != t2.NumIn()) || (t1.NumOut() != t2.NumOut()) {
		return false
	}
	for i := 0; i < t1.NumIn(); i++ {
		if t1.In(i) != t2.In(i) {
			return false
		}
	}
	return true
}

func (serv *Weixin) check_and_set(l_func reflect.Value, uc reflect.Value, name string, resp func([]reflect.Value) []reflect.Value) {
	f := l_func.Elem()
	if fn := uc.MethodByName(name); fn.IsValid() && param_match(f, fn) {
		f.Set(fn)
		return
	}
	v := reflect.MakeFunc(f.Type(), resp)
	f.Set(v)
}

func (serv *Weixin) init() bool {
	if serv.user_custom == nil {
		return false
	}
	values := reflect.ValueOf(serv.user_custom)

	serv.check_and_set(reflect.ValueOf(&serv.onValidateFail), values, "OnValidateFail", serv.resp_error)
	serv.check_and_set(reflect.ValueOf(&serv.unsupported), values, "UnsupportedRequest", serv.resp_unsupport)
	serv.check_and_set(reflect.ValueOf(&serv.onRequestError), values, "OnRequestError", serv.resp_error)
	serv.check_and_set(reflect.ValueOf(&serv.onTextRequest), values, "OnTextRequest", serv.resp_unsupport)
	serv.check_and_set(reflect.ValueOf(&serv.onImageRequest), values, "OnImageRequest", serv.resp_unsupport)
	serv.check_and_set(reflect.ValueOf(&serv.onVoiceRequest), values, "OnVoiceRequest", serv.resp_unsupport)
	serv.check_and_set(reflect.ValueOf(&serv.onVideoRequest), values, "OnVideoRequest", serv.resp_unsupport)
	serv.check_and_set(reflect.ValueOf(&serv.onLocationRequest), values, "OnLocationRequest", serv.resp_unsupport)
	serv.check_and_set(reflect.ValueOf(&serv.onLinkRequest), values, "OnLinkRequest", serv.resp_unsupport)
	serv.check_and_set(reflect.ValueOf(&serv.onSubscribeEvent), values, "OnSubscribeEvent", serv.resp_unsupport)
	serv.check_and_set(reflect.ValueOf(&serv.onQRScanEvent), values, "OnQRScanEvent", serv.resp_unsupport)
	serv.check_and_set(reflect.ValueOf(&serv.onLocationEvent), values, "OnLocationEvent", serv.resp_unsupport)
	serv.check_and_set(reflect.ValueOf(&serv.onMenuEvent), values, "OnMenuEvent", serv.resp_unsupport)
	serv.check_and_set(reflect.ValueOf(&serv.onLinkEvent), values, "OnLinkEvent", serv.resp_unsupport)

	go serv.get_access_token()
	return true
}

type access_token_json struct {
	AccessToken string `json:"access_token"`
	Expire      int    `json:"expires_in"`
}

type js_token_json struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Ticket  string `json:"ticket"`
	Expire  int    `json:"expires_in"`
}

func (serv *Weixin) get_access_token() {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", serv.AppID, serv.AppSecret)
	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		d := access_token_json{}
		if err = json.Unmarshal(body, &d); err == nil {
			serv.AccessToken = d.AccessToken
			go serv.get_jsapi_token(1)
			go serv.get_wxcard_token(1)
			time.Sleep(time.Duration(d.Expire-100) * time.Second)
			go serv.get_access_token()
			return
		}
	}
	time.Sleep(10 * time.Second)
	go serv.get_access_token()
}

func (serv *Weixin) get_jsapi_token(count int) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi", serv.AccessToken)
	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		d := js_token_json{}
		if err = json.Unmarshal(body, &d); err == nil {
			serv.JSApiToken = d.Ticket
			time.Sleep(time.Duration(d.Expire) * time.Second)
			go func(token string) {
				if serv.JSApiToken == token {
					serv.JSApiToken = ""
				}
			}(serv.JSApiToken)
			return
		}
	}
	if count <= 3 {
		time.Sleep(10 * time.Second)
		go serv.get_jsapi_token(count + 1)
	}
}

func (serv *Weixin) get_wxcard_token(count int) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=wx_card", serv.AccessToken)
	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		d := js_token_json{}
		if err = json.Unmarshal(body, &d); err == nil {
			serv.WXCardToken = d.Ticket
			time.Sleep(time.Duration(d.Expire) * time.Second)
			go func(token string) {
				if serv.WXCardToken == token {
					serv.WXCardToken = ""
				}
			}(serv.WXCardToken)
			return
		}
	}
	if count <= 3 {
		time.Sleep(10 * time.Second)
		go serv.get_jsapi_token(count + 1)
	}
}

func (serv *Weixin) gen_app_info(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-javascript")
	s := fmt.Sprintf("var app_info = {}")
	w.Write([]byte(s))
}

func (serv *Weixin) onGet(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	if _, ok := v["echostr"]; ok {
		logger.GetLogger("weixin").Info("服务器来验证了")
		w.Write([]byte(v.Get("echostr")))
	} else {
		// 这是个什么请求？只有验证，没有echostr
		w.WriteHeader(200)
	}
}

func (serv *Weixin) onPost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	req := WXRequest{}
	logger.GetLogger("weixin").Trace("收到一次POST请求, 内容是:")
	logger.GetLogger("weixin").Trace(string(body))
	if err := xml.Unmarshal(body, &req); err != nil {
		logger.GetLogger("weixin").Error("解析xml出错:", err.Error())
		serv.onRequestError(w, r)
		return
	}
	if strings.ToLower(req.MsgType) == "text" {
		logger.GetLogger("weixin").Info("用户发送了一个文本请求")
		serv.onTextRequest(w, &req.WXRequestInfo, req.Content)
	} else if strings.ToLower(req.MsgType) == "location" {
		logger.GetLogger("weixin").Info("用户发送了一个位置请求")
		serv.onLocationRequest(w, &req.WXRequestInfo, &req.WXLocationRequest)
	} else if strings.ToLower(req.MsgType) == "image" {
		serv.unsupported(w, &req.WXRequestInfo)
		logger.GetLogger("weixin").Info("用户发上来一张图片: ", req.MediaId)
		serv.onImageRequest(w, &req.WXRequestInfo, req.MediaId, req.PicUrl)
	} else if strings.ToLower(req.MsgType) == "voice" {
		logger.GetLogger("weixin").Info("用户发上来一段语音: ", req.MediaId)
		serv.onVoiceRequest(w, &req.WXRequestInfo, req.MediaId, req.Format)
	} else if strings.ToLower(req.MsgType) == "video" {
		logger.GetLogger("weixin").Info("用户发上来一段视频: ", req.MediaId)
		serv.onVideoRequest(w, &req.WXRequestInfo, req.MediaId, req.ThumbMediaId, false)
	} else if strings.ToLower(req.MsgType) == "shortvideo" {
		logger.GetLogger("weixin").Info("用户发上来一段小视频: ", req.MediaId)
		serv.onVideoRequest(w, &req.WXRequestInfo, req.MediaId, req.ThumbMediaId, true)
	} else if strings.ToLower(req.MsgType) == "link" {
		logger.GetLogger("weixin").Info("用户发上来一个链接")
		serv.onLinkRequest(w, &req.WXRequestInfo, &req.WXLinkRequest)
	} else if strings.ToLower(req.MsgType) == "event" {
		logger.GetLogger("weixin").Info("用户发上来一个Event")
		serv.onEvent(w, &req.WXRequestInfo, &req.WXEvent)
	} else if strings.ToLower(req.MsgType) == "event" {
		serv.onEvent(w, &req.WXRequestInfo, &req.WXEvent)
	} else {
		logger.GetLogger("weixin").Error("不支持的POST请求:", req.MsgType)
		serv.unsupported(w, &req.WXRequestInfo)
	}
}

func (serv *Weixin) onEvent(w http.ResponseWriter, info *WXRequestInfo, e *WXEvent) {
	if strings.ToLower(e.Event) == "subscribe" {
		logger.GetLogger("weixin").Info("用户订阅了我们的号")
		serv.onSubscribeEvent(w, info, true)
		if len(e.EventKey) > 0 {
			serv.onQRScanEvent(w, info, e.EventKey, e.Ticket)
		}
	} else if strings.ToLower(e.Event) == "unsubscribe" {
		logger.GetLogger("weixin").Info("用户取消订阅")
		serv.onSubscribeEvent(w, info, false)
	} else if strings.ToLower(e.Event) == "scan" {
		logger.GetLogger("weixin").Info("用户触发了扫描二维码事件")
		serv.onQRScanEvent(w, info, e.EventKey, e.Ticket)
		logger.GetLogger("weixin").Info("用户触发了位置事件")
	} else if strings.ToLower(e.Event) == "location" {
		serv.onLocationEvent(w, info, &e.WXLocationEvent)
	} else if strings.ToLower(e.Event) == "click" {
		logger.GetLogger("weixin").Info("用户触发了点击菜单事件")
		serv.onMenuEvent(w, info, e.EventKey)
	} else if strings.ToLower(e.Event) == "view" {
		logger.GetLogger("weixin").Info("用户触发了点击菜单链接事件")
		serv.onLinkEvent(w, info, e.EventKey)
	} else {
		logger.GetLogger("weixin").Trace("不支持的Event: ", e.Event)
	}
}
