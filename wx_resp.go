package weixin

type (
	Response struct {
		ToUserName   string
		FromUserName string
		CreateTime   int
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
