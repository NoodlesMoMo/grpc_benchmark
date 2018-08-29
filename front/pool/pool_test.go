package pool

import (
	"fmt"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"srv.android.shouji.sogou.com/rpc/grpclb"
)

func TestPool(t *testing.T) {

	dialWithLoadBalancer := func() (*grpc.ClientConn, error) {
		r := grpclb.NewResolver("configure")
		b := grpc.RoundRobin(r)
		c, e := grpc.Dial("http://127.0.0.1:2379", grpc.WithInsecure(), grpc.WithBalancer(b), grpc.WithBlock())
		return c, e
	}

	connDetect := func(conn *grpc.ClientConn, i time.Time) error {

		state := conn.GetState()
		switch state {
		case connectivity.Shutdown:
			t.Error(ErrConnShutdown)
		case connectivity.TransientFailure:
			t.Error(ErrConnFailture)
		default:
			fmt.Println(">>> debug grpc connection: ", state)
		}

		return nil
	}

	p, err := New(dialWithLoadBalancer, WithPoolSize(48), WithPoolCapacity(48),
		WithPoolTimeout(3*time.Second), WithTestOnBorrow(connDetect))

	fmt.Println(p, err)
}
