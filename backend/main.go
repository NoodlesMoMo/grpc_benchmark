package main

import (
	"grpc_benchmark/backend/rpc"
	"sync"
)

var (
	waitGroup = sync.WaitGroup{}
)

func withWaitServer(server func() error) {
	if server == nil {
		return
	}

	waitGroup.Add(1)
	go func() {
		if e := server(); e != nil {
			panic(e)
		}
		waitGroup.Done()
	}()
}

func main() {

	//withWaitServer(http.ServeHTTPForever)

	rpc.ServeRPCForever()

	//waitGroup.Wait()
}
