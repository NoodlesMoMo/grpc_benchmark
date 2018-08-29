package rpc

import (
	"fmt"
	"grpc_benchmark/backend/rpc/proto"
	"net"

	"google.golang.org/grpc"
)

func ServeRPCForever() error {
	var (
		err      error
		listener net.Listener
	)

	srv := grpc.NewServer()

	proto.RegisterGetDummyDataServer(srv, new(DummyRPCServer))

	if listener, err = net.Listen("tcp4", ":9090"); err != nil {
		return err
	}

	fmt.Println("rpc will serve on :9090")
	return srv.Serve(listener)
}
