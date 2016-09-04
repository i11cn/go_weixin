package weixin

import (
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
	ScanHandle        func(user string, t time.Time, key uint32, ticker string) (interface{}, error)
	LocationHandle    func(user string, t time.Time, lat, long, precision float64) (interface{}, error)
	MenuClickHandle   func(user string, t time.Time, key string) (interface{}, error)
	MenuViewHandle    func(user string, t time.Time, url string) (interface{}, error)

	TplMessageHandle func(user string, t time.Time, msg_id int64, status string)

	WXHandler struct {
		WXComponent

		verify_handle          VerifyHandle
		text_req_handle        TextReqHandle
		image_req_handle       ImageReqHandle
		voice_req_handle       VoiceReqHandle
		video_req_handle       VideoReqHandle
		short_video_req_handle ShortVideoReqHandle
		location_req_handle    LocationReqHandle
		link_req_handle        LinkReqHandle

		subscribe_handle   SubscribeHandle
		unsubscribe_handle UnsubscribeHandle
		scan_handle        ScanHandle
		location_handle    LocationHandle
		menu_click_handle  MenuClickHandle
		menu_view_handle   MenuViewHandle

		tpl_msg_handle TplMessageHandle
	}
)

func NewHandler(info *WXGlobalInfo, wx *Weixin) *WXHandler {
	ret := &WXHandler{WXComponent: WXComponent{info, wx}}
	ret.verify_handle = func(token, sign, ts, nonce string) bool {
		return check_sign(token, sign, ts, nonce, ret.Logger)
	}
	return ret
}

func (wh *WXHandler) Handle(fn interface{}) error {
	return wh.handle(fn)
}
