module github.com/jenkins-x/octant-jx

go 1.15

require (
	cloud.google.com/go/container v1.9.0 // indirect
	cloud.google.com/go/monitoring v1.9.0 // indirect
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/jenkins-x/jx-api/v4 v4.0.25
	github.com/jenkins-x/jx-helpers/v3 v3.0.84
	github.com/jenkins-x/jx-logging/v3 v3.0.3
	github.com/jenkins-x/jx-pipeline v0.0.109
	github.com/jenkins-x/jx-preview v0.0.160
	github.com/jenkins-x/jx-secret v0.0.230
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.8.1
	github.com/tektoncd/pipeline v0.20.0
	github.com/vmware-tanzu/octant v0.17.0
	helm.sh/helm/v3 v3.10.3
	k8s.io/api v0.25.2
	k8s.io/apimachinery v0.25.2
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	sigs.k8s.io/yaml v1.3.0
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v14.0.1+incompatible

	// override the go-scm from tekton
	github.com/jenkins-x/go-scm => github.com/jenkins-x/go-scm v1.5.223
	github.com/tektoncd/pipeline => github.com/jenkins-x/pipeline v0.3.2-0.20210118090417-1e821d85abf6
	k8s.io/api => k8s.io/api v0.19.3
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.3

	k8s.io/client-go => k8s.io/client-go v0.19.3
)
