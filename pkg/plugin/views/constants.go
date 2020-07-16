package views

const (
	AnnotationPipelineType           = "jenkins.io/pipelineType"
	AnnotationValuePipelineTypeMetaa = "meta"

	// AnnotationHome the link to the chart home page
	AnnotationHome = "jenkins.io/home"

	// AnnotationHost used to indicate the host if using NodePort Ingress resources on premise without a LoadBalancer
	AnnotationHost = "jenkins.io/host"

	// AnnotationIcon used to annotate a Deployment with its icon
	AnnotationIcon = "jenkins.io/icon"

	// LabelHelmChart label used by helm for the name of the chart
	LabelHelmChart = "helm.sh/chart"

	// LabelJXChart label for the name of the chart for jx deployments
	LabelJXChart = "chart"

	// LabelAppName standard label for the app name
	// see https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/
	// see https://helm.sh/docs/chart_best_practices/labels/
	LabelAppName = "app.kubernetes.io/name"

	// LabelAppVersion standard label for the app version
	// see https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/
	// see https://helm.sh/docs/chart_best_practices/labels/
	LabelAppVersion = "app.kubernetes.io/version"

	// DefaultIcon default icon image if none specified
	DefaultIcon = "https://clarity.design/.netlify/functions/download-icon?set=technology&shape=container-line"
)
