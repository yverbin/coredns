package kubernetes

import (
	"errors"


	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	api "k8s.io/client-go/pkg/api/v1"
)

func nodeNameIndexFunc(obj interface{}) ([]string, error) {
	s, ok := obj.(*api.Node)
	if !ok {
		return nil, errors.New("obj was not an *api.Node")
	}
	return []string{s.ObjectMeta.Name}, nil
}

func nodeListFunc(c *kubernetes.Clientset, s *labels.Selector) func(meta.ListOptions) (runtime.Object, error) {
	return func(opts meta.ListOptions) (runtime.Object, error) {
		if s != nil {
			opts.LabelSelector = (*s).String()
		}
		listV1, err := c.Nodes().List(opts)
		if err != nil {
			return nil, err
		}
		return listV1, err
	}
}

func nodeWatchFunc(c *kubernetes.Clientset, s *labels.Selector) func(options meta.ListOptions) (watch.Interface, error) {
	return func(options meta.ListOptions) (watch.Interface, error) {
		if s != nil {
			options.LabelSelector = (*s).String()
		}
		w, err := c.Nodes().Watch(options)
		if err != nil {
			return nil, err
		}
		return w, nil
	}
}

func (dns *dnsControl) NodeIndex(idx string) (nodes []*api.Node) {
	if dns.nodeLister == nil {
		return nil
	}
	os, err := dns.nodeLister.ByIndex(nodeNameIndex, idx)
	if err != nil {
		return nil
	}
	for _, o := range os {
		s, ok := o.(*api.Node)
		if !ok {
			continue
		}
		nodes = append(nodes, s)
	}
	return nodes
}

