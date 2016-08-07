package weixin

type (
	MsgRequest struct {
		ToUserName   string
		FromUserName string
		CreateTime   int
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
		Title       string
		Description string
		Url         string
	}
)
