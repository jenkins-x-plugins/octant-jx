package views // import "github.com/jenkins-x/octant-jx/pkg/plugin/views"

import (
	"context"
	"strings"

	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"github.com/jenkins-x/octant-jx/pkg/common/pipelines"

	v1 "github.com/jenkins-x/jx-api/v3/pkg/apis/jenkins.io/v1"
	"github.com/jenkins-x/octant-jx/pkg/common/links"
	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx"
	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/jenkins-x/octant-jx/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
)

// TODO terminals don't currently render in plugins in octant
// so lets disable until they do work
const terminalsWork = false

func BuildPipelineView(request service.Request, pluginContext pluginctx.Context) (component.Component, error) {
	name := strings.TrimPrefix(request.Path(), "/")
	name = strings.TrimPrefix(name, plugin.PipelinesPath+"/")
	ctx := request.Context()
	client := request.DashboardClient()
	ns := pluginContext.Namespace

	log.Logger().Debugf("BuildPipelineView querying for PipelineActivity %s in namespace %s\n", name, ns)

	pa, err := pipelines.GetPipeline(ctx, client, ns, name)
	if err != nil {
		log.Logger().Info(err)
		return nil, err
	}

	// lets try find the pod for the pipeline
	pod, err := findPodForPipeline(ctx, client, pluginContext, pa)
	if err != nil {
		log.Logger().Info(err)
	}

	s := &pa.Spec
	breadcrumbs := []string{
		plugin.RootBreadcrumb,
		viewhelpers.ToMarkdownLink("Pipelines", plugin.GetPipelinesLink()),
		s.GitOwner,
		s.GitRepository,
		ToNameMarkdown(pa),
	}
	if pod != nil {
		podName := pod.GetName()
		//nolint:gocritic
		breadcrumbs = append(breadcrumbs, viewhelpers.ToMarkdownLink("Pod", links.GetPodLink(ns, podName)))
		breadcrumbs = append(breadcrumbs, viewhelpers.ToMarkdownLink("Steps", plugin.GetPipelineContainersLink(ns, pa.Name, podName)))
		if terminalsWork {
			breadcrumbs = append(breadcrumbs, viewhelpers.ToMarkdownLink("Terminal", plugin.GetPipelineTerminalLink(ns, pa.Name, podName)))
		}
	}
	breadcrumbs = append(breadcrumbs, viewhelpers.ToMarkdownLink("Logs", plugin.GetPipelineLogLink(pa.Name)))
	header := viewhelpers.NewMarkdownText(viewhelpers.ToBreadcrumbMarkdown(breadcrumbs...))

	detailSummarySections := []component.SummarySection{
		{Header: "Status", Content: ToPipelineLastStepStatus(pa, true, true)},
		{Header: "Source", Content: ToRepository(pa)},
	}
	statusSummarySections := []component.SummarySection{
		{Header: "Started", Content: viewhelpers.NewMarkdownText(ToPipelineStartCompleteTimeMarkdown(pa))},
	}
	if pa.Spec.CompletedTimestamp != nil {
		statusSummarySections = append(statusSummarySections, component.SummarySection{Header: "Duration", Content: ToDuration(pa)})
	}
	detailsSummary := component.NewSummary("Status", detailSummarySections...)
	statusSummary := component.NewSummary("Timings", statusSummarySections...)

	notesCard := component.NewCard(nil)
	notesCard.SetBody(ToStepsView(pa, pod))

	flexLayout := component.NewFlexLayout("")
	flexLayout.AddSections(component.FlexLayoutSection{
		{Width: component.WidthFull, View: header},
		{Width: component.WidthHalf, View: detailsSummary},
		{Width: component.WidthHalf, View: statusSummary},
		{Width: component.WidthFull, View: notesCard},
	})
	return flexLayout, nil
}

func findPodForPipeline(ctx context.Context, client service.Dashboard, pluginContext pluginctx.Context, pa *v1.PipelineActivity) (*unstructured.Unstructured, error) {
	if pa.Labels != nil {
		podName := pa.Labels["podName"]
		if podName != "" {
			u, err := viewhelpers.GetResourceByName(ctx, client, "v1", "Pod", podName, pluginContext.Namespace)
			if err == nil {
				return u, nil
			}
			log.Logger().Warnf("failed to find pod %s/%s with error %s", pluginContext.Namespace, podName, err.Error())
		}
	}
	s := &pa.Spec
	selector := labels.Set{
		"branch":     s.GitBranch,
		"build":      s.Build,
		"owner":      s.GitOwner,
		"repository": s.GitRepository,
	}
	ul, err := viewhelpers.ListResourcesBySelector(ctx, client, "v1", "Pod", pluginContext.Namespace, selector)
	if err != nil {
		return nil, err
	}
	if len(ul.Items) > 0 {
		for i, u := range ul.Items {
			l := u.GetLabels()
			if l[AnnotationPipelineType] != AnnotationValuePipelineTypeMetaa {
				return &ul.Items[i], nil
			}
		}
		return &ul.Items[0], nil
	}
	log.Logger().Infof("could not find pod for PipelineActivity %s with selector %s", pa.Name, selector.String())
	return nil, nil
}
