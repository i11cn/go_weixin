package weixin

import (
	"crypto/sha1"
	"fmt"
	rest "github.com/i11cn/go_rest_client"
	"net/http"
	"strings"
)

type (
	Tokens struct {
		AccessToken string
		JSApiToken  string
		WXCardToken string
	}
	Tokens2 struct {
		AccessToken2 *WXToken
		JSApiToken2  *WXToken
		WXCardToken2 *WXToken
	}

	Weixin struct {
		WXConfig
		Tokens2
		Tokens

		user_custom interface{}

		onValidateFail OnValidateFail
		unsupported    UnsupportedRequest
		onRequestError OnRequestError

		onTextRequest     OnTextRequest
		onImageRequest    OnImageRequest
		onVoiceRequest    OnVoiceRequest
		onVideoRequest    OnVideoRequest
		onLocationRequest OnLocationRequest
		onLinkRequest     OnLinkRequest
		onSubscribeEvent  OnSubscribeEvent
		onQRScanEvent     OnQRScanEvent
		onLocationEvent   OnLocationEvent
		onMenuEvent       OnMenuEvent
		onLinkEvent       OnLinkEvent
	}
)

func (serv *Weixin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("收到一次", r.Method, "请求 : ")
	fmt.Println(r.URL.Query())
	v := r.URL.Query()
	if len(v) >= 3 && exist_all_values(v, []string{"signature", "timestamp", "nonce"}) {
		strs := []string{serv.Token, v.Get("timestamp"), v.Get("nonce")}
		if strs[0] > strs[2] {
			strs[0], strs[2] = strs[2], strs[0]
		}
		if strs[0] > strs[1] {
			strs[0], strs[1] = strs[1], strs[0]
		} else if strs[1] > strs[2] {
			strs[1], strs[2] = strs[2], strs[1]
		}
		str := strings.Join(strs, "")
		sign := fmt.Sprintf("%x", sha1.Sum([]byte(str)))
		if v.Get("signature") != sign {
			serv.onValidateFail(w, r)
			return
		}
	}
	switch strings.ToUpper(r.Method) {
	case "GET":
		serv.onGet(w, r)

	case "POST":
		serv.onPost(w, r)

	default:
		w.WriteHeader(500)
	}
}

func (serv *Weixin) GetRestClient() *rest.RestClient {
	return nil
}

func (serv *Weixin) GetAccessToken() string {
	return serv.AccessToken
}
