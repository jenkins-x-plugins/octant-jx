package viewhelpers

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/jenkins-x/jx-logging/pkg/log"

	v1 "github.com/jenkins-x/jx-api/pkg/apis/jenkins.io/v1"
	"github.com/pkg/errors"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"knative.dev/pkg/apis/duck"
)

// ToPipelineActivity converts to a PipelineActivity
func ToPipelineActivity(u *unstructured.Unstructured) (*v1.PipelineActivity, error) {
	r := &v1.PipelineActivity{}
	err := ToStructured(u, r)
	return r, err
}

// ToPod converts to a Pod
func ToPod(u *unstructured.Unstructured) (*corev1.Pod, error) {
	r := &corev1.Pod{}
	err := ToStructured(u, r)
	return r, err
}

// ToSecret converts to a Secret
func ToSecret(u *unstructured.Unstructured) (*corev1.Secret, error) {
	r := &corev1.Secret{}
	err := ToStructured(u, r)
	return r, err
}

// ToStructured converts an unstructured object to a pointer to a structured type
func ToStructured(u *unstructured.Unstructured, structured interface{}) error {
	if err := duck.FromUnstructured(u, structured); err != nil {
		return errors.Wrapf(err, "failed to convert unstructured object to %#v", structured)
	}
	return nil
}

// ResourceTimeLessThan returns whether the first time is less than the second time
func ResourceTimeLessThan(t1, t2 *metav1.Time) bool {
	if t1 == nil {
		if t2 == nil {
			return false
		} else {
			return true
		}
	}
	if t2 == nil {
		return false
	}
	return t1.Time.Before(t2.Time)
}

func ToTimestamp(t *metav1.Time) component.Component {
	if t == nil {
		return component.NewText("")
	}
	return component.NewTimestamp(t.Time)
}

func ToDurationString(u time.Time) string {
	return time.Since(u).Truncate(time.Second).String()
}

func ToDurationMarkdown(time time.Time, titlePrefix string) string {
	return fmt.Sprintf(`<span title="%s%s">%s</span>`, titlePrefix, time.String(), ToDurationString(time))
}

// ListPodsBySelector loads the pods for the given selector
func ListPodsBySelector(ctx context.Context, client service.Dashboard, namespace string, selector labels.Set) ([]*corev1.Pod, error) {
	ul, err := ListResourcesBySelector(ctx, client, "v1", "Pod", namespace, selector)
	if err != nil {
		return nil, err
	}

	pods := []*corev1.Pod{}

	for k, v := range ul.Items {
		pod, err := ToPod(&ul.Items[k])
		if err != nil {
			return pods, errors.Wrapf(err, "failed to convert pod %s", v.GetName())
		}
		pods = append(pods, pod)
	}
	return pods, nil
}

func FindLatestPodForSelector(ctx context.Context, client service.Dashboard, namespace string, selector labels.Set) (*corev1.Pod, error) {
	pods, err := ListPodsBySelector(ctx, client, namespace, selector)
	if err != nil {
		return nil, err
	}
	if len(pods) == 0 {
		log.Logger().Infof("could not find pod in namespace %s with selector %s", namespace, selector.String())
		return nil, nil
	}

	// sort into first created
	sort.Slice(pods, func(i, j int) bool {
		p1 := pods[i]
		p2 := pods[j]
		return p1.CreationTimestamp.After(p2.CreationTimestamp.Time)
	})
	return pods[0], nil
}
