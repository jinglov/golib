package utils

import "testing"

func TestByte2Int32(t *testing.T) {
	params := []struct {
		in  []byte
		out int32
	}{
		{
			in:  []byte{0, 0, 0, 10},
			out: 10,
		},
	}
	for _, param := range params {
		out := Byte2Int32(param.in)
		if out != param.out {
			t.Error("want:", param.out, "res:", out)
			t.FailNow()
		}
	}
}
