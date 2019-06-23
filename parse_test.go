package main

import "testing"

func TestSetVal(t *testing.T) {
	var err error
	i, err := setVal("int", "5")
	f64, err := setVal("float", "1.1")
	b, err := setVal("bool", "true")
	s, err := setVal("string", "ss")
	// c, err := setVal("byte", "a")
	if err != nil {
		t.Error(err.Error())
	} else {
		if i != 5 || f64 != 1.1 || b != true || s != "ss" {
			t.Error("setVal err")
		} else {
			t.Log("setVal ok")
		}
	}
}

func TestParseSecretKeys(t *testing.T) {
	ks, err := ParseSecretKeys("key1,key2,key3")
	if err != nil {
		t.Error(err.Error())
	} else {
		t.Log(ks)
	}
}

func TestKVSparse(t *testing.T) {
	k := KVS{
		KV{"user[string]", "admin[string]"},
		KV{"id[string]", "0[int]"},
	}
	o, err := k.parse()
	if err != nil {
		t.Error(err.Error())
	} else {
		t.Log(o)
	}
}

func TestParseObjString(t *testing.T) {
	s := "{user[string]: admin[string], id[string]: 0[int]}"
	o, err := ParseObjString(s)
	if err != nil {
		t.Error(err.Error())
	} else {
		t.Log(o)
	}
}
