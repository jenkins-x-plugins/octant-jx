package viewhelpers

import (
	"context"

	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/store"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
)

// GetResourceByName looks up the resource by name
func GetResourceByName(ctx context.Context, client service.Dashboard, apiVersion, kind, name, ns string) (*unstructured.Unstructured, error) {
	u, err := client.Get(ctx, store.Key{
		APIVersion: apiVersion,
		Kind:       kind,
		Namespace:  ns,
		Name:       name,
	})
	if err != nil {
		return nil, err
	}
	return u, nil
}

// ListResourcesBySelector lists resources using a selector
func ListResourcesBySelector(ctx context.Context, client service.Dashboard, apiVersion, kind, ns string, selector labels.Set) (*unstructured.UnstructuredList, error) {
	ul, err := client.List(ctx, store.Key{
		APIVersion: apiVersion,
		Kind:       kind,
		Namespace:  ns,
		Selector:   &selector,
	})
	if err != nil {
		return nil, err
	}
	answer := &unstructured.UnstructuredList{}
	for _, u := range ul.Items {
		l := u.GetLabels()
		if l != nil {
			if MatchesSelector(l, selector) {
				answer.Items = append(answer.Items, u)
			}
		}
	}
	return answer, nil
}

// MatchesSelector returns true if the given labels match the selector
func MatchesSelector(labels map[string]string, selector labels.Set) bool {
	for k, v := range selector {
		if labels[k] != v {
			return false
		}
	}
	return true
}
