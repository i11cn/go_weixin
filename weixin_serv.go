package weixin

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"strings"
)

type Weixin struct {
	WXConfig
	AccessToken string
	UserCustom  interface{}

	onValidateFail OnValidateFail
	unsupported    UnsupportedRequest
	onRequestError OnRequestError

	onTextRequest     OnTextRequest
	onLocationRequest OnLocationRequest
}

func (serv *Weixin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func (serv *Weixin) GetAccessToken() string {
	return serv.AccessToken
}
