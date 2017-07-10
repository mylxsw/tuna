package libs

import "testing"

func TestIsValidURL(t *testing.T) {

	validURLs := []string{
		"http://aicode.cc",
		"https://aicode.cc",
		"http://aicode.cc/?id=123&name=测试",
		"https://aicode.cc?id=555&name=测试+测试2",
	}

	for _, url := range validURLs {
		if !IsValidURL(url) {
			t.Fail()
		}
	}

	invalidURLs := []string{
		"http//aicode.cc",
		"xxx.ccc",
		"http:/aicode.cc",
		"http://aicode.cc?id=xxxxxxxxxxxxxxxxxxxxxxxadfdfadfsdfasdfasdfsadfasdfsadfasdfasdfasdfasdfsadfsadfsadfasdfsadfsadfasdfasdfasfxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
	}

	for _, url := range invalidURLs {
		if IsValidURL(url) {
			t.Fail()
		}
	}
}

func TestIsValidHash(t *testing.T) {
	validHashs := []string{
		"xxafcadg",
		"x",
		"2c38b9e45cec1b324dde4e3d5b22c648",
	}

	for _, hash := range validHashs {
		if !IsValidHash(hash) {
			t.Fail()
		}
	}

	invalidHashs := []string{
		"xxxa3#",
		"",
		"2c38b9e45cec1b324dde4e3d5b22c648x",
	}

	for _, hash := range invalidHashs {
		if IsValidHash(hash) {
			t.Fail()
		}
	}
}
