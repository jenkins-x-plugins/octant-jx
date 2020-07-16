package views_test

import (
	"testing"

	"github.com/jenkins-x/octant-jx/pkg/plugin/views"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/util/intstr"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	ns = "jx"
)

func TestToResourceNameVersion(t *testing.T) {
	testCases := []struct {
		Name            string
		ObjectMeta      *metav1.ObjectMeta
		ExpectedName    string
		ExpectedVersion string
	}{
		{
			Name: "lighthouse-foghorn",
			ObjectMeta: &metav1.ObjectMeta{
				Name:      "ighthouse-foghorn",
				Namespace: ns,
				Annotations: map[string]string{
					"meta.helm.sh/release-name": "lighthouse",
				},
				Labels: map[string]string{
					"app":                       "lighthouse-foghorn",
					"app.kubernetes.io/version": "0.0.563",
					"chart":                     "lighthouse-0.0.563",
				},
			},
			ExpectedName:    "foghorn",
			ExpectedVersion: "0.0.563",
		},
		{
			// see https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/
			Name: "standardLabels",
			ObjectMeta: &metav1.ObjectMeta{
				Name:      "some-random-name",
				Namespace: ns,
				Labels: map[string]string{
					"app.kubernetes.io/name":    "myapp",
					"app.kubernetes.io/version": "1.2.3",
				},
			},
			ExpectedName:    "myapp",
			ExpectedVersion: "1.2.3",
		},
		{
			Name: "useJXChartLabel",
			ObjectMeta: &metav1.ObjectMeta{
				Name:      "nexus-nexus",
				Namespace: ns,
				Labels: map[string]string{
					"app":     "nexus",
					"release": "nexus",
					"chart":   "nexus-0.1.21",
				},
			},
			ExpectedName:    "nexus",
			ExpectedVersion: "0.1.21",
		},
		{
			Name: "useHelmChartLabel",
			ObjectMeta: &metav1.ObjectMeta{
				Name:      "nexus-nexus",
				Namespace: ns,
				Labels: map[string]string{
					"helm.sh/chart": "nexus-0.1.21",
				},
			},
			ExpectedName:    "nexus",
			ExpectedVersion: "0.1.21",
		},
		{
			Name: "noLabels",
			ObjectMeta: &metav1.ObjectMeta{
				Name:      "random-deployment",
				Namespace: ns,
			},
			ExpectedName:    "random-deployment",
			ExpectedVersion: "",
		},
	}

	for _, tc := range testCases {
		require.NotNil(t, tc.ObjectMeta, "no Deployment for test %s", tc.Name)
		name, version := views.ToResourceNameVersion(tc.ObjectMeta)
		assert.Equal(t, tc.ExpectedName, name, "ToDeploymentNameVersion name for test %s", tc.Name)
		assert.Equal(t, tc.ExpectedVersion, version, "ToDeploymentNameVersion version for test %s", tc.Name)
	}
}

func TestFindIngressLinkByAppName(t *testing.T) {
	ingresses := []*v1beta1.Ingress{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "http",
				Labels:    map[string]string{},
				Namespace: ns,
			},
			Spec: v1beta1.IngressSpec{
				Rules: []v1beta1.IngressRule{
					{
						Host: "hook-jx.1.2.3.4.nip.io",
						IngressRuleValue: v1beta1.IngressRuleValue{
							HTTP: &v1beta1.HTTPIngressRuleValue{
								Paths: []v1beta1.HTTPIngressPath{
									{
										Path: "",
										Backend: v1beta1.IngressBackend{
											ServiceName: "hook",
											ServicePort: intstr.IntOrString{
												IntVal: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "http-and-https-ingress",
				Labels:    map[string]string{},
				Namespace: ns,
			},
			Spec: v1beta1.IngressSpec{
				Rules: []v1beta1.IngressRule{
					{
						Host: "hook-jx.1.2.3.4.nip.io",
						IngressRuleValue: v1beta1.IngressRuleValue{
							HTTP: &v1beta1.HTTPIngressRuleValue{
								Paths: []v1beta1.HTTPIngressPath{
									{
										Path: "",
										Backend: v1beta1.IngressBackend{
											ServiceName: "hook",
											ServicePort: intstr.IntOrString{
												IntVal: 80,
											},
										},
									},
								},
							},
						},
					},
				},
				TLS: []v1beta1.IngressTLS{
					{
						Hosts:      []string{"hook-jx.1.2.3.4.nip.io"},
						SecretName: "",
					},
				},
			},
		},
	}

	testCases := []struct {
		Name     string
		Expected string
	}{
		{
			Name:     "http",
			Expected: "http://hook-jx.1.2.3.4.nip.io",
		},
		{
			Name:     "http-and-https-ingress",
			Expected: "https://hook-jx.1.2.3.4.nip.io",
		},
	}

	for _, tc := range testCases {
		actual := views.FindIngressLinkByAppName(tc.Name, ingresses)
		assert.Equal(t, tc.Expected, actual, "FindIngressLinkByAppName for name %s", tc.Name)
	}
}
