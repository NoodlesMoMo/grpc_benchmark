package grpclb

import (
	"context"
	"fmt"
	"strings"
	etcd3 "github.com/coreos/etcd/clientv3"
	"net"
)

var Prefix = "service.sogou.ime"
var doneChan = make(chan struct{})

func genLocalIPV4Addr(port string) string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}

	var ret []string

	for _, addr := range addrs {
		if ip, ok := addr.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() == nil{
				continue
			}
			ret = append(ret, net.JoinHostPort(ip.IP.String(), port))
		}
	}

	return strings.Join(ret, ",")
}

func Register(name, port string, target string, ttl int) error {
	addr := genLocalIPV4Addr(port)
	serviceKey := fmt.Sprintf("/%s/%s/%s", Prefix, name, addr)

	var err error
	client, err := etcd3.New(etcd3.Config{
		Endpoints: strings.Split(target, ","),
	})
	if err != nil {
		return err
	}
	resp, err := client.Grant(context.TODO(), int64(ttl))
	if err != nil {
		return err
	}

	if _, err := client.Put(context.TODO(), serviceKey, addr, etcd3.WithLease(resp.ID)); err != nil {
		return err
	}

	if _, err := client.KeepAlive(context.TODO(), resp.ID); err != nil {
		return err
	}

	go func() {
		<-doneChan
		client.Delete(context.Background(), serviceKey)
		doneChan <- struct{}{}
	}()


	return nil
}

func UnRegister() {
	doneChan <- struct{}{}
	<-doneChan
}

