package etcd

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"strings"

	"github.com/coreos/etcd/clientv3"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/internal/backoff"
	"google.golang.org/grpc/resolver"
)

const (
	defaultPort = "2379"
	defaultKey  = `hello_mgj`
)

var (
	defaultMinFrequency = 120 * time.Second
)

func init() {
	resolver.Register(NewETCDBuilder())
}

type etcdBuilder struct {
}

func NewETCDBuilder() resolver.Builder {
	return &etcdBuilder{}
}

func (b *etcdBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	etcdProxy, port, err := parseTarget(target.Endpoint)
	if err != nil {
		return nil, err
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdProxy + ":" + port},
		DialTimeout: 3 * time.Second,
	})
	if err != nil {
		return nil, errors.New("connect to etcd proxy error")
	}

	ctx, cancel := context.WithCancel(context.Background())
	rlv := &etcdResolver{
		cc:           cc,
		cli:          cli,
		ctx:          ctx,
		cancel:       cancel,
		freq:         5 * time.Minute,
		backoff:      backoff.Exponential{MaxDelay: defaultMinFrequency},
		proxyAddress: etcdProxy + ":" + port,
		t:            time.NewTimer(0),
		rn:           make(chan struct{}, 1),
		im:           make(chan []resolver.Address),
		wg:           sync.WaitGroup{},
	}

	rlv.wg.Add(2)
	go rlv.watcher()
	go rlv.FetchBackendsWithWatch()

	return rlv, nil
}

func (b *etcdBuilder) Scheme() string {
	return "etcd"
}

type etcdResolver struct {
	retry        int
	proxyAddress string
	freq         time.Duration
	backoff      backoff.Exponential
	ctx          context.Context
	cancel       context.CancelFunc
	cc           resolver.ClientConn
	cli          *clientv3.Client
	t            *time.Timer

	rn chan struct{}
	im chan []resolver.Address

	wg sync.WaitGroup
}

func (r *etcdResolver) ResolveNow(opt resolver.ResolveNowOption) {
	select {
	case r.rn <- struct{}{}:
	default:
	}
}

func (r *etcdResolver) Close() {
	r.cancel()
	r.wg.Wait()
	r.t.Stop()
}

func (r *etcdResolver) watcher() {
	defer r.wg.Done()

	for {
		select {
		case <-r.ctx.Done():
			return
		case addrs := <-r.im:
			if len(addrs) > 0 {
				r.retry = 0
				r.t.Reset(r.freq)
				r.cc.NewAddress(addrs)
				continue
			}
		case <-r.t.C:
		case <-r.rn:
		}

		result := r.FetchBackends()

		if len(result) == 0 {
			r.retry++
			r.t.Reset(r.backoff.Backoff(r.retry))
		} else {
			r.retry = 0
			r.t.Reset(r.freq)
		}

		r.cc.NewAddress(result)
	}
}

func (r *etcdResolver) FetchBackendsWithWatch() {
	defer r.wg.Done()

	for {
		select {
		case <-r.ctx.Done():
			return
		case data := <-r.cli.Watch(r.ctx, defaultKey):
			result := make([]resolver.Address, 0)
			for _, ev := range data.Events {
				//if ev.Type == mvccpb.PUT || ev.Type == mvccpb.DELETE {
					for _, addr := range strings.Split(string(ev.Kv.Value), ",") {
						addr := strings.TrimSpace(addr)
						if addr == "" {
							continue
						}
						result = append(result, resolver.Address{Addr: addr})
					}
					grpclog.Infoln("data changed:", result)
					r.im <- result
				//}
			}
		}
	}
}

func (r *etcdResolver) FetchBackends() []resolver.Address {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result := make([]resolver.Address, 0)

	resp, err := r.cli.Get(ctx, defaultKey)
	if err != nil {
		grpclog.Errorln("Fetch etcd proxy error:", err)
		return result
	}

	for _, kv := range resp.Kvs {
		for _, addr := range strings.Split(string(kv.Value), ",") {
			addr := strings.TrimSpace(addr)
			if addr == "" {
				continue
			}
			result = append(result, resolver.Address{Addr: addr})
		}
	}

	return result
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
