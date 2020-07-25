package admin

import (
	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/jenkins-x/octant-jx/pkg/plugin"
)

const (
	PluginName = "ojx" // This should stay lowercase for routing purposes

	OverviewPath        = "overview"
	BootPipelinesPath   = "pipelines/boot"
	FailedPipelinesPath = "pipelines/failed"
	HealthPath          = "health"
	BootJobsPath        = "jobs/boot"
	GCPipelineJobsPath  = "jobs/gcpipeline"
	GCPodJobsPath       = "jobs/gcpod"
	GCPreviewJobsPath   = "jobs/gcpreview"
	UpgradeJobsPath     = "jobs/upgrades"
	WorkspacesPath      = "workspaces"
)

var (
	// RootBreadcrumb the root breadcrumb for the operator plugin
	RootBreadcrumb = viewhelpers.ToMarkdownLink("JX OPS", OverviewLink())
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
