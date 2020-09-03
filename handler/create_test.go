package handler

import (
	"testing"
)

func TestGenURLHash(t *testing.T) {
	if createURLHash("https://aicode.cc", false) != "18dbf40d1b8e2121c56d71b52a3ef169" {
		t.FailNow()
	}

	if createURLHash("https://aicode.cc", false) != createURLHash("https://aicode.cc", false) {
		t.FailNow()
	}

	if createURLHash("https://aicode.cc", true) == createURLHash("https://aicode.cc", true) {
		t.Errorf("hash should not equal")
	}
}
