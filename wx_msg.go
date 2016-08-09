package weixin

import (
	"time"
)

type (
	MsgRequest struct {
		ToUserName   string
		FromUserName string
		CreateTime   int64
		MsgType      string
		MsgId        int64
	}
	TextReq struct {
		MsgRequest
		Content string
	}
	ImageReq struct {
		MsgRequest
		PicUrl  string
		MediaId string
	}
	VoiceReq struct {
		MsgRequest
		MediaId     string
		Format      string
		Recognition string
	}
	VideoReq struct {
		MsgRequest
		MediaId      string
		ThumbMediaId string
	}
	ShortVideoReq struct {
		MsgRequest
		MediaId      string
		ThumbMediaId string
	}
	LocationReq struct {
		MsgRequest
		Location_X float64
		Location_Y float64
		Scale      int
		Label      string
	}
	LinkReq struct {
		MsgRequest
		Url         string
		Title       string
		Description string
	}
)

func (r *MsgRequest) UserName() string {
	return r.FromUserName
}

func (r *MsgRequest) GetCreateTime() time.Time {
	return time.Unix(r.CreateTime, 0)
}

func (r *MsgRequest) ID() int64 {
	return r.MsgId
}
