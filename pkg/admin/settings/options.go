package settings // import "github.com/jenkins-x/octant-jx/pkg/plugin/settings"

import (
	"github.com/jenkins-x/octant-jx/pkg/admin"
	"github.com/jenkins-x/octant-jx/pkg/admin/router"
	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx/actioners"
	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
)

func GetOptions(h *router.Handlers) []service.PluginOption {
	return []service.PluginOption{
		service.WithNavigation(
			func(_ *service.NavigationRequest) (navigation.Navigation, error) {
				return navigation.Navigation{
					Title:    "JX OPS",
					Path:     admin.PluginName + "/" + admin.OverviewPath,
					IconName: rootNavIcon,
					Children: []navigation.Navigation{
						{
							Title: "Overview",
							Path:  admin.PluginName + "/" + admin.OverviewPath,
						},
						{
							Title: "Workspaces",
							Path:  admin.PluginName + "/" + admin.WorkspacesPath,
						},
						{
							Title: "Boot Jobs",
							Path:  admin.PluginName + "/" + admin.BootJobsPath,
						},
						{
							Title: "Boot Pipelines",
							Path:  admin.PluginName + "/" + admin.BootPipelinesPath,
						},
						{
							Title: "Health",
							Path:  admin.PluginName + "/" + admin.HealthPath,
						},
						{
							Title: "GC Pipeline Jobs",
							Path:  admin.PluginName + "/" + admin.GCPipelineJobsPath,
						},
						{
							Title: "GC Pod Jobs",
							Path:  admin.PluginName + "/" + admin.GCPodJobsPath,
						},
						{
							Title: "GC Preview Jobs",
							Path:  admin.PluginName + "/" + admin.GCPreviewJobsPath,
						},
						{
							Title: "Upgrade Jobs",
							Path:  admin.PluginName + "/" + admin.UpgradeJobsPath,
						},
						{
							Title: "Failed Release Pipelines",
							Path:  admin.PluginName + "/" + admin.FailedPipelinesPath,
						},
					},
				}, nil
			},
			h.InitRoutes,
		),
		service.WithActionHandler(actioners.CreateHandler(h.Context)),
	}
}
