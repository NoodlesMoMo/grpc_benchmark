package grpclb

import (
	"errors"
	"strings"

	etcd3 "github.com/coreos/etcd/clientv3"
	"google.golang.org/grpc/naming"
)

type Resolver struct {
	serviceName string
}

func NewResolver(serviceName string) *Resolver {
	return &Resolver{serviceName: serviceName}
}

func (re *Resolver) Resolve(target string) (naming.Watcher, error) {
	if re.serviceName == "" {
		return nil, errors.New("service name is empty")
	}

	client, err := etcd3.New(etcd3.Config{
		Endpoints: strings.Split(target, ","),
	})

	if err != nil {
		return nil, err
	}

	return &watcher{re: re, client: *client}, nil
}
