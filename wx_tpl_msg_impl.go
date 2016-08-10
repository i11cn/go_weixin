package weixin

type (
	tpl_mgr_impl struct {
		client  *WXClient
		client2 *WXClient
	}
)

func (tm *tpl_mgr_impl) SetIndustry(int, int) error {
	return nil
}

func (tm *tpl_mgr_impl) GetIndustry() (int, int, error) {
	return 0, 0, nil
}

type (
	get_tpl_id_req struct {
		ShortID string `json:"template_id_short"`
	}
	get_tpl_id_resp struct {
		ErrorCode int    `json:"errcode"`
		ErrorMsg  string `json:"errmsg"`
		TplID     string `json:"template_id"`
	}
)

func (tm *tpl_mgr_impl) GetTemplateID(short string) (string, error) {
	req := &get_tpl_id_req{short}
	resp := &get_tpl_id_resp{}
	_, err := tm.client.Post("/api_add_template", req, resp)
	if err != nil {
		return "", nil
	} else {
		return resp.TplID, nil
	}
}

func (tm *tpl_mgr_impl) GetTemplateList() error {
	return nil
}

func (tm *tpl_mgr_impl) DeleteTemplate(string) error {
	return nil
}

type (
	send_tpl_req struct {
		User   string            `json:"touser"`
		TplID  string            `json:"template_id"`
		Url    string            `json:"url,omitempty"`
		Values map[string]TplVal `json:"data"`
	}
	send_tpl_resp struct {
		ErrorCode int    `json:"errcode"`
		ErrorMsg  string `json:"errmsg"`
		MsgID     int64  `json:"msgid"`
	}
)

func (tm *tpl_mgr_impl) SendMessage(user string, tpl_id string, url string, vals []TplVal) (int64, error) {
	req := &send_tpl_req{User: user, TplID: tpl_id, Url: url}
	req.Values = make(map[string]TplVal)
	for _, v := range vals {
		req.Values[v.Name] = v
	}
	resp := &send_tpl_resp{}
	_, err := tm.client2.Post("send", req, resp)
	if err != nil {
		return 0, err
	}
	return resp.MsgID, nil
}
