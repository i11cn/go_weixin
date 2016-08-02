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
		WXTokenMgr
		cfg     WXConfig
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
	g_log.AddAppender(logger.NewConsoleAppender("%T [%N] %L (%f) : %M"))
	//g_log.AddAppender(logger.NewSplittedFileAppender("%T [%N] %L (%f) : %M", "weixin.log", 24*time.Hour))
}

func NewWeixin(cfg WXConfig) *Weixin {
	ret := &Weixin{cfg: cfg, rc: rest.NewClient("api.weixin.qq.com", 0, "/cgi-bin"), log: g_log}
	ret.rc.SSL = true
	ret.WXTokenMgr = &default_token_mgr{rc: ret.rc, log: ret.log}
	return ret
}

func (wx *Weixin) SetLogger(log *logger.Logger) *Weixin {
	wx.log = log
	return wx
}

func (wx *Weixin) GetHandler() (ret *WXHandler, err error) {
	if wx.handler == nil {
		ret, err = NewHandler(wx.cfg, wx.GetAccessToken(), wx.log)
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
