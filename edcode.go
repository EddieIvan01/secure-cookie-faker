package main

import (
	"fmt"

	"github.com/gorilla/securecookie"
)

type (
	SecretPairs [][]byte
	Obj         map[interface{}]interface{}

	Manager struct {
		SecretPairs
		SessionName string
	}
)

func (sp SecretPairs) getCodecs() []securecookie.Codec {
	return securecookie.CodecsFromPairs(sp...)
}

func (m Manager) Encode(obj Obj) (string, error) {
	cs := m.getCodecs()
	dst, err := securecookie.EncodeMulti(m.SessionName, obj, cs...)
	return dst, err
}

func (m Manager) Decode(cookie string) (Obj, error) {
	cs := m.getCodecs()
	dst := make(Obj)
	err := securecookie.DecodeMulti(m.SessionName, cookie, &dst, cs...)
	return dst, err
}

func DisplayObj(obj Obj) {
	fmt.Println(obj)
	fmt.Println("type detail: ")
	fmt.Println("{")
	for k, v := range obj {
		fmt.Printf("    %v[%T]: %v[%T],\n", k, k, v, v)
	}
	fmt.Println("}")
}
