package utils

import "testing"

func TestMd5Sum(t *testing.T) {
	params := []struct {
		filename string
		md5      string
	}{
		{filename: "http.go", md5: ""},
	}
	for _, param := range params {
		md5 := FileMd5Sum(param.filename)
		if md5 != param.md5 {
			t.Error("want:", param.md5, "res:", md5)
			t.FailNow()
			return
		}
	}
}

func TestExecCommand(t *testing.T) {
	params := []struct {
		cmd    string
		hasErr bool
	}{
		/*		{
				cmd:    "ls",
				hasErr: false,
			},*/
		{
			cmd:    "lmn",
			hasErr: true,
		},
	}
	for _, param := range params {
		out, err := ExecCommand(param.cmd)
		if param.hasErr && err == nil || !param.hasErr && err != nil {
			t.Error("want:", param.hasErr, "res:", err)
			t.FailNow()
			return
		}
		if err != nil {
			t.Log(err)
		}
		t.Log(string(out))
	}
}
