package admin

import (
	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx"
	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/jenkins-x/octant-jx/pkg/plugin"
)

const (
	PluginName = "ojx" // This should stay lowercase for routing purposes

	OverviewPath        = "overview"
	BootPipelinesPath   = "pipelines/boot"
	FailedPipelinesPath = "pipelines/failed"
	HealthPath          = "health"
	SecretsPath         = "secrets"
	BootJobsPath        = "jobs/boot"
	GCPipelineJobsPath  = "jobs/gcpipeline"
	GCPodJobsPath       = "jobs/gcpod"
	GCPreviewJobsPath   = "jobs/gcpreview"
	UpgradeJobsPath     = "jobs/upgrades"
	WorkspacesPath      = "workspaces"

	BootJobNamespace = "jx-git-operator"
)

var (
	// RootBreadcrumb the root breadcrumb for the operator plugin
	RootBreadcrumb = viewhelpers.ToMarkdownLink("JX OPS", OverviewLink())

	// BootPluginContext the context used for boot jobs
	BootPluginContext = pluginctx.Context{
		Namespace: BootJobNamespace,
		Composite: false,
	}
)

func OverviewLink() string {
	return JobsViewLink(OverviewPath)
}

func JobsViewLink(path string) string {
	return plugin.PathPrefix + "/" + PluginName + "/" + path
}

func JobsLogsViewLink(path, jobName string) string {
	link := JobsViewLink(path) + "/logs"
	if jobName != "" {
		return link + "/" + jobName
	}
	return link
}
