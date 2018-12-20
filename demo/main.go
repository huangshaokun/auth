package main

import (
	"authhandle"
	"flag"
	"fmt"
	"log"
	"loghandle"
	"net/http"

	"git.algor.tech/algor/algorlib/version"

	auth "github.com/abbot/go-http-auth"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	port := flag.Int("port", 8888, "http port")

	flag.Parse()

	//创建实例
	authobject, err := authhandle.NewAuth("userpasswrod", "log", nil)
	if err != nil {
		log.Println(err)
		return
	}

	//使用实例的验证用户的函数
	authenticator := auth.NewBasicAuthenticator("example.com", authobject.Secret)

	//请求来时验证
	http.HandleFunc("/", authenticator.Wrap(func(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
		loghandle.SendJSON(nil, w, nil, version.GetVersion("demoserver", "1.0.0", "描述"), 0)
	}))

	err = http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Println(err)
		return
	}
}
