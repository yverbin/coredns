package etcd3

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/coredns/coredns/middleware/etcd/msg"
	"github.com/coredns/coredns/middleware/pkg/tls"
	"github.com/coredns/coredns/middleware/proxy"

	"github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
)

func (e *Etcd3) set(t *testing.T, k string, m msg.Service) {
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	path, _ := msg.PathWithWildcard(k, e.PathPrefix)

	_, err = e.kv.Put(context.TODO(), path, string(b))
	if err != nil {
		log.Fatal(err)
	}
}

func (e *Etcd3) delete(t *testing.T, k string) {
	path, _ := msg.PathWithWildcard(k, e.PathPrefix)

	e.kv.Delete(context.TODO(), path, clientv3.WithPrefix())
}

func testEtcd3() *Etcd3 {
	endpoints := []string{"http://localhost:2379"}
	tc, _ := tls.NewTLSConfigFromArgs()
	client, _ := newEtcd3KV(endpoints, tc)

	return &Etcd3{
		Proxy:      proxy.NewLookup([]string{"8.8.8.8:53"}),
		PathPrefix: "skydns",
		Zones:      []string{"skydns.test.", "skydns_extra.test.", "in-addr.arpa."},
		kv:         client,
	}
}

func TestEtcd3(t *testing.T) {
	svcs := []msg.Service{
		{Host: "10.0.0.1", Port: 8080, Key: "a.server1.prod.region1.skydns.test."},
		{Host: "10.0.0.2", Port: 8080, Key: "b.server1.prod.region1.skydns.test."},
	}
	e := testEtcd3()
	for _, s := range svcs {
		e.set(t, s.Key, s)
		defer e.delete(t, s.Key)
	}
}
