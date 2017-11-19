package kubernetes

import (
	"strings"

	"github.com/coredns/coredns/plugin/etcd/msg"
	api "k8s.io/client-go/pkg/api/v1"
)

func (k *Kubernetes) findNodes(r recordRequest, zone string) (ret []msg.Service, err error) {
	zonePath := msg.Path(zone, "coredns")
	err = errNoItems
	nodeList := k.APIConn.NodeIndex(r.endpoint)
	for _, node := range nodeList {
		if node.Name == r.endpoint {
			for _, nodeAddr := range node.Status.Addresses {
				if string(nodeAddr.Type) == "InternalIP" {
					nodeA := msg.Service{Host: nodeAddr.Address, TTL: k.ttl}
					nodeA.Key = strings.Join([]string{zonePath, Nodes, node.Name}, "/")
					ret = append(ret, nodeA)
					err = nil
					return
				}
			}
		}
	}
	return
}

func (k *Kubernetes) findIngress(r recordRequest, zone string) (nodes []msg.Service, err error) {
	zonePath := msg.Path(zone, "coredns")
	err = errNoItems
	if wildcard(r.service) || wildcard(r.namespace) {
		return
	}
	if r.endpoint == "" {
		return
	}

	nodeList := k.APIConn.NodeIndex(r.endpoint)
	if nodeList == nil || len(nodeList) == 0 {
		return
	}

	var (
		endpointsListFunc func() []*api.Endpoints
		endpointsList     []*api.Endpoints
		serviceList       []*api.Service
	)
	idx := r.service + "." + r.namespace
	serviceList = k.APIConn.SvcIndex(idx)
	endpointsListFunc = func() []*api.Endpoints { return k.APIConn.EpIndex(idx) }

	for _, svc := range serviceList {
		if !(match(r.namespace, svc.Namespace) && match(r.service, svc.Name)) {
			continue
		}

		if _, ok := svc.Annotations["exportNodesDomain"]; !ok {
			continue
		}
		if endpointsList == nil {
			endpointsList = endpointsListFunc()
		}
		for _, ep := range endpointsList {
			if ep.ObjectMeta.Name != svc.Name || ep.ObjectMeta.Namespace != svc.Namespace {
				continue
			}
			for _, eps := range ep.Subsets {
				for _, addr := range eps.Addresses {
					nodeList := k.APIConn.NodeIndex(*addr.NodeName)
					if nodeList == nil || len(nodeList) == 0 {
						continue
					}

					for _, node := range nodeList {
						if node.Name == r.endpoint {
							n := msg.Service{ Key: strings.Join([]string{zonePath, Ingress, svc.Namespace, svc.Name, node.Name}, "/"),
								Host: strings.Join([]string{node.Name, svc.Annotations["exportNodesDomain"]}, "."),
								TTL:  k.ttl}

							nodes = append(nodes, n)
							err = nil
							return nodes, err
						}
					}

				}

			}
		}
	}
	return

}

const (
	// Nodes is the DNS schema for kubernetes services
	Nodes = "nodes"
	// Ingress is the DNS schema for kubernetes pods
	Ingress = "ingress"
)
