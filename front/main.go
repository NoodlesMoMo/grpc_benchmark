package main

import (
	"github.com/qiangxue/fasthttp-routing"
	"fmt"
	"github.com/valyala/fasthttp"
	"grpc_benchmark/front/controller"
)

func main() {

	router := routing.New()

	router.Get("/data", controller.DummyFrontHandler)

	fmt.Println("HTTP front server will listen :8000")

	fasthttp.ListenAndServe(":8000", router.HandleRequest)
}
