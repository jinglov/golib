package cmdservice

import "testing"

func TestNewService(t *testing.T) {
	NewService("unix", "../test.sock")
	Start()
	select {}
}
