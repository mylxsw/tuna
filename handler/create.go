package handler

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/mylxsw/tuna/libs"
	"github.com/mylxsw/tuna/storage"
	log "github.com/sirupsen/logrus"
)

var r = rand.New(rand.NewSource(9999999))

type respForCreate struct {
	Link   string `json:"link"`
	Expire int64  `json:"expire"`
}

// Create 函数用于创建一个hash与url的对应关系
func Create(w http.ResponseWriter, r *http.Request) {

	url := r.PostFormValue("url")
	expire, _ := strconv.ParseInt(r.PostFormValue("expire"), 10, 64)

	link := ""
	if driver := storage.Default(); driver != nil {
		urlHash := genURLHash(url, true)

		i := 6
		link = urlHash[:i]
		for driver.Get(link) != "" {
			log.Warningf("hash collision detected [%s] for %s", link, url)

			if i >= 32 {
				log.Warningf("oops, url [%s] has the same hash %s with someothers", url, link)
				goto ERR
			}
			i++
			link = urlHash[:i]
		}

		driver.Set(link, url, expire)

		log.Debugf("create new link %s for %s expired at ", link, url, expire)

		w.Write(libs.Success(respForCreate{
			Link:   link,
			Expire: expire,
		}))

		return
	}

ERR:
	w.Write(libs.Failed("操作失败"))
}

// 生成URL哈希值
func genURLHash(url string, unique bool) string {
	salt := ""
	if unique {
		salt = genSalt()
	}

	digest := md5.New()
	digest.Write([]byte(url + salt))
	urlHash := hex.EncodeToString(digest.Sum(nil))

	return urlHash
}

// 生成一个随机值
func genSalt() string {
	return time.Now().Format("Mon Jan 2 15:04:05 MST 2006") + strconv.Itoa(r.Int())
}
