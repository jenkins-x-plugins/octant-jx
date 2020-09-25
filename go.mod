module github.com/jenkins-x/octant-jx

go 1.13

require (
	github.com/coreos/bbolt v1.3.3 // indirect
	github.com/coreos/etcd v3.3.17+incompatible // indirect
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32 // indirect
	github.com/go-yaml/yaml v2.1.0+incompatible // indirect
	github.com/golang/lint v0.0.0-20180702182130-06c8688daad7 // indirect
	github.com/gotestyourself/gotestyourself v2.2.0+incompatible // indirect
	github.com/hashicorp/go-plugin v1.0.1 // indirect
	github.com/jenkins-x/jx-api v0.0.20
	github.com/jenkins-x/jx-helpers v1.0.70
	github.com/jenkins-x/jx-logging v0.0.11
	github.com/klauspost/cpuid v1.2.2 // indirect
	github.com/knative/build v0.1.2 // indirect
	github.com/natefinch/lumberjack v2.0.0+incompatible // indirect
	github.com/nats-io/gnatsd v1.4.1 // indirect
	github.com/nats-io/go-nats v1.7.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.6.1
	github.com/tektoncd/pipeline v0.16.3
	github.com/vmware-tanzu/octant v0.12.2-0.20200506154048-420def050373
	gopkg.in/yaml.v1 v1.0.0-20140924161607-9f9df34309c0 // indirect
	helm.sh/helm/v3 v3.2.4
	k8s.io/api v0.18.6
	k8s.io/apimachinery v0.18.6
	knative.dev/pkg v0.0.0-20200702222342-ea4d6e985ba0
	sigs.k8s.io/structured-merge-diff v1.0.1 // indirect
	sigs.k8s.io/testing_frameworks v0.1.1 // indirect
	sigs.k8s.io/yaml v1.2.0
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v14.0.1+incompatible
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20191016225534-b1267f8c42b4
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190620085101-78d2af792bab
)
