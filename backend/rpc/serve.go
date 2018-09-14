package rpc

import (
	"flag"
	"fmt"
	"grpc_benchmark/backend/rpc/proto"
	"net"

	"google.golang.org/grpc"
)

var (
	flagPort string
)

func init() {
	flag.StringVar(&flagPort, "rpc_port", "9090", "rpc listen port")

	flag.Parse()
}

func ServeRPCForever() error {
	var (
		err      error
		listener net.Listener
	)

	srv := grpc.NewServer()

	proto.RegisterGetDummyDataServer(srv, new(DummyRPCServer))

	if listener, err = net.Listen("tcp4", ":"+flagPort); err != nil {
		return err
	}

	fmt.Println("rpc will serve on :9090")
	return srv.Serve(listener)
}
