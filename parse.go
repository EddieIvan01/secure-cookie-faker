package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrEmptyParam       = errors.New("param is empty string")
	ErrInvalidObjFormat = errors.New("obj string format is invalid")
	ErrUnsupportedType  = "unsupported field type: %v"
)

type (
	KV  [2]string
	KVS []KV
)

// TODO: support slice type
func setVal(t string, v string) (interface{}, error) {
	var ret interface{}

	switch t {
	case "int":
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		ret = i

	case "uint":
		u, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return nil, err
		}
		ret = u

	case "float":
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, err
		}
		ret = f

	case "bool":
		b, err := strconv.ParseBool(v)
		if err != nil {
			return nil, err
		}
		ret = b

	case "string":
		ret = v

	case "byte":
		ret = v[0]

	default:
		return nil, fmt.Errorf(ErrUnsupportedType, t)
	}

	return ret, nil
}

// KVS is like [["k1[t1]", "v1[t2]"], ["k1[t3]", "v1[t4]"]]
func (s KVS) parseGob() (Obj, error) {
	var key, value string
	var keyT, valueT string
	var k, v interface{}
	var tmp []string
	var err error

	o := make(Obj)

	for _, kv := range s {
		tmp = strings.Split(kv[0], "[")
		key = tmp[0]
		if len(tmp) == 1 {
			keyT = "string"
		} else {
			keyT = tmp[1][:len(tmp[1])-1]
		}

		k, err = setVal(keyT, key)
		if err != nil {
			return nil, err
		}

		tmp = strings.Split(kv[1], "[")
		value = tmp[0]
		if len(tmp) == 1 {
			valueT = "string"
		} else {
			valueT = tmp[1][:len(tmp[1])-1]
		}

		v, err = setVal(valueT, value)
		if err != nil {
			return nil, err
		}

		o[k] = v
	}

	return o, nil
}

// parse JSON: map[string]interface{}
func (s KVS) parseJSON() (JSON, error) {
	var key, value string
	var keyT, valueT string
	var k, v interface{}
	var tmp []string
	var err error

	o := make(JSON)

	for _, kv := range s {
		tmp = strings.Split(kv[0], "[")
		key = tmp[0]
		keyT = "string"

		k, err = setVal(keyT, key)
		if err != nil {
			return nil, err
		}

		tmp = strings.Split(kv[1], "[")
		value = tmp[0]
		if len(tmp) == 1 {
			valueT = "string"
		} else {
			valueT = tmp[1][:len(tmp[1])-1]
		}

		v, err = setVal(valueT, value)
		if err != nil {
			return nil, err
		}

		o[k.(string)] = v
	}

	return o, nil
}

// "key" or "key1, key2, key3"
func ParseSecretKeys(ks string) ([][]byte, error) {
	if ks == "" {
		return nil, ErrEmptyParam
	}

	keys := strings.Split(ks, ",")
	ret := make([][]byte, 0)

	for _, k := range keys {
		k = strings.TrimSpace(k)
		ret = append(ret, []byte(k))
	}

	return ret, nil
}

// "{key1[type]: value1[type], key2[type]: value2[type]}"
func ParseObjString(s string, serialize string) (interface{}, error) {
	// trim white space
	if s == "" {
		return nil, ErrEmptyParam
	}
	s = strings.TrimSpace(s)

	// trim "{}"
	if s[0] != '{' || s[len(s)-1] != '}' {
		return nil, ErrInvalidObjFormat
	}
	s = s[1:]
	s = s[:len(s)-1]

	// split k-v
	// key[type]: value[type]
	kvStrs := strings.Split(s, ",")
	kvs := make(KVS, 0, 4)
	tmp := []string{}

	for _, kv := range kvStrs {
		tmp = strings.Split(kv, ":")
		tmp[0] = strings.TrimSpace(tmp[0])
		tmp[1] = strings.TrimSpace(tmp[1])
		kvs = append(kvs, KV{
			tmp[0], tmp[1],
		})
	}

	switch serialize {
	case "gob":
		return kvs.parseGob()
	default:
		return kvs.parseJSON()
	}

}
