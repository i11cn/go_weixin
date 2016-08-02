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
		Log        *logger.Logger
		RestClient *rest.RestClient
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

func NewWeixin(cfg WXConfig) *Weixin {
	ret := &Weixin{WXGlobalInfo: WXGlobalInfo{cfg, g_log, rest.NewClient("api.weixin.qq.com", 0, "/cgi-bin")}}
	ret.WXGlobalInfo.RestClient.SSL = true
	ret.WXTokenMgr = DefaultTokenMgr(&ret.WXGlobalInfo)
	return ret
}

func (wx *Weixin) SetLogger(log *logger.Logger) *Weixin {
	wx.Log = log
	return wx
}

func (wx *Weixin) GetHandler() (ret *WXHandler, err error) {
	if wx.handler == nil {
		ret, err = NewHandler(wx.Config, wx.GetAccessToken(), wx.Log)
		if err != nil {
			wx.Log.Error(err.Error())
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
