package http

import (
	"fmt"

	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

func ServeHTTPForever() error {

	router := routing.New()

	router.Get("/dummy/data", DummyHandler)

	fmt.Println("HTTP server will listen :8080")
	return fasthttp.ListenAndServe(":8080", router.HandleRequest)
}
