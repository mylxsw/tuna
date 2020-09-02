package handler

import (
	"net/http"

	"github.com/mylxsw/tuna/libs"
	"github.com/mylxsw/tuna/storage"
)

type respForWelcome struct {
	URLCount int `json:"url_count"`
}

// Welcome 用于输出欢迎页面
func Welcome(w http.ResponseWriter, r *http.Request) {
	if driver := storage.Default(); driver != nil {
		w.Write(libs.Success(respForWelcome{
			URLCount: driver.Count(),
		}))

		return
	}

	w.Write(libs.Failed("operation failed"))
}
