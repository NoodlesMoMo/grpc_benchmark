package http

import (
	"grpc_benchmark/backend/common"

	"github.com/qiangxue/fasthttp-routing"
)

func DummyHandler(ctx *routing.Context) error {

	ctx.Write(common.GenerateDummyData())

	return nil
}
