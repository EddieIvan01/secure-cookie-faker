# Secure-Cookie-Faker

security tool to encode/decode Golang web-frameworks' client-side session cookie which use [gorilla/securecookie](https://github.com/gorilla/securecookie) or [gorilla/sessions](https://github.com/gorilla/sessions), such as Gin, Echo or Iris

***

## Stats

top stars Go web-frameworks' using of [gorilla/securecookie](https://github.com/gorilla/securecookie) or [gorilla/sessions](https://github.com/gorilla/sessions)

(stars count comes from [go-web-framework-stars](https://github.com/mingrammer/go-web-framework-stars))

| web framework                                    | stars | uses gorilla's securecookie or sessions? |
| ------------------------------------------------ | ----- | ---------------------------------------- |
| [gin](https://github.com/gin-gonic/gin)          | 28126 | ✔                                        |
| [beego](https://github.com/astaxie/beego)        | 20805 | ✘                                        |
| [iris](https://github.com/kataras/iris)          | 15102 | ✔                                        |
| [echo](https://github.com/labstack/echo)         | 14180 | ✔                                        |
| [kit](https://github.com/go-kit/kit)             | 13920 | ✘                                        |
| [revel](https://github.com/revel/revel)          | 11125 | ✘                                        |
| [martini](https://github.com/go-martini/martini) | 10572 | ✔                                        |

and many personal application using them as a basic web application toolkit [link](https://github.com/search?q=import+%22github.com%2Fgorilla%2Fsessions%22&type=Code)

***

## Usage

```
Secure Cookie Faker v0.1

Usage: faker [enc/dec] [-n cookie_name] [-k secret_key] [-o object_string / -c cookie_string]

Mode: 
  dec
        decode mode, cookie => object
  enc
        encode mode, object => cookie

Options:
  --help    show help
  -k string
        secret keys, string like "key" or multiple keys like "key1, key2, key3"
  -n string
        the cookie name
  -o string
        object to be encoded, string like "{key1[type]: value1[type], key2[type]: value2[type]}"
        type hint could be `int`, `float`, `bool`, `string`, `byte`
        when type is `string`, it could be omitted. like this {str1: str2}
        if mode is encode, this param is required
  -c string
        cookie to be decoded
        if mode is decode, this param is required
  -way string                                                                                   serialize way: gob | json | nop, default is gob (default "gob")
```

## Example

choosing a mode is required: enc or dec

```
./faker dec -c "MTU2MTE4NjQzNHxFXy1CQkFFQkEwOWlhZ0hfZ2dBQkVBRVFBQUJUXzRJQUF3WnpkSEpwYm1jTUJnQUVkWE5sY2daemRISnBibWNNQndBRllXUnRhVzRHYzNSeWFXNW5EQVFBQW1sa0EybHVkQVFDQUFBR2MzUnlhVzVuREFjQUJYQnZhVzUwQTJsdWRBUUZBUDBERFQ0PXwKR14WwPjXeUBZlZ0sKcEfRu-n7_va9drjsFaIEVahmA=="
```

`-c` : the cookie to be decoded

***

encode object

```
./faker enc -n "mysession" -k "secret" -o "{user: admin, id: 0[int]}"
```

`-o `: the object string，its like a K-V map, it should have type hints

`-n` : the cookie name, its required because the HMAC hash's generation relies on it

`-k` : the secret key(s), could be multiple: `-k "key1, key2"`, the first is hash key, the second is encrypt block key

when the element is `string` type, the type tag can be omitted

type tag can only be `int`,  `uint`,  `float`,  `bool`,  `string`,  `byte`

***

select a serializer

```
./faker enc -n "mysession" -k "secret" -o "some-string" -w json
./faker enc -n "mysession" -k "secret" -o "{id: 0[int]}" -w json
./faker enc -n "mysession" -k "secret" -o "some-string" -w nop
```
