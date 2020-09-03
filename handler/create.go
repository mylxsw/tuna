package handler

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/tuna/conf"
	"github.com/mylxsw/tuna/libs"
	"github.com/mylxsw/tuna/storage"
)

var r = rand.New(rand.NewSource(9999999))

type respForCreate struct {
	Link   string `json:"link"`
	Expire int64  `json:"expire"`
}

// Create 函数用于创建一个hash与url的对应关系
func Create(w http.ResponseWriter, r *http.Request) {
	// 默认有效时间 0，长期有效
	var expire int64 = 0
	var err error

	url := r.PostFormValue("url")
	expireStr := r.PostFormValue("expire")
	if expireStr != "" {
		expire, err = strconv.ParseInt(expireStr, 10, 64)
		if err != nil {
			libs.SendFormInvalidResponse(w, fmt.Sprintf("expire is invalid: %v", err))
			return
		}
	}

	if !libs.IsValidURL(url) {
		libs.SendFormInvalidResponse(w, "url is invalid")
		return
	}

	uniq := r.PostFormValue("uniq")

	if driver := storage.Default(); driver != nil {
		urlHash := createURLHash(url, uniq == "1" || uniq == "true")
		i := 6

		existedURL := driver.Get(urlHash[:i])
		for existedURL != "" && existedURL != url {
			log.Warningf("hash collision detected [%s] for %s", urlHash[:i], url)

			if i >= 32 {
				log.Warningf("oops, url [%s] has the same hash %s with someothers", url, urlHash[:i])
				goto ERR
			}
			i++

			existedURL = driver.Get(urlHash[:i])
		}

		if _, err := driver.Set(urlHash[:i], url, expire); err != nil {
			libs.SendInternalServerErrorResponse(w, fmt.Sprintf("driver set failed: %v", err))
			log.Errorf("driver set for %s=%s(%d) failed: %v", urlHash[:i], url, expire, err)
			return
		}

		log.Debugf("create new link %s for %s expired at %d", urlHash[:i], url, expire)

		_, _ = w.Write(libs.Success(respForCreate{
			Link:   fmt.Sprintf("%s/%s", conf.GetConf().PublicURL, urlHash[:i]),
			Expire: expire,
		}))

		return
	}

ERR:
	w.Write(libs.Failed("operation failed"))
}

// 生成URL哈希值
func createURLHash(url string, unique bool) string {
	salt := ""
	if unique {
		salt = randomSalt()
	}

	digest := md5.New()
	digest.Write([]byte(url + salt))
	urlHash := hex.EncodeToString(digest.Sum(nil))

	return urlHash
}

// 生成一个随机值
func randomSalt() string {
	return time.Now().Format("Mon Jan 2 15:04:05 MST 2006") + strconv.Itoa(r.Int())
}
