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
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	req := WXRequest{}
	if err := xml.Unmarshal(body, &req); err != nil {
		go_logger.GetLogger("weixin").Error("解析xml出错:", err.Error())
		serv.onRequestError(w, r)
		return
	}
	if strings.ToLower(req.MsgType) == "text" {
		serv.onTextRequest(w, &req.WXRequestInfo, &req.WXTextRequest)
	} else if strings.ToLower(req.MsgType) == "location" {
		serv.onLocationRequest(w, &req.WXRequestInfo, &req.WXLocationRequest)
	} else {
		go_logger.GetLogger("weixin").Error("不支持的POST请求:", req.MsgType)
		serv.unsupported(w, &req.WXRequestInfo)
		fmt.Println(body)
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

func (serv *Weixin) postText(w http.ResponseWriter, info *WXRequestInfo, text *WXTextRequest) {
	go_logger.GetLogger("weixin").Trace("用户发送了一个文本:", text.Content)
	serv.unsupported(w, info)
}

func (serv *Weixin) postLocation(w http.ResponseWriter, info *WXRequestInfo, pos *WXLocationRequest) {
	go_logger.GetLogger("weixin").Trace("用户发送了一个坐标:", pos.Location_X, ",", pos.Location_Y)
	serv.unsupported(w, info)
}
