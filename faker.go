package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"os"

	"github.com/gorilla/securecookie"
)

const version = "v0.1"

var (
	enc = flag.NewFlagSet("enc", flag.ExitOnError)
	dec = flag.NewFlagSet("dec", flag.ExitOnError)

	sessionName string

	// string like "key" or "key1, key2, key3"
	// there can be no space between the keys
	secretKey string

	cookie string

	// string like "{key1[type]: value1[type], key2[type]: value2[type]}"
	object string

	serializeWay string
	serializer   securecookie.Serializer
)

func init() {
	enc.StringVar(&object, "o", "",
		"object to be encoded, `string` like \"{key1[type]: value1[type], key2[type]: value2[type]}\"\n"+
			"type hint could be `int`, `uint`, `float`, `bool`, `string`, `byte`\n"+
			"when type is `string`, it could be omitted. like this {str1: str2}\n")
	enc.StringVar(&sessionName, "n", "", "cookie name")
	enc.StringVar(&secretKey, "k", "",
		"secret keys, string like \"key\" or \"key1, key2, key3\"")
	enc.StringVar(&serializeWay, "way", "gob",
		"serialize way: gob | json | nop, default is gob")

	dec.StringVar(&cookie, "c", "", "cookie to be decoded")
	dec.StringVar(&sessionName, "n", "", "the cookie name")
	dec.StringVar(&secretKey, "k", "",
		"secret keys, string like \"key\" or multiple keys like \"key1, key2, key3\"")
	dec.StringVar(&serializeWay, "way", "gob",
		"serialize way: gob | json | nop")

	enc.Usage = usage
	dec.Usage = usage

	gob.Register([]interface{}{})
}

func usage() {
	fmt.Printf("Secure Cookie Faker %v\n\n"+
		"Usage: faker [enc/dec] [-n cookie_name] [-k secret_key] [-o object_string / -c cookie_string]\n\n"+
		"Options:\n", version)

	if enc.Parsed() {
		enc.PrintDefaults()
	} else {
		dec.PrintDefaults()
	}
}

func checkParams() bool {
	if !enc.Parsed() && !dec.Parsed() {
		fmt.Println(os.Args[0] + " [enc/dec] --help to")
		return false
	}

	if enc.Parsed() && object == "" {
		fmt.Println("-o param is required")
		return false
	}

	if dec.Parsed() && cookie == "" {
		fmt.Println("-c param is required")
		return false
	}
	return true
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println(os.Args[0] + " [enc/dec] --help")
		return
	}

	switch os.Args[1] {
	case "enc":
		enc.Parse(os.Args[2:])
	case "dec":
		dec.Parse(os.Args[2:])
	default:
		fmt.Println(os.Args[0] + " [enc/dec] --help")
		return
	}

	ok := checkParams()
	if !ok {
		return
	}

	switch serializeWay {
	case "gob":
		serializer = securecookie.GobEncoder{}
	case "json":
		serializer = securecookie.JSONEncoder{}
	case "nop":
		serializer = securecookie.NopEncoder{}
	default:
		fmt.Println("unrecognized serialized way")
		return
	}

	sks, err := ParseSecretKeys(secretKey)
	if err != nil && enc.Parsed() {
		fmt.Println(err.Error())
		return
	}

	manager := Manager{
		sks,
		sessionName,
		serializer,
	}

	if enc.Parsed() {
		var o interface{}
		var err error

		if serializeWay == "gob" || serializeWay == "json" {
			o, err = ParseObjString(object, serializeWay)
			if err != nil {
				if err == ErrInvalidObjFormat && serializeWay == "json" {
					o = []byte(object)
				} else {
					fmt.Println(err.Error())
					return
				}
			}
		} else {
			o = []byte(object)
		}

		ret, err := manager.Encode(o)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(ret)

	} else {
		var ret interface{}
		var err error

		if secretKey == "" {
			ret, err = DecodeWithoutKey(cookie)
		} else {
			ret, err = manager.Decode(cookie)
		}
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if serializeWay == "gob" || serializeWay == "json" {
			DisplayObj(ret)
		} else {
			fmt.Println(string(ret.([]byte)))
		}
	}
}
