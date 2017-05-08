package driver

import "testing"

func TestInitMysql(t *testing.T) {
	cfg := &MysqlConfig{
		DbHost:         "127.0.0.1:3306",
		DbUser:         "root",
		DbPwd:          "111",
		DbName:         "test",
		DbMaxIdleConns: 1,
		DbMaxOpenConns: 1,
	}
	err := RegisterMysql("test", cfg)
	if err != nil {
		t.Log(err)
	}
}
