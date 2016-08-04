package weixin

import ()

type (
	WXTokenMgr interface {
		GetAccessToken() WXToken
		SetGlobalInfo(*WXGlobalInfo)
	}

	TokenStorage interface {
		SetToken(string)
		GetToken() string
		Start()
		Stop()
	}

	WXToken interface {
		Expired()
		GetToken() string
		Close()
		Primary(bool)
		SetSource(TokenStorage)
		SetGlobalInfo(*WXGlobalInfo)
	}
)

func DefaultTokenMgr(info *WXGlobalInfo) WXTokenMgr {
	ret := new_default_token_mgr()
	ret.SetGlobalInfo(info)
	return ret
}
