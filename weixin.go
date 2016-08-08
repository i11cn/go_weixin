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

	WXGlobalInfo struct {
		Config     WXConfig
		Logger     *logger.Logger
		RestClient *rest.RestClient
	}

	WXComponent struct {
		*WXGlobalInfo
		*Weixin
	}

	WXMessage struct {
	}

	Weixin struct {
		WXGlobalInfo
		WXTokenMgr
		handler *WXHandler
		msg     *WXMessage
	}
)

var (
	g_log *logger.Logger = logger.GetLogger("weixin")
)

func init() {
	rand.Seed(time.Now().UnixNano())
	g_log.AddAppender(logger.NewConsoleAppender("%T [%N] %L (%f) : %M"))
	//g_log.AddAppender(logger.NewSplittedFileAppender("%T [%N] %L (%f) : %M", "weixin.log", 24*time.Hour))
}

func (wc WXComponent) SetGlobalInfo(info *WXGlobalInfo) {
	wc.WXGlobalInfo = info
}

func (wc WXComponent) GetToken() string {
	return wc.GetAccessToken().GetToken()
}

func NewWeixin(cfg WXConfig) *Weixin {
	ret := &Weixin{WXGlobalInfo: WXGlobalInfo{cfg, g_log, rest.NewClient("api.weixin.qq.com", 0, "/cgi-bin")}}
	ret.WXGlobalInfo.RestClient.SSL = true
	ret.WXTokenMgr = DefaultTokenMgr(&ret.WXGlobalInfo)
	return ret
}

func (wx *Weixin) SetLogger(log *logger.Logger) *Weixin {
	wx.Logger = log
	return wx
}

func (wx *Weixin) GetHandler() *WXHandler {
	if wx.handler == nil {
		wx.handler = NewHandler(&wx.WXGlobalInfo, wx)
	}
	return wx.handler
}

func (wx *Weixin) Start() error {
	server := &http.Server{Handler: wx.GetHandler()}
	return server.ListenAndServe()
}

func (wx *Weixin) StartTLS(cert, key string) error {
	server := &http.Server{Handler: wx.GetHandler()}
	return server.ListenAndServeTLS(cert, key)
}

func (wx *Weixin) StartAt(addr string) error {
	server := &http.Server{Addr: addr, Handler: wx.GetHandler()}
	return server.ListenAndServe()
}

func (wx *Weixin) StartTLSAt(addr string, cert, key string) error {
	server := &http.Server{Addr: addr, Handler: wx.GetHandler()}
	return server.ListenAndServeTLS(cert, key)
}
