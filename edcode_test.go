package main

import "testing"

var m = Manager{
	[][]byte{
		[]byte("secret"),
	},
	"session",
}

func TestEncode(t *testing.T) {
	initDec := map[interface{}]interface{}{
		"user": "admin",
		"id":   0,
	}

	enc, err := m.Encode(initDec)
	if err != nil {
		t.Errorf("%s\n", err.Error())
	} else {
		t.Logf("encode ok: %s\n", enc)
	}
}

func TestDecode(t *testing.T) {
	initDec := map[interface{}]interface{}{
		"user": "admin",
		"id":   0,
	}
	enc, err := m.Encode(initDec)

	dec, err := m.Decode(enc)
	if err != nil {
		t.Errorf("%v\n", err.Error())
	} else {
		t.Logf("decode ok: %v\n", dec)
	}
}
