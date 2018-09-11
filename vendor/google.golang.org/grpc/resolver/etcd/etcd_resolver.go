package etcd

import (
	"errors"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc/resolver"
)

const (
	defaultPort = "2379"
)

func init() {
	resolver.Register(NewETCDBuilder())
}

type etcdBuilder struct {
	ttl time.Duration
}

func NewETCDBuilder() resolver.Builder {
	return &etcdBuilder{ttl: 30 * time.Second}
}

func (b *etcdBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	etcdProxy, port, err := parseTarget(target.Endpoint)
	if err != nil {
		return nil, err
	}

	rlv := &etcdResolver{
		cc:           cc,
		proxyAddress: etcdProxy + ":" + port,
		rn:           make(chan struct{}, 1),
		q:            make(chan struct{}),
	}

	go rlv.watcher()

	return rlv, nil
}

func (b *etcdBuilder) Scheme() string {
	return "etcd"
}

type etcdResolver struct {
	proxyAddress string
	cc           resolver.ClientConn
	ip           []resolver.Address

	rn chan struct{}
	q  chan struct{}
}

func (r *etcdResolver) ResolveNow(opt resolver.ResolveNowOption) {
	select {
	case r.rn <- struct{}{}:
	default:
	}
}

func (r *etcdResolver) Close() {
	close(r.q)
}

func (r *etcdResolver) watcher() {
	for {
		select {
		case <-r.q:
			return
		case <-r.rn:
		}

		r.cc.NewAddress(r.ip)
	}
}

func parseTarget(target string) (host, port string, err error) {
	if target == "" {
		return "", "", errors.New("invalid target")
	}

	if ip := net.ParseIP(target); ip != nil {
		return target, defaultPort, nil
	}
	if host, port, err = net.SplitHostPort(target); err == nil {
		if port == "" {
			return "", "", errors.New("Invalid address format")
		}
		if host == "" {
			host = "localhost"
		}
		return host, port, nil
	}
	if host, port, err = net.SplitHostPort(target + ":" + defaultPort); err == nil {
		return host, port, nil
	}
	return "", "", fmt.Errorf("invalid target address %v, error info: %v", target, err)
}
