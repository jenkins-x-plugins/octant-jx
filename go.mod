module github.com/jenkins-x/octant-jx

go 1.13

require (
	github.com/Azure/go-autorest/autorest v0.10.0 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.8.3 // indirect
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32 // indirect
	github.com/hashicorp/go-hclog v0.9.2 // indirect
	github.com/hashicorp/go-plugin v1.0.1 // indirect
	github.com/jenkins-x/jx-api v0.0.20
	github.com/jenkins-x/jx-helpers v1.0.70
	github.com/jenkins-x/jx-logging v0.0.11
	github.com/mitchellh/mapstructure v1.2.2 // indirect
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.6.1
	github.com/tektoncd/pipeline v0.14.2
	github.com/vmware-tanzu/octant v0.12.2-0.20200506154048-420def050373
	helm.sh/helm/v3 v3.2.4
	k8s.io/api v0.18.6
	k8s.io/apimachinery v0.19.2
	knative.dev/pkg v0.0.0-20200528142800-1c6815d7e4c9
	sigs.k8s.io/yaml v1.2.0
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v14.0.1+incompatible
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20191016225534-b1267f8c42b4
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190620085101-78d2af792bab
)
