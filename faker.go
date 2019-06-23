package main

import (
	"flag"
	"fmt"
)

const version = "v0.1"

var (
	enc bool
	dec bool

	sessionName string

	// string like "key" or "key1, key2, key3"
	// there can be no space between the keys
	secretKey string

	cookie string

	// string like "{key1[type]: value1[type], key2[type]: value2[type]}"
	object string

	h bool
)

func init() {
	flag.BoolVar(&enc, "enc", false, "encode mode, object => cookie")
	flag.BoolVar(&dec, "dec", false, "decode mode, cookie => object")

	flag.StringVar(&object, "o", "",
		"object to be encoded, `string` like \"{key1[type]: value1[type], key2[type]: value2[type]}\"\n"+
			"type could be `int`, `float`, `bool`, `string`, `byte`\n"+
			"when type is `string`, it could be omitted. like this {str1: str2}\n"+
			"if mode is encode, this param is required")
	flag.StringVar(&cookie, "c", "", "cookie to be decoded\n"+
		"if mode is decode, this param is required")

	flag.StringVar(&sessionName, "n", "", "cookie name")
	flag.StringVar(&secretKey, "k", "",
		"secret keys, string like \"key\" or \"key1, key2, key3\"")

	flag.BoolVar(&h, "h", false, "show help")

	flag.Usage = usage
}

func usage() {
	fmt.Printf("Secure Cookie Faker %v\n\n"+
		"Usage: faker [-enc/dec] [-n cookie_name] [-k secret_key] [-o object_string / -c cookie_string]\n\n"+
		"Options:\n", version)
	flag.PrintDefaults()
}

func checkParams() bool {
	if (!enc && !dec) || (enc && dec) {
		fmt.Println("must only choose one mode, -enc/dec")
		return false
	}

	if enc && object == "" {
		fmt.Println("encode mode must have -o param")
		return false
	}

	if dec && cookie == "" {
		fmt.Println("decode mode must have -c param")
		return false
	}

	if sessionName == "" || secretKey == "" {
		fmt.Println("cookie name and secret key is requeired")
		return false
	}

	return true
}

func main() {
	flag.Parse()
	if h {
		flag.Usage()
		return
	}
	ok := checkParams()
	if !ok {
		return
	}

	sks, err := ParseSecretKeys(secretKey)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	manager := Manager{
		sks,
		sessionName,
	}

	if enc {
		o, err := ParseObjString(object)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		ret, err := manager.Encode(o)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(ret)
	} else {
		ret, err := manager.Decode(cookie)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		DisplayObj(ret)
	}
}
