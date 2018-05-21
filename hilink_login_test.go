package hilink

import (
	"testing"
)

func Test_CreatePassword(t *testing.T) {
	x := CreateLoginPassword("SampleToken", "admin", "admin")
	if x != "ZDkyZTE5YTAwZGM5YWUzNjZiMDQ3MGU1NmNmYzFhMTUzMTVhNTBkMzdhZWM0YWMwNzM5YjcxODE4M2I0YmMxNQ==" {
		t.Errorf("Bad result: %s", x)
	}
}
