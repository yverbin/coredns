package etcd3

import (
	"github.com/miekg/dns"
	"golang.org/x/net/context"
)

func (e *Etcd3) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	return 0, nil
}

func (e *Etcd3) Name() string { return "etcd3" }
