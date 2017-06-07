package runservice

import "testing"

func TestNewClient(t *testing.T) {
	NewClient("unix", "../test.sock")
	Send(2, "OK")
	Send(3, "OK")
	Send(4, "OK")
	Send(5, "OK")
	Send(6, "OK")
	Send(1, "")
	Close()
}
