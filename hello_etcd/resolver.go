package hello_etcd

import (
	"errors"
	"fmt"
	etcd3 "github.com/coreos/etcd/clientv3"
	"google.golang.org/grpc/naming"
	"strings"
)

type resolver struct {
	serviceName string
}

// NewResolver命名解析实现
func NewResolver(serviceName string) *resolver {

	return &resolver{serviceName}
}

func (re *resolver) Resolve(target string) (naming.Watcher, error) {

	if re.serviceName == "" {
		return nil, errors.New("grpclb:no service name provided")
	}
	client, err := etcd3.New(etcd3.Config{
		Endpoints: strings.Split(target, ","),
	})
	if err != nil {
		return nil, fmt.Errorf("grpclb:create etcd3 client failed:%s", err)
	}
	return &watcher{re: re, client: *client}, nil
}
