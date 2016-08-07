package weixin

type (
	EventRequest struct {
		ToUserName   string
		FromUserName string
		CreateTime   int
		MsgType      string
		Event        string
	}
	Subscribe struct {
		EventRequest
		EventKey *string
		Ticket   *string
	}
	Unsubscribe struct {
		EventRequest
	}
	Scan struct {
		EventRequest
		EventKey uint32
		Ticket   string
	}
	Location struct {
		EventRequest
		Latitude  float64
		Longitude float64
		Precision float64
	}
	MenuClick struct {
		EventRequest
		Key string `xml:"EventKey"`
	}
	MenuView struct {
		EventRequest
		Url string `xml:"EventKey"`
	}
)
