package weixin

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type (
	pre_process struct {
		ToUserName string
		MsgType    string
		Event      string
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

func dummy_scan(string, time.Time, uint32, string) (interface{}, error) {
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
		wh.text_req_handle = fn.(func(string, time.Time, int64, string) (interface{}, error))
	case reflect.TypeOf(dummy_image_req):
		wh.image_req_handle = fn.(func(string, time.Time, int64, string, string) (interface{}, error))
	case reflect.TypeOf(dummy_voice_req):
		wh.voice_req_handle = fn.(func(string, time.Time, int64, string, string, string) (interface{}, error))
	case reflect.TypeOf(dummy_video_req):
		wh.video_req_handle = fn.(func(string, time.Time, int64, string, string) (interface{}, error))
	case reflect.TypeOf(dummy_short_video_req):
		wh.short_video_req_handle = fn.(func(string, time.Time, int64, string, string) (interface{}, error))
	case reflect.TypeOf(dummy_location_req):
		wh.location_req_handle = fn.(func(string, time.Time, int64, float64, float64, int, string) (interface{}, error))
	case reflect.TypeOf(dummy_link_req):
		wh.link_req_handle = fn.(func(string, time.Time, int64, string, string, string) (interface{}, error))
	case reflect.TypeOf(dummy_subscribe):
		wh.subscribe_handle = fn.(func(string, time.Time, string, string) (interface{}, error))
	case reflect.TypeOf(dummy_unsubscribe):
		wh.unsubscribe_handle = fn.(func(string, time.Time) (interface{}, error))
	case reflect.TypeOf(dummy_scan):
		wh.scan_handle = fn.(func(string, time.Time, uint32, string) (interface{}, error))
	case reflect.TypeOf(dummy_location):
		wh.location_handle = fn.(func(string, time.Time, float64, float64, float64) (interface{}, error))
	case reflect.TypeOf(dummy_menu_click):
		wh.menu_click_handle = fn.(func(string, time.Time, string) (interface{}, error))
	case reflect.TypeOf(dummy_menu_view):
		wh.menu_view_handle = fn.(func(string, time.Time, string) (interface{}, error))
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
		w.WriteHeader(400)
	}
}

func (wh *WXHandler) do_request(w http.ResponseWriter, body []byte, t string) {
	wrapper := func(f interface{}, o interface{}, fn func() (interface{}, error)) {
		if reflect.ValueOf(f).IsNil() {
			wh.Logger.Error("没有Handle该请求: ", t)
			w.WriteHeader(501)
			return
		}
		if err := xml.Unmarshal(body, o); err != nil {
			wh.Logger.Error("解析xml出错: ", err.Error())
			w.WriteHeader(500)
			return
		}
		resp, err := fn()
		if err != nil {
			wh.Logger.Error("处理", t, "请求失败: ", err.Error())
			w.WriteHeader(500)
			return
		}
		d, err := xml.Marshal(resp)
		if err != nil {
			wh.Logger.Error("处理", t, "响应失败: ", err.Error())
			w.WriteHeader(500)
			return
		}
		if len(d) > 0 {
			w.Header().Set("Content-Type", "application/xml;charset=utf-8")
			w.Write(d)
		} else {
			w.Write([]byte(""))
		}
	}
	switch t {
	case "TEXT":
		o := &TextReq{}
		wrapper(wh.text_req_handle, o, func() (interface{}, error) {
			return wh.text_req_handle(o.UserName(), o.GetCreateTime(), o.ID(), o.Content)
		})
	case "IMAGE":
		o := &ImageReq{}
		wrapper(wh.image_req_handle, o, func() (interface{}, error) {
			return wh.image_req_handle(o.UserName(), o.GetCreateTime(), o.ID(), o.PicUrl, o.MediaId)
		})
	case "VOICE":
		o := &VoiceReq{}
		wrapper(wh.voice_req_handle, o, func() (interface{}, error) {
			return wh.voice_req_handle(o.UserName(), o.GetCreateTime(), o.ID(), o.MediaId, o.Format, o.Recognition)
		})
	case "VIDEO":
		o := &VideoReq{}
		wrapper(wh.video_req_handle, o, func() (interface{}, error) {
			return wh.video_req_handle(o.UserName(), o.GetCreateTime(), o.ID(), o.MediaId, o.ThumbMediaId)
		})
	case "SHORTVIDEO":
		o := &ShortVideoReq{}
		wrapper(wh.short_video_req_handle, o, func() (interface{}, error) {
			return wh.short_video_req_handle(o.UserName(), o.GetCreateTime(), o.ID(), o.MediaId, o.ThumbMediaId)
		})
	case "LOCATION":
		o := &LocationReq{}
		wrapper(wh.location_req_handle, o, func() (interface{}, error) {
			return wh.location_req_handle(o.UserName(), o.GetCreateTime(), o.ID(), o.Location_X, o.Location_Y, o.Scale, o.Label)
		})
	case "LINK":
		o := &LinkReq{}
		wrapper(wh.link_req_handle, o, func() (interface{}, error) {
			return wh.link_req_handle(o.UserName(), o.GetCreateTime(), o.ID(), o.Url, o.Title, o.Description)
		})
	default:
		wh.Logger.Error("客户端发送了不正确的请求类型：", t)
		w.WriteHeader(400)
	}
}

func (wh *WXHandler) do_event(w http.ResponseWriter, body []byte, e string) {
	wrapper := func(f interface{}, o interface{}, fn func() (interface{}, error)) {
		if reflect.ValueOf(f).IsNil() {
			wh.Logger.Error("没有Handle该事件: ", e)
			w.WriteHeader(501)
			return
		}

		if err := xml.Unmarshal(body, o); err != nil {
			wh.Logger.Error("解析xml出错: ", err.Error())
			w.WriteHeader(500)
			return
		}
		resp, err := fn()
		if err != nil {
			wh.Logger.Error("处理", e, "事件失败: ", err.Error())
			w.WriteHeader(500)
			return
		}
		d, err := xml.Marshal(resp)
		if err != nil {
			wh.Logger.Error("处理", e, "事件失败: ", err.Error())
			w.WriteHeader(500)
			return
		}
		if len(d) > 0 {
			w.Header().Set("Content-Type", "application/xml;charset=utf-8")
			w.Write(d)
		} else {
			w.Write([]byte(""))
		}
	}
	switch e {
	case "SUBSCRIBE":
		o := &Subscribe{}
		wrapper(wh.subscribe_handle, o, func() (interface{}, error) {
			return wh.subscribe_handle(o.UserName(), o.GetCreateTime(), o.EventKey, o.Ticket)
		})
	case "UNSUBSCRIBE":
		o := &Unsubscribe{}
		wrapper(wh.unsubscribe_handle, o, func() (interface{}, error) {
			return wh.unsubscribe_handle(o.UserName(), o.GetCreateTime())
		})
	case "SCAN":
		o := &Scan{}
		wrapper(wh.scan_handle, o, func() (interface{}, error) {
			return wh.scan_handle(o.UserName(), o.GetCreateTime(), o.EventKey, o.Ticket)
		})
	case "LOCATION":
		o := &Location{}
		wrapper(wh.location_handle, o, func() (interface{}, error) {
			return wh.location_handle(o.UserName(), o.GetCreateTime(), o.Latitude, o.Longitude, o.Precision)
		})
	case "CLICK":
		o := &MenuClick{}
		wrapper(wh.menu_click_handle, o, func() (interface{}, error) {
			return wh.menu_click_handle(o.UserName(), o.GetCreateTime(), o.Key)
		})
	case "VIEW":
		o := &MenuView{}
		wrapper(wh.menu_view_handle, o, func() (interface{}, error) {
			return wh.menu_view_handle(o.UserName(), o.GetCreateTime(), o.Url)
		})
	default:
		wh.Logger.Error("客户端发送了不正确的事件类型：", e)
		w.WriteHeader(400)
	}
}

func (wh *WXHandler) doPost(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		wh.Logger.Error(err.Error())
		w.WriteHeader(500)
		return
	}
	wh.Logger.Trace(string(body))
	req := pre_process{}
	if err = xml.Unmarshal(body, &req); err != nil {
		wh.Logger.Error("解析xml出错:", err.Error())
		w.WriteHeader(500)
		return
	}
	wh.Logger.Trace(req)
	t := strings.ToUpper(req.MsgType)
	if t == "EVENT" {
		wh.do_event(w, body, strings.ToUpper(req.Event))
	} else {
		wh.do_request(w, body, t)
	}
}

func (wh *WXHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wh.Logger.Info("收到一次", r.Method, "请求 : ", r.URL)
	v := r.URL.Query()
	if err := exist_all_values(v, "signature", "timestamp", "nonce"); err != nil {
		wh.Logger.Error("没有校验的签名，请求不是来自微信")
		w.WriteHeader(403)
		return
	}
	if !check_sign(wh.Config.Token, v.Get("signature"), v.Get("timestamp"), v.Get("nonce"), wh.Logger) {
		wh.Logger.Error("签名验证失败")
		w.WriteHeader(400)
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
