package etcd3

import (
	"context"
	"fmt"
	"time"

	"github.com/coredns/coredns/middleware"
	"github.com/coredns/coredns/middleware/etcd/msg"
	"github.com/coredns/coredns/middleware/pkg/cache"
	"github.com/coredns/coredns/request"

	etcdc "github.com/coreos/etcd/client"
	"github.com/coreos/etcd/clientv3"
	"github.com/mholt/caddy/caddyhttp/proxy"
)

// Etcd3 is a middleware talks to an etcd cluster.
type Etcd3 struct {
	Next        middleware.Handler
	Zones       []string
	Client      *clientv3.KV
	PathPrefix  string
	Proxy       proxy.Proxy             // Proxy for looking up names during the resolution process
	Stubmap     *map[string]proxy.Proxy // list of proxies for stub resolving.
	Fallthrough bool

	endpoints []string // Stored here as well, to aid in testing.
}

func (e *Etcd3) Records(state request.Request, exact bool) ([]msg.Service, error) {
	return nil, nil
}

/*
cli, err := clientv3.New(clientv3.Config{
	Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
	DialTimeout: 5 * time.Second,
})
if err != nil {
	// handle error!
}
*/
/*
 ctx, _ := context.WithTimeout(context.Background(), requestTimeout)
    cli, _ := clientv3.New(clientv3.Config{
        DialTimeout: dialTimeout,
        Endpoints: []string{"127.0.0.1:2379"},
    })
    defer cli.Close()
    kv := clientv3.NewKV(cli)
*/

func (e *Etcd3) get(path string, recursive bool) (*etcdc.Response, error) {

	hash := cache.Hash([]byte(path))

	resp, err := e.Inflight.Do(hash, func() (interface{}, error) {
		ctx, cancel := context.WithTimeout(e.Ctx, etcdTimeout)
		defer cancel()
		r, e := e.Client.Get(ctx, path, &etcdc.GetOptions{Sort: false, Recursive: recursive})
		if e != nil {
			return nil, e
		}
		return r, e
	})
	if err != nil {
		return nil, err
	}
	gr, _ := kv.Get(ctx, "key")
	fmt.Println("Value: ", string(gr.Kvs[0].Value), "Revision: ", gr.Header.Revision)

	// Modify the value of an existing key (create new revision)
	return resp.(*etcdc.Response), err
}

const (
	defaultPriority = 10 // default priority when nothing is set
	defaultTTL      = 5  // default ttl when nothing is set
	defaultTimeout  = 5 * time.Second
)
