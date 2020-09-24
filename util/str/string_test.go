package str

import "testing"

func TestRandomStr(t *testing.T) {
	str := RandomStr(16)
	t.Log(str)
}
