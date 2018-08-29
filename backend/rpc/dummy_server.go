package rpc

import (
	"context"
	"grpc_benchmark/backend/rpc/proto"
	"grpc_benchmark/backend/common"
)

type DummyRPCServer struct {

}

func (d *DummyRPCServer) Get(ctx context.Context, req *proto.DummyRequest) (*proto.DummyResponse, error) {
	resp := &proto.DummyResponse{
		Payload: common.GenerateDummyData(),
	}
	return resp, nil
}
