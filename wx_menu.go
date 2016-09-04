package weixin

type (
	MenuItem struct {
		Name    string     `json:"name"`
		Key     string     `json:"key"`
		Type    string     `json:"type,omitempty"`
		Url     string     `json:"url,omitempty"`
		MediaId string     `json:"media_id,omitempty"`
		SubMenu []MenuItem `json:"sub_button,omitempty"`
	}
	MatchRule struct {
		Group    string `json:"group_id,omitempty"`
		Sex      string `json:"sex,omitempty"`
		Country  string `json:"country,omitempty"`
		Province string `json:"province,omitempty"`
		City     string `json:"city,omitempty"`
		Client   string `json:"client_platform_type,omitempty"`
		Language string `json:"language,omitempty"`
	}
	Menu struct {
		Items []MenuItem `json:"button"`
		Rule  *MatchRule `json:"matchrule,omitempty"`
	}

	MenuMgr interface {
		SubmitMenu(*Menu) error
		QueryMenu() (*Menu, error)
		DeleteMenu() error
		I17NMenu() error
	}
)

func (i *MenuItem) IsValid() bool {
	if len(i.Name) > 16 {
		return false
	}
	if len(i.SubMenu) > 5 {
		return false
	}
	for _, s := range i.SubMenu {
		if len(s.Name) > 40 {
			return false
		}
		if len(s.SubMenu) > 0 {
			return false
		}
	}
	return true
}

func (i *MenuItem) AddItem(mi *MenuItem) *MenuItem {
	i.Type = ""
	i.SubMenu = append(i.SubMenu, *mi)
	return i
}

func (r *MatchRule) IsValid() bool {
	if len(r.Group) < 1 && len(r.Sex) < 1 && len(r.Client) < 1 && len(r.Country) < 1 && len(r.Province) < 1 && len(r.City) < 1 && len(r.Language) < 1 {
		return false
	}
	if len(r.Country) < 1 && (len(r.Province) > 0 || len(r.City) > 0) {
		return false
	}
	if len(r.Province) < 1 && len(r.City) > 0 {
		return false
	}
	return true
}

func (m *Menu) IsValid() bool {
	if len(m.Items) > 3 {
		return false
	}
	for _, i := range m.Items {
		if !i.IsValid() {
			return false
		}
	}
	if m.Rule != nil {
		return m.Rule.IsValid()
	}
	return true
}

func (m *Menu) AddItem(i *MenuItem) *Menu {
	m.Items = append(m.Items, *i)
	return m
}
