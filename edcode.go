package main

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"github.com/gorilla/securecookie"
)

type (
	SecretPairs [][]byte
	Obj         map[interface{}]interface{}
	JSON        map[string]interface{}

	Manager struct {
		SecretPairs
		SessionName string
		Serializer  securecookie.Serializer
	}
)

func (sp SecretPairs) GetSecCookies() []*securecookie.SecureCookie {
	secCookies := make([]*securecookie.SecureCookie, len(sp)/2+len(sp)%2)
	for i := 0; i < len(sp); i += 2 {
		var blockKey []byte
		if i+1 < len(sp) {
			blockKey = sp[i+1]
		}
		secCookies[i/2] = securecookie.New(sp[i], blockKey)
	}
	return secCookies
}

func EncodeMulti(name string, value interface{},
	cs ...*securecookie.SecureCookie) (string, error) {
	var err error
	var encoded string
	for _, c := range cs {
		encoded, err = c.Encode(name, value)
		if err == nil {
			return encoded, nil
		}
	}
	return "", err
}

func DecodeMulti(name string, value string, dst interface{},
	cs ...*securecookie.SecureCookie) error {
	var err error
	for _, c := range cs {
		err = c.Decode(name, value, dst)
		if err == nil {
			return nil
		}
	}
	return err
}

func (m Manager) Encode(obj interface{}) (string, error) {
	cs := m.GetSecCookies()
	for _, c := range cs {
		c.SetSerializer(serializer)
	}

	dst, err := EncodeMulti(m.SessionName, obj, cs...)
	return dst, err
}

func (m Manager) Decode(cookie string) (interface{}, error) {
	cs := m.GetSecCookies()
	for _, c := range cs {
		c.SetSerializer(serializer)
	}

	switch m.Serializer {
	case securecookie.GobEncoder{}:
		var dst Obj
		err := DecodeMulti(m.SessionName, cookie, &dst, cs...)
		return dst, err
	case securecookie.JSONEncoder{}:
		var dst JSON
		err := DecodeMulti(m.SessionName, cookie, &dst, cs...)
		return dst, err
	default:
		var dst []byte
		err := DecodeMulti(m.SessionName, cookie, &dst, cs...)
		return dst, err
	}
}

func b64decode(value []byte) ([]byte, error) {
	decoded := make([]byte, base64.URLEncoding.DecodedLen(len(value)))
	b, err := base64.URLEncoding.Decode(decoded, value)
	if err != nil {
		return nil, fmt.Errorf("base64 decode error")
	}
	return decoded[:b], nil
}

func DecodeWithoutKey(cookie string) (interface{}, error) {
	raw, err := b64decode([]byte(cookie))
	if err != nil {
		return nil, err
	}

	srzData, err := b64decode(bytes.SplitN(raw, []byte("|"), 3)[1])
	if err != nil {
		return nil, err
	}

	switch serializer {
	case securecookie.GobEncoder{}:
		var dst Obj
		err = serializer.Deserialize(srzData, &dst)
		if err != nil {
			return nil, err
		}
		return dst, nil
	case securecookie.JSONEncoder{}:
		var dst JSON
		err = serializer.Deserialize(srzData, &dst)
		if err != nil {
			return nil, err
		}
		return dst, nil
	default:
		var dst []byte
		err = serializer.Deserialize(srzData, &dst)
		if err != nil {
			return nil, err
		}
		return dst, nil
	}
}

func DisplayObj(obj interface{}) {
	fmt.Println(obj)
	fmt.Println("type detail: ")
	fmt.Println("{")

	switch obj := obj.(type) {
	case Obj:
		for k, v := range obj {
			fmt.Printf("    %v[%T]: %v[%T],\n", k, k, v, v)
		}
	case JSON:
		for k, v := range obj {
			fmt.Printf("    %v[%T]: %v[%T],\n", k, k, v, v)
		}
	}

	fmt.Println("}")
}
