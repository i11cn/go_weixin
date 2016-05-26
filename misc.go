package weixin

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/i11cn/go_logger"
	"net/url"
	"strings"
)

func exist_all_values(values url.Values, keys ...string) error {
	for _, key := range keys {
		if _, ok := values[key]; !ok {
			return errors.New(fmt.Sprintf("%s 不存在", key))
		}
	}
	return nil
}

func check_sign(token, sign, ts, nonce string, log *logger.Logger) bool {
	log.Trace("准备教研签名: ", "signature=", sign, " token=", token, " timestamp=", ts, " nonce=", nonce)
	strs := []string{token, ts, nonce}
	if strs[0] > strs[2] {
		strs[0], strs[2] = strs[2], strs[0]
	}
	if strs[0] > strs[1] {
		strs[0], strs[1] = strs[1], strs[0]
	} else if strs[1] > strs[2] {
		strs[1], strs[2] = strs[2], strs[1]
	}
	str := strings.Join(strs, "")
	s := fmt.Sprintf("%x", sha1.Sum([]byte(str)))
	log.Trace("计算出来的签名是: ", s)
	return sign == s
}
