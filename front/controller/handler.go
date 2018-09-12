package controller

import (
	"context"
	"flag"
	"grpc_benchmark/front/pool"
	"grpc_benchmark/front/proto"
	"sync"
	"time"

	"github.com/qiangxue/fasthttp-routing"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/resolver/etcd"
	"google.golang.org/grpc/grpclog"
	"os"
)

var (
	AccessType, HTTPBackend, RPCBackend string
)

var (
	clientCache *sync.Pool
	rpcConnPool *pool.Pool
)

var (
	singleRpcConn proto.GetDummyDataClient
)

func init() {

	flag.StringVar(&AccessType, "access", "grpc", "access backend type")
	flag.StringVar(&HTTPBackend, "http_backend", "http://127.0.0.1:8080", "http backend server address")
	//flag.StringVar(&RPCBackend, "rpc_backend", "127.0.0.1:9090", "rpc backend server address")
	flag.StringVar(&RPCBackend, "rpc_backend", "etcd:///10.135.54.224:2379", "rpc backend server address")

	flag.Parse()

	grpclog.SetLoggerV2(grpclog.NewLoggerV2(os.Stdout, os.Stdout, os.Stderr))

	singleRpcConn = func() proto.GetDummyDataClient {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		conn, err := grpc.DialContext(ctx, RPCBackend, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			panic(err)
		}
		return proto.NewGetDummyDataClient(conn)
	}()

	clientCache = &sync.Pool{
		New: func() interface{} {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			activeConn, err := rpcConnPool.Get(ctx)
			if err != nil {
				return nil
			}
			defer activeConn.Close()

			conn, err := activeConn.RpcConnection()
			if err != nil {
				return nil
			}
			return proto.NewGetDummyDataClient(conn)
		},
	}
}

func DummyFrontHandler(ctx *routing.Context) error {
	var (
		response []byte
		err      error
	)
	switch AccessType {
	case `grpc`:
		response, err = callWithGRPC(ctx)
	case `http`:
		response, err = callWithHTTP(ctx)
	}

	if err == nil {
		ctx.Write(response)
	} else {
		ctx.WriteString(err.Error())
	}

	return nil
}

func callWithHTTP(ctx *routing.Context) ([]byte, error) {

	return HttpGet200ok(HTTPBackend + "/inner/data")
}

func callWithGRPC(ctx *routing.Context) ([]byte, error) {

	req := &proto.DummyRequest{}

	rpcCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := singleRpcConn.Get(rpcCtx, req)

	return resp.GetPayload(), err
}
