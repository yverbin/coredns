package etcd3

import (
	"crypto/tls"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/middleware"
	"github.com/coredns/coredns/middleware/pkg/dnsutil"
	mwtls "github.com/coredns/coredns/middleware/pkg/tls"
	"github.com/coredns/coredns/middleware/proxy"

	"github.com/coreos/etcd/clientv3"
	"github.com/mholt/caddy"
)

func init() {
	caddy.RegisterPlugin("etcd3", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	e, stubzones, err := etcd3Parse(c)
	if err != nil {
		return middleware.Error("etcd3", err)
	}

	/*
		if stubzones {
			c.OnStartup(func() error {
				e.UpdateStubZones()
				return nil
			})
		}
	*/

	dnsserver.GetConfig(c).AddMiddleware(func(next middleware.Handler) middleware.Handler {
		e.Next = next
		return e
	})

	return nil
}

func etcd3Parse(c *caddy.Controller) (*Etcd3, bool, error) {
	stub := make(map[string]proxy.Proxy)
	e3 := Etcd3{
		PathPrefix: "skydns",
		Stubmap:    &stub,
	}
	var (
		tlsConfig *tls.Config
		err       error
		endpoints = []string{defaultEndpoint}
		stubzones = false
	)
	for c.Next() {
		e3.Zones = c.RemainingArgs()
		if len(e3.Zones) == 0 {
			e3.Zones = make([]string, len(c.ServerBlockKeys))
			copy(e3.Zones, c.ServerBlockKeys)
		}
		for i, str := range e3.Zones {
			e3.Zones[i] = middleware.Host(str).Normalize()
		}

		if c.NextBlock() {
			for {
				switch c.Val() {
				case "stubzones":
					stubzones = true
				case "fallthrough":
					e3.Fallthrough = true
				case "path":
					if !c.NextArg() {
						return nil, false, c.ArgErr()
					}
					e3.PathPrefix = c.Val()
				case "endpoint":
					args := c.RemainingArgs()
					if len(args) == 0 {
						return nil, false, c.ArgErr()
					}
					endpoints = args
				case "upstream":
					args := c.RemainingArgs()
					if len(args) == 0 {
						return nil, false, c.ArgErr()
					}
					ups, err := dnsutil.ParseHostPortOrFile(args...)
					if err != nil {
						return nil, false, err
					}
					e3.Proxy = proxy.NewLookup(ups)
				case "tls": // cert key cacertfile
					args := c.RemainingArgs()
					tlsConfig, err = mwtls.NewTLSConfigFromArgs(args...)
					if err != nil {
						return nil, false, err
					}
				default:
					if c.Val() != "}" {
						return nil, false, c.Errf("unknown property '%s'", c.Val())
					}
				}
			}

		}
		client, err := newEtcd3KV(endpoints, tlsConfig)
		if err != nil {
			return nil, false, err
		}
		e3.kv = client
		e3.endpoints = endpoints

		return &e3, stubzones, nil
	}
	return nil, false, nil
}

func newEtcd3KV(endpoints []string, cc *tls.Config) (clientv3.KV, error) {
	cli, err := clientv3.New(clientv3.Config{
		DialTimeout: defaultTimeout,
		Endpoints:   endpoints,
		TLS:         cc,
	})
	if err != nil {
		return nil, err
	}
	kv := clientv3.NewKV(cli)
	return kv, nil
}

const defaultEndpoint = "http://localhost:2379"
