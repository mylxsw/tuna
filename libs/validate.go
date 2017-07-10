package libs

import (
	"regexp"
)

var urlRegex = regexp.MustCompile(`^http(s)?://`)
var hashRegex = regexp.MustCompile(`^[a-zA-Z0-9-]{1,32}$`)

// IsValidURL 判断url地址是否合法
func IsValidURL(url string) bool {
	if len(url) <= 8 || len(url) >= 255 {
		return false
	}

	return urlRegex.MatchString(url)
}

// IsValidHash 判断hash值是否合法
func IsValidHash(hash string) bool {
	if len(hash) < 1 || len(hash) > 32 {
		return false
	}

	return hashRegex.MatchString(hash)
}
