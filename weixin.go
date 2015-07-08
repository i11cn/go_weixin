package weixin

import (
	"encoding/xml"
	"github.com/i11cn/go_logger"
	"net/http"
	"time"
)

type (
	WXRequestInfo struct {
		ToUserName   string
		FromUserName string
		CreateTime   time.Duration
		MsgId        int64
	}

	WXLocationRequest struct {
		Location_Y float64
		Location_X float64
		Scale      int
		Label      string
	}

	WXLinkRequest struct {
		Title       string
		Description string
		Url         string
	}

	WXLocationEvent struct {
		Longitude float64
		Latitude  float64
		Precision float64
	}

	WXEvent struct {
		Event    string
		EventKey string
		Ticket   string
		WXLocationEvent
	}

	WXRequest struct {
		WXRequestInfo
		MsgType string
		Content string
		WXLocationRequest
		WXLinkRequest
		MediaId      string
		PicUrl       string
		Format       string
		ThumbMediaId string
		WXEvent
	}

	WXResponseInfo struct {
		XMLName      xml.Name `xml:"xml"`
		ToUserName   string
		FromUserName string
		CreateTime   int64
		MsgType      string
	}

	WXTextResponse struct {
		WXResponseInfo
		Content string
	}

	WXImageResponse struct {
		WXResponseInfo
		MediaId string `xml:"Image>MediaId"`
	}

	WXVoiceResponse struct {
		WXResponseInfo
		MediaId string `xml:"Voice>MediaId"`
	}

	WXVideoResponse struct {
		WXResponseInfo
		MediaId     string `xml:"Video>MediaId"`
		Title       string `xml:"Video>Title"`
		Description string `xml:"Video>Description"`
	}

	WXMusicResponse struct {
		WXResponseInfo
		Title        string `xml:"Music>Title"`
		Description  string `xml:"Music>Description"`
		MusicUrl     string `xml:"Music>MusicUrl"`
		HQMusicUrl   string `xml:"Music>HQMusicUrl"`
		ThumbMediaId string `xml:"Music>ThumbMediaId"`
	}

	WXNewsItem struct {
		Title       string
		Description string
		PicUrl      string
		Url         string
	}

	WXNewsResponse struct {
		WXResponseInfo
		ArticleCount int
		Items        []WXNewsItem `xml:"Articles>item"`
	}

	WXConfig struct {
		Token     string
		MsgKey    string
		AppID     string
		AppSecret string
	}
)

type (
	OnValidateFail     func(w http.ResponseWriter, r *http.Request)
	UnsupportedRequest func(w http.ResponseWriter, info *WXRequestInfo)
	OnRequestError     func(w http.ResponseWriter, r *http.Request)

	OnTextRequest     func(w http.ResponseWriter, info *WXRequestInfo, content string)
	OnLocationRequest func(w http.ResponseWriter, info *WXRequestInfo, pos *WXLocationRequest)
	OnImageRequest    func(w http.ResponseWriter, info *WXRequestInfo, mid, url string)
	OnVoiceRequest    func(w http.ResponseWriter, info *WXRequestInfo, mid, format string)
	OnVideoRequest    func(w http.ResponseWriter, info *WXRequestInfo, mid, thumb string, short bool)
	OnLinkRequest     func(w http.ResponseWriter, info *WXRequestInfo, req *WXLinkRequest)

	OnSubscribeEvent func(w http.ResponseWriter, info *WXRequestInfo, sub bool)
	OnQRScanEvent    func(w http.ResponseWriter, info *WXRequestInfo, key, ticket string)
	OnLocationEvent  func(w http.ResponseWriter, info *WXRequestInfo, pos *WXLocationEvent)
	OnMenuEvent      func(w http.ResponseWriter, info *WXRequestInfo, key string)
	OnLinkEvent      func(w http.ResponseWriter, info *WXRequestInfo, url string)
)

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

func (info *WXRequestInfo) ResponseImage(w http.ResponseWriter, id string) {
	resp := WXImageResponse{}
	resp.ToUserName = info.FromUserName
	resp.FromUserName = info.ToUserName
	resp.CreateTime = time.Now().Unix()
	resp.MsgType = "image"
	resp.MediaId = id
	info.Response(w, resp)
}

func (info *WXRequestInfo) ResponseVoice(w http.ResponseWriter, id string) {
	resp := WXVoiceResponse{}
	resp.ToUserName = info.FromUserName
	resp.FromUserName = info.ToUserName
	resp.CreateTime = time.Now().Unix()
	resp.MsgType = "voice"
	resp.MediaId = id
	info.Response(w, resp)
}

func (info *WXRequestInfo) ResponseVideo(w http.ResponseWriter, id, title, desc string) {
	resp := WXVideoResponse{}
	resp.ToUserName = info.FromUserName
	resp.FromUserName = info.ToUserName
	resp.CreateTime = time.Now().Unix()
	resp.MsgType = "video"
	resp.MediaId = id
	resp.Title = title
	resp.Description = desc
	info.Response(w, resp)
}

func (info *WXRequestInfo) ResponseMusic(w http.ResponseWriter, id, title, desc, url, hqurl string) {
	resp := WXMusicResponse{}
	resp.ToUserName = info.FromUserName
	resp.FromUserName = info.ToUserName
	resp.CreateTime = time.Now().Unix()
	resp.MsgType = "music"
	resp.ThumbMediaId = id
	resp.Title = title
	resp.Description = desc
	resp.MusicUrl = url
	resp.HQMusicUrl = hqurl
	info.Response(w, resp)
}

func (info *WXRequestInfo) ResponseNews(w http.ResponseWriter, articles []WXNewsItem) {
	resp := WXNewsResponse{}
	resp.ToUserName = info.FromUserName
	resp.FromUserName = info.ToUserName
	resp.CreateTime = time.Now().Unix()
	resp.MsgType = "news"
	resp.ArticleCount = len(articles)
	resp.Items = articles
	info.Response(w, resp)
}

func NewWeixinServ(conf *WXConfig, uc interface{}) *Weixin {
	serv := &Weixin{WXConfig: *conf, user_custom: uc}
	if serv.init() {
		return serv
	} else {
		return nil
	}
}
