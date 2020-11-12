package router // import "github.com/jenkins-x/octant-jx/pkg/plugin/router"

import (
	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx"
	"github.com/jenkins-x/octant-jx/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"

	"github.com/jenkins-x/octant-jx/pkg/plugin/views"
)

type Handlers struct {
	Context *pluginctx.Context
}

func (h *Handlers) InitRoutes(router *service.Router) {
	h.handleView(router, plugin.OverviewPath, views.BuildOverview)
	h.handleView(router, plugin.AppsPath, views.BuildAppsView)
	h.handleView(router, plugin.EnvironmentsPath, views.BuildEnvironmentsView)
	h.handleView(router, plugin.EnvironmentsPath+"/*", views.BuildEnvironmentView)
	h.handleView(router, plugin.HelmPath, views.BuildHelmReleasesView)
	h.handleView(router, plugin.HelmPath+"/*", views.BuildHelmReleaseView)
	h.handleView(router, plugin.PipelinesPath, views.BuildPipelinesViewDefault)
	h.handleView(router, plugin.PipelinesRecentPath, views.BuildPipelinesViewRecent)
	h.handleView(router, plugin.LogsPath+"/*", views.BuildPipelineLog)
	h.handleView(router, plugin.PipelinesPath+"/*", views.BuildPipelineView)
	h.handleView(router, plugin.PipelineContainersPath+"/*", views.BuildPipelineContainersView)
	h.handleView(router, plugin.PipelineContainerPath+"/*", views.BuildPipelineContainerView)
	h.handleView(router, plugin.PipelineTerminalPath+"/*", views.BuildPipelineTerminalView)
	h.handleView(router, plugin.PreviewsPath, views.BuildPreviewsView)
	h.handleView(router, plugin.RepositoriesPath, views.BuildRepositoriesView)
}

func (h *Handlers) handleView(router *service.Router, path string, fn func(request service.Request, pluginContext pluginctx.Context) (component.Component, error)) {
	router.HandleFunc("/"+path, func(request service.Request) (component.ContentResponse, error) {
		view, err := fn(request, *h.Context)
		if err != nil {
			return component.EmptyContentResponse, err
		}
		response := component.NewContentResponse(nil)
		response.Add(view)
		return *response, nil
	})
}
