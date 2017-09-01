# etcd3

*etcd3* enables reading zone data from an etcd3 instance. The data in etcd3 has to be encoded as
a [message](https://github.com/skynetservices/skydns/blob/2fcff74cdc9f9a7dd64189a447ef27ac354b725f/msg/service.go#L26)
like [SkyDNS](https://github.com/skynetservices/skydns).

The *etcd3* middleware makes extensive use of the proxy middleware to forward and query other servers
in the network.

## Syntax

~~~
etcd [ZONES...]
~~~

* **ZONES** zones etcd should be authoritative for.

The path will default to `/skydns` the local etcd3 proxy (http://localhost:2379).
If no zones are specified the block's zone will be used as the zone.

~~~
etcd3 [ZONES...] {
    stubzones
    fallthrough
    path PATH
    endpoint ENDPOINT...
    upstream ADDRESS...
    tls CERT KEY CACERT
}
~~~

* `stubzones` enables the stub zones feature. The stubzone is *only* done in the etcd tree located
    under the *first* zone specified.
* `fallthrough` If zone matches but no record can be generated, pass request to the next middleware.
* **PATH** the path inside etcd. Defaults to "/skydns".
* **ENDPOINT** the etcd endpoints. Defaults to "http://localhost:2397".
* `upstream` upstream resolvers to be used resolve external names found in etcd (think CNAMEs)
  pointing to external names. If you want CoreDNS to act as a proxy for clients, you'll need to add
  the proxy middleware. **ADDRESS** can be an IP address, and IP:port or a string pointing to a file
  that is structured as /etc/resolv.conf.
* `tls` followed by:
  * no arguments, if the server certificate is signed by a system-installed CA and no client cert is needed
  * a single argument that is the CA PEM file, if the server cert is not signed by a system CA and no client cert is needed
  * two arguments - path to cert PEM file, the path to private key PEM file - if the server certificate is signed by a system-installed CA and a client certificate is needed
  * three arguments - path to cert PEM file, path to client private key PEM file, path to CA PEM file - if the server certificate is not signed by a system-installed CA and client certificate is needed

## Examples
