package runservice

import "testing"

func TestNewClient(t *testing.T) {
	NewClient("unix", "../test.sock")
	Send("reload", []byte("OK"))
	Send("ping", []byte("OK"))
	Close()
}
