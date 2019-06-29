package main

import (
	"flag"
	"fmt"
	"os"
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
)

func init() {
	enc.StringVar(&object, "o", "",
		"object to be encoded, `string` like \"{key1[type]: value1[type], key2[type]: value2[type]}\"\n"+
			"type could be `int`, `uint`, `float`, `bool`, `string`, `byte`\n"+
			"when type is `string`, it could be omitted. like this {str1: str2}\n")
	enc.StringVar(&sessionName, "n", "", "cookie name")
	enc.StringVar(&secretKey, "k", "",
		"secret keys, string like \"key\" or \"key1, key2, key3\"")

	dec.StringVar(&cookie, "c", "", "cookie to be decoded")
	dec.StringVar(&sessionName, "n", "", "cookie name")
	dec.StringVar(&secretKey, "k", "",
		"secret keys, string like \"key\" or \"key1, key2, key3\"")

	enc.Usage = usage
	dec.Usage = usage
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

	if sessionName == "" || secretKey == "" {
		fmt.Println("cookie name and secret key are requeired")
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

	sks, err := ParseSecretKeys(secretKey)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	manager := Manager{
		sks,
		sessionName,
	}

	if enc.Parsed() {
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
