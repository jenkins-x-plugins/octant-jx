package views

import (
	"fmt"
	"strings"

	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"github.com/jenkins-x/octant-jx/pkg/common/pipelines"

	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx"
	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/jenkins-x/octant-jx/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func BuildPipelineLog(request service.Request, pluginContext pluginctx.Context) (component.Component, error) {
	path := strings.TrimPrefix(request.Path(), "/")
	path = strings.TrimPrefix(path, plugin.LogsPath+"/")

	paths := strings.Split(path, "/")
	name := paths[0]

	ctx := request.Context()
	client := request.DashboardClient()

	log.Logger().Infof("BuildPipelineLog querying for PipelineActivity %s\n", name)

	pa, err := pipelines.GetPipeline(ctx, client, pluginContext.Namespace, name)
	if err != nil {
		log.Logger().Info(err)
		return nil, err
	}

	s := &pa.Spec
	header := viewhelpers.NewMarkdownText(viewhelpers.ToBreadcrumbMarkdown(
		plugin.RootBreadcrumb,
		viewhelpers.ToMarkdownLink("Pipelines", plugin.GetPipelinesLink()),
		s.GitOwner,
		s.GitRepository,
		viewhelpers.ToMarkdownLink(ToNameMarkdown(pa), plugin.GetPipelineLink(pa.Name)),
		"Logs"))

	// lets try find the pod for the pipeline
	var logsView component.Component
	pod, err := findPodForPipeline(ctx, client, pluginContext, pa)
	if err != nil {
		log.Logger().Info(err)
	}
	if pod != nil {
		ns := pa.Namespace
		podName := pod.GetName()
		if len(paths) > 1 {
			logsView, err = viewhelpers.ViewPipelineLogs(ns, podName, paths[1])
		} else {
			logsView, err = viewhelpers.ViewPipelineLogs(ns, podName)
		}
		if err != nil {
			log.Logger().Info(err)
			logsView = component.NewText(fmt.Sprintf("could not find pod: %s", err.Error()))
		}
	} else {
		logsView = component.NewText("could not find pod")
	}
	notesCard := component.NewCard(component.TitleFromString("Steps"))
	notesCard.SetBody(ToStepsView(pa, pod))

	flexLayout := component.NewFlexLayout("")
	flexLayout.AddSections(component.FlexLayoutSection{
		{Width: component.WidthFull, View: header},
		{Width: component.WidthFull, View: logsView},
	})
	return flexLayout, nil
}
