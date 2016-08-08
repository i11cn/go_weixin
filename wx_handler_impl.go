package weixin

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"
)

type (
	pre_process struct {
		ToUserName   string
		FromUserName string
		CreateTime   int64
		MsgType      string
		Event        string
	}
)

func dummy_verify(string, string, string, string) bool {
	return true
}

func dummy_text_req(string, time.Time, int64, string) (interface{}, error) {
	return nil, nil
}

func dummy_image_req(string, time.Time, int64, string, string) (interface{}, error) {
	return nil, nil
}

func dummy_voice_req(string, time.Time, int64, string, string, string) (interface{}, error) {
	return nil, nil
}

func dummy_video_req(string, time.Time, int64, string, string) (interface{}, error) {
	return nil, nil
}

func dummy_short_video_req(string, time.Time, int64, string, string) (interface{}, error) {
	return nil, nil
}

func dummy_location_req(string, time.Time, int64, float64, float64, int, string) (interface{}, error) {
	return nil, nil
}

func dummy_link_req(string, time.Time, int64, string, string, string) (interface{}, error) {
	return nil, nil
}

func dummy_subscribe(string, time.Time, string, string) (interface{}, error) {
	return nil, nil
}

func dummy_unsubscribe(string, time.Time) (interface{}, error) {
	return nil, nil
}

func dummy_scan(string, time.Time, string, string) (interface{}, error) {
	return nil, nil
}

func dummy_location(string, time.Time, float64, float64, float64) (interface{}, error) {
	return nil, nil
}

func dummy_menu_click(string, time.Time, string) (interface{}, error) {
	return nil, nil
}

func dummy_menu_view(string, time.Time, string) (interface{}, error) {
	return nil, nil
}

func (wh *WXHandler) handle(fn interface{}) error {
	switch reflect.TypeOf(fn) {
	case reflect.TypeOf(dummy_text_req):
	case reflect.TypeOf(dummy_image_req):
	case reflect.TypeOf(dummy_voice_req):
	case reflect.TypeOf(dummy_video_req):
	case reflect.TypeOf(dummy_short_video_req):
	case reflect.TypeOf(dummy_location_req):
	case reflect.TypeOf(dummy_link_req):
	case reflect.TypeOf(dummy_subscribe):
	case reflect.TypeOf(dummy_unsubscribe):
	case reflect.TypeOf(dummy_scan):
	case reflect.TypeOf(dummy_location):
	case reflect.TypeOf(dummy_menu_click):
	case reflect.TypeOf(dummy_menu_view):
	default:
		return errors.New("注册的方法类型不正确")
	}
	return nil

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
	req := pre_process{}
	if err = xml.Unmarshal(body, &req); err != nil {
		wh.Logger.Error("解析xml出错:", err.Error())
		w.WriteHeader(500)
		return
	}
	wh.Logger.Trace(req)
}
