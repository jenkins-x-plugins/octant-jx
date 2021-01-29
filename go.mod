module github.com/jenkins-x/octant-jx

go 1.15

require (
	//github.com/Azure/go-autorest/autorest v0.10.0 // indirect
	//github.com/Azure/go-autorest/autorest/adal v0.8.3 // indirect
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/hashicorp/go-plugin v1.0.1 // indirect
	github.com/hashicorp/vault v1.2.3 // indirect
	github.com/jenkins-x/gen-crd-api-reference-docs v0.1.6 // indirect
	github.com/jenkins-x/go-scm v1.5.202 // indirect
	github.com/jenkins-x/jx-api/v4 v4.0.14
	github.com/jenkins-x/jx-gitops v0.0.508 // indirect
	github.com/jenkins-x/jx-helpers/v3 v3.0.47
	github.com/jenkins-x/jx-logging/v3 v3.0.2
	github.com/jenkins-x/jx-pipeline v0.0.75
	github.com/jenkins-x/jx-preview v0.0.138
	github.com/jenkins-x/jx-secret v0.0.198
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.6.1
	github.com/tektoncd/pipeline v0.16.3
	github.com/vmware-tanzu/octant v0.16.3
	go.mozilla.org/sops v0.0.0-20190912205235-14a22d7a7060 // indirect
	golang.org/x/build v0.0.0-20190111050920-041ab4dc3f9d // indirect
	helm.sh/helm/v3 v3.5.1
	k8s.io/api v0.20.1
	k8s.io/apimachinery v0.20.1
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	sigs.k8s.io/yaml v1.2.0
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v14.0.1+incompatible

	github.com/tektoncd/pipeline => github.com/jenkins-x/pipeline v0.0.0-20201002150609-ca0741e5d19a
	k8s.io/api => k8s.io/api v0.19.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.2

	k8s.io/client-go => k8s.io/client-go v0.19.2
)
