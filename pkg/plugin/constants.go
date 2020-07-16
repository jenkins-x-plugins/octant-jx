package plugin

import (
	"github.com/vmware-tanzu/octant/pkg/navigation"
)

const (
	PluginName = "jx" // This should stay lowercase for routing purposes

	// PathPrefix the initial prefix of all paths
	PathPrefix = "/#"

	AppsPath               = "apps"
	EnvironmentsPath       = "environments"
	HelmPath               = "helm"
	LogsPath               = "logs"
	OverviewPath           = "overview"
	PipelinesPath          = "pipelines"
	PipelineContainersPath = "pipeline/containers"
	PipelineContainerPath  = "pipeline/container"
	PipelineTerminalPath   = "pipeline/terminal"
	PipelinesRecentPath    = "pipelines-recent"
	RepositoriesPath       = "repositories"

	// RootBreadcrumb the root breadcrumb for the developer plugin
	RootBreadcrumb = `<a href="/#/jx/overview">Jenkins X</a>`
)

var (
	// Navigations the default navigations
	Navigations = []navigation.Navigation{
		{
			Title: "Apps",
			Path:  PluginName + "/" + AppsPath,
		},
		{
			Title: "Environments",
			Path:  PluginName + "/" + EnvironmentsPath,
		},
		{
			Title: "Helm",
			Path:  PluginName + "/" + HelmPath,
		},
		{
			Title: "Pipelines",
			Path:  PluginName + "/" + PipelinesPath,
		},
		{
			Title: "Pipelines: Recent",
			Path:  PluginName + "/" + PipelinesRecentPath,
		},
		{
			Title: "Repositories",
			Path:  PluginName + "/" + RepositoriesPath,
		},
	}
)
