package etcd3

import (
	"time"

	"github.com/coredns/coredns/middleware"
	"github.com/coredns/coredns/middleware/etcd/msg"
	"github.com/coredns/coredns/middleware/proxy"
	"github.com/coredns/coredns/request"

	"github.com/coreos/etcd/clientv3"
)

// Etcd3 is a middleware talks to an etcd cluster.
type Etcd3 struct {
	Next        middleware.Handler
	Zones       []string
	PathPrefix  string
	Proxy       proxy.Proxy             // Proxy for looking up names during the resolution process
	Stubmap     *map[string]proxy.Proxy // list of proxies for stub resolving.
	Fallthrough bool

	kv        *clientv3.KV
	endpoints []string // Stored here as well, to aid in testing.
}

func (e *Etcd3) Records(state request.Request, exact bool) ([]msg.Service, error) {
	return nil, nil
}

/*
func (e *Etcd3) get(path string, recursive bool) (*etcdc.Response, error) {

	gr, _ := kv.Get(ctx, "key")
	fmt.Println("Value: ", string(gr.Kvs[0].Value), "Revision: ", gr.Header.Revision)

}
*/

const (
	defaultPriority = 10 // default priority when nothing is set
	defaultTTL      = 5  // default ttl when nothing is set
	defaultTimeout  = 5 * time.Second
)
