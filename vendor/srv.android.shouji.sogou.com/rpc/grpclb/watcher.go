package grpclb

import (
	"context"
	"fmt"

	etcd3 "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"google.golang.org/grpc/naming"
	"strings"
)

type watcher struct {
	re            *Resolver
	client        etcd3.Client
	isInitialized bool
}

func (w *watcher) Close() {
}

func (w *watcher) Next() ([]*naming.Update, error) {
	prefix := fmt.Sprintf("/%s/%s/", Prefix, w.re.serviceName)

	if !w.isInitialized {
		resp, err := w.client.Get(context.Background(), prefix, etcd3.WithPrefix())
		w.isInitialized = true
		if err == nil {
			addrs := extractAddrs(resp)
			if l := len(addrs); l != 0 {
				updates := make([]*naming.Update, l)
				for i := range addrs {
					updates[i] = &naming.Update{Op: naming.Add, Addr: addrs[i]}
				}
				return updates, nil
			}
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	rch := w.client.Watch(ctx, prefix, etcd3.WithPrefix())
	defer cancel()
	for wresp := range rch {
		for _, ev := range wresp.Events {
			//fmt.Println(ev)
			switch ev.Type {
			case mvccpb.PUT:
				return []*naming.Update{{Op: naming.Add, Addr: string(ev.Kv.Value)}}, nil
			case mvccpb.DELETE:
				return []*naming.Update{{Op: naming.Delete, Addr: string(ev.Kv.Value)}}, nil
			}
		}
	}
	return nil, nil
}

func extractAddrs(resp *etcd3.GetResponse) []string {
	addrs := []string{}

	if resp == nil || resp.Kvs == nil {
		return addrs
	}

	for i := range resp.Kvs {
		if v := string(resp.Kvs[i].Value); v != "" {
			if strings.Contains(v, ",") {
				for _, iv:= range strings.Split(v, ",") {
					addrs = append(addrs, iv)
				}
			}else{
				addrs = append(addrs, v)
			}
		}
	}

	return addrs
}
