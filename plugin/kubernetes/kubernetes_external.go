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
			if ip, ok := getNodeIP(node); ok {
				nodeA := msg.Service{Host: ip, TTL: k.ttl}
				nodeA.Key = strings.Join([]string{zonePath, Nodes, node.Name}, "/")
				ret = append(ret, nodeA)
				err = nil
				return
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

	idx := r.service + "." + r.namespace
	serviceList := k.APIConn.SvcIndex(idx)

	if serviceList == nil || len(serviceList) == 0 {
		return
	}
	var (
		nodeList []*api.Node
		svc      *api.Service
	)
	svc = serviceList[0]
	if !(match(r.namespace, svc.Namespace) && match(r.service, svc.Name)) {
		return
	}

	switch svc.Annotations["exportNodesMode"] {
	case ExportModeCName:
		if r.endpoint == "" || svc.Annotations["exportNodesDomain"] == "" {
			return
		}
		nodeList = k.APIConn.NodeIndex(r.endpoint)
	case ExportModeHost:
		if r.endpoint != "" {
			return
		}
		nodeList = k.APIConn.NodesList()
	default:
		return
	}

	if nodeList == nil || len(nodeList) == 0 {
		return
	}
	endpoints := k.APIConn.EpIndex(idx)

	for _, ep := range endpoints {
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
					switch svc.Annotations["exportNodesMode"] {
					case ExportModeCName:
						if node.Name == r.endpoint {
							//export as cname
							n := msg.Service{Key: strings.Join([]string{zonePath, Ingress, svc.Namespace, svc.Name, node.Name}, "/"),
								Host: strings.Join([]string{node.Name, svc.Annotations["exportNodesDomain"]}, "."),
								TTL:  k.ttl}

							nodes = append(nodes, n)
							err = nil
							//returns A-record  if exportNodesDomain equal nodes primary zone at once
							if strings.Join(strings.Split(svc.Annotations["exportNodesDomain"], "."), "/") == strings.Join([]string{zone, Nodes}, "/") {
								if ip, ok := getNodeIP(node); ok {
									nodeA := msg.Service{Host: ip, TTL: k.ttl}
									nodeA.Key = strings.Join([]string{zonePath, Nodes, node.Name}, "/")
									nodes = append(nodes, n)
								}
							}
							return nodes, err
						}
					case ExportModeHost:
						if node.Name == *addr.NodeName {
							if ip, ok := getNodeIP(node); ok {
								//export as A records
								n := msg.Service{Key: strings.Join([]string{zonePath, Ingress, svc.Namespace, svc.Name}, "/"),
									Host: ip,
									TTL:  k.ttl}
								err = nil
								nodes = append(nodes, n)
							}

						}

					}
				}

			}
		}
	}

	return

}

func getNodeIP(node *api.Node) (ip string, ok bool) {
	for _, nodeAddr := range node.Status.Addresses {
		if string(nodeAddr.Type) == "InternalIP" {
			ip = nodeAddr.Address
			ok = true
		}
	}
	return
}

const (
	// Nodes is the DNS schema for cluster nodes
	Nodes = "nodes"
	// Ingress is the DNS schema for kubernetes ingress services
	Ingress = "ingress"
	//ExportModeHost export nodes as A records
	ExportModeHost = "host"
	//ExportModeCName export nodes as Cname records
	ExportModeCName = "cname"
)
