package weixin

type (
	TplVal struct {
		Name  string `json:"-"`
		Value string `json:"value,omitempty"`
		Color string `json:"color,omitempty"`
	}

	TemplateMgr interface {
		SetIndustry(primary, secondary int) error
		GetIndustry() (int, int, error)
		GetTemplateID(short string) (string, error)
		GetTemplateList() error
		DeleteTemplate(id string) error
		SendMessage(user, tpl_id, url string, vals []TplVal) (int64, error)
	}

	TemplateSendJob struct {
		EventRequest
		MsgID  int64
		Status string
	}
)
