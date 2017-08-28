# etcd3

*etcd3* enables reading zone data from an etcd3 instance. The data in etcd3 has to be encoded as
a [message](https://github.com/skynetservices/skydns/blob/2fcff74cdc9f9a7dd64189a447ef27ac354b725f/msg/service.go#L26)
like [SkyDNS](https://github.com/skynetservices/skydns).

The etcd middleware makes extensive use of the proxy middleware to forward and query other servers
in the network.

## Syntax

~~~
etcd [ZONES...]
~~~

* **ZONES** zones etcd should be authoritative for.

The path will default to `/skydns` the local etcd proxy (http://localhost:2379).
If no zones are specified the block's zone will be used as the zone.

If you want to `round robin` A and AAAA responses look at the `loadbalance` middleware.

