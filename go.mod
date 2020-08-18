module github.com/jenkins-x/octant-jx

go 1.13

require (
	cloud.google.com/go v0.47.0 // indirect
	github.com/Azure/go-autorest/autorest v0.10.0 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.8.3 // indirect
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32 // indirect
	github.com/golang/groupcache v0.0.0-20191027212112-611e8accdfc9 // indirect
	github.com/golang/protobuf v1.3.3 // indirect
	github.com/googleapis/gnostic v0.3.1 // indirect
	github.com/gorilla/mux v1.7.3 // indirect
	github.com/hashicorp/go-hclog v0.9.2 // indirect
	github.com/hashicorp/go-plugin v1.0.1 // indirect
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/jenkins-x/jx-api v0.0.17
	github.com/jenkins-x/jx-helpers v1.0.44
	github.com/jenkins-x/jx-logging v0.0.11
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/mitchellh/mapstructure v1.2.2 // indirect
	github.com/pelletier/go-toml v1.4.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.6.1
	github.com/vmware-tanzu/octant v0.12.2-0.20200506154048-420def050373
	go.uber.org/atomic v1.5.1 // indirect
	go.uber.org/multierr v1.4.0 // indirect
	go.uber.org/zap v1.13.0 // indirect
	golang.org/x/lint v0.0.0-20200130185559-910be7a94367 // indirect
	golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	gomodules.xyz/jsonpatch/v2 v2.1.0 // indirect
	google.golang.org/genproto v0.0.0-20191028173616-919d9bdd9fe6 // indirect
	helm.sh/helm/v3 v3.2.4
	k8s.io/api v0.18.6
	k8s.io/apimachinery v0.18.6
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible // indirect
	knative.dev/pkg v0.0.0-20200423200331-4945766b290c
	sigs.k8s.io/yaml v1.2.0
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v14.0.1+incompatible
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20191016225534-b1267f8c42b4
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190620085101-78d2af792bab
)
