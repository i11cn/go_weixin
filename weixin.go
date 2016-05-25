package weixin

import (
	"github.com/i11cn/go_logger"
	rest "github.com/i11cn/go_rest_client"
	"math/rand"
	"net/http"
	"time"
)

type (
	WXConfig struct {
		Token          string
		EncodingAESKey string
		AppID          string
		AppSecret      string
	}

	WXMessage struct {
	}

	Weixin struct {
		cfg     WXConfig
		tokens  Tokens
		handler *WXHandler
		msg     *WXMessage
		rc      *rest.RestClient
		log     *logger.Logger
	}
)

var (
	g_log *logger.Logger = logger.GetLogger("weixin")
)

func init() {
	rand.Seed(time.Now().UnixNano())
	g_log.AddAppender(logger.NewConsoleAppender("[%N] %L (%f) : %M"))
	//g_log.AddAppender(logger.NewSplittedFileAppender("[%N] %L (%f) : %M", "weixin.log", 24*time.Hour))
}

func NewWeixin(cfg WXConfig) *Weixin {
	ret := &Weixin{cfg: cfg, tokens: Tokens{}, rc: rest.NewClient("api.weixin.qq.com", 0, "/cgi-bin"), log: g_log}
	ret.rc.SSL = true
	ret.tokens.AccessToken = ret.SetTokenSource(AccessToken, Local, true)
	return ret
}

func (wx *Weixin) SetLogger(log *logger.Logger) *Weixin {
	wx.log = log
	return wx
}

func (wx *Weixin) GetAccessToken() string {
	if wx.tokens.AccessToken != nil {
		return wx.tokens.AccessToken.GetToken()
	}
	return ""
}

func (wx *Weixin) GetJSApiToken() string {
	if wx.tokens.JSApiToken != nil {
		return wx.tokens.JSApiToken.GetToken()
	}
	return ""
}

func (wx *Weixin) GetWXCardToken() string {
	if wx.tokens.WXCardToken != nil {
		return wx.tokens.WXCardToken.GetToken()
	}
	return ""
}

func (wx *Weixin) GetHandler() (ret *WXHandler, err error) {
	if wx.handler == nil {
		ret, err = NewHandler(wx.cfg, wx.tokens.AccessToken, wx.log)
		if err != nil {
			wx.log.Error(err.Error())
		}
		wx.handler = ret
	}
	ret = wx.handler
	return
}

func (wx *Weixin) Start() error {
	h, err := wx.GetHandler()
	if err != nil {
		return err
	}
	server := &http.Server{Handler: h}
	return server.ListenAndServe()
}

func (wx *Weixin) StartTLS(cert, key string) error {
	h, err := wx.GetHandler()
	if err != nil {
		return err
	}
	server := &http.Server{Handler: h}
	return server.ListenAndServeTLS(cert, key)
}
