package handler

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	"github.com/mylxsw/tuna/libs"
	"github.com/mylxsw/tuna/storage"
)

type respForCreate struct {
	Link   string `json:"link"`
	Expire int    `json:"expire"`
}

// Create 函数用于创建一个hash与url的对应关系
func Create(w http.ResponseWriter, r *http.Request) {

	url := r.PostFormValue("url")
	expire, _ := strconv.Atoi(r.PostFormValue("expire"))

	digest := md5.New()
	digest.Write([]byte(url + time.Now().Format("Mon Jan 2 15:04:05 MST 2006")))
	urlHash := hex.EncodeToString(digest.Sum(nil))

	if driver := storage.Default(); driver != nil {
		i := 8
		link := urlHash[:i]
		for driver.Get(link) != "" {
			if i >= 32 {
				goto ERR
			}
			i++
			link = urlHash[:i]
		}

		driver.Set(link, url, expire)

		w.Write(libs.Success(respForCreate{
			Link:   link,
			Expire: expire,
		}))

		return
	}

ERR:
	w.Write(libs.Failed("操作失败"))
}
