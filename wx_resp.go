package weixin

import (
	"encoding/xml"
	"time"
)

type (
	Response struct {
		XMLName      xml.Name `xml:"xml"`
		ToUserName   string
		FromUserName string
		CreateTime   int64
		MsgType      string
	}
	TextResp struct {
		Response
		Content string
	}
	ImageResp struct {
		Response
		MediaId string `xml:"image>MediaId"`
	}
	VoiceResp struct {
		Response
		MediaId string `xml:"Voice>MediaId"`
	}
	VideoResp struct {
		Response
		MediaId     string `xml:"Video>MediaId"`
		Title       string `xml:"Video>Title,omitempty"`
		Description string `xml:"Video>Description,omitempty"`
	}
	MusicResp struct {
		Response
		Title        string `xml:"Music>Title,omitempty"`
		Description  string `xml:"Music>Description,omitempty"`
		MusicUrl     string `xml:"Music>MusicUrl,omitempty"`
		HQMusicUrl   string `xml:"Music>HQMusicUrl,omitempty"`
		ThumbMediaId string `xml:"Music>ThumbMediaId,omitempty"`
	}
	Article struct {
		Title       string `xml:"omitempty"`
		Description string `xml:"omitempty"`
		PicUrl      string `xml:"omitempty"`
		Url         string `xml:"omitempty"`
	}
	ArticleResp struct {
		Response
		ArticleCount int
		Articles     []Article `xml:"Articles>item,omitempty"`
	}
)

var (
	g_mp_id string = ""
)

func new_response(user, msg_type string) Response {
	return Response{ToUserName: user, FromUserName: g_mp_id, CreateTime: time.Now().Unix(), MsgType: msg_type}
}

func NewTextResponse(user, msg string) *TextResp {
	return &TextResp{new_response(user, "text"), msg}
}

func NewImageResponse(user, img_id string) *ImageResp {
	return &ImageResp{new_response(user, "image"), img_id}
}

func NewVoiceResponse(user, voice_id string) *VoiceResp {
	return &VoiceResp{new_response(user, "voice"), voice_id}
}

func NewVideoResponse(user, video_id, title, desc string) *VideoResp {
	return &VideoResp{new_response(user, "video"), video_id, title, desc}
}

func NewMusicResponse(user, music_url, hq_url string) *MusicResp {
	return &MusicResp{Response: new_response(user, "music"), MusicUrl: music_url, HQMusicUrl: hq_url}
}

func NewArticleResp(user string, arts ...Article) *ArticleResp {
	if len(arts) < 1 {
		return nil
	}
	return &ArticleResp{new_response(user, "news"), len(arts), arts}
}
