package main

/*
import (
	"encoding/json"
	"testing"
)

func TestT(t *testing.T) {
	jstr := `{"a":9391394749691704,"b":"12345","c":0.1}`
	jm := make(map[string]*Value)
	json.Unmarshal([]byte(jstr), &jm)
	for k, v := range jm {
		t.Log(k)
		t.Log(*v)
	}
	t.Log(jm)
}

type Value string

func (p *Value) UnmarshalJSON(data []byte) error {
	if data[0] == '"' && data[len(data)-1] == '"' {
		*p = Value(data[1 : len(data)-1])
	} else {
		*p = Value(data)
	}
	return nil
}
*/
