package views // import "github.com/jenkins-x/octant-jx/pkg/plugin/views"

import (
	"context"
	"log"
	"strings"

	v1 "github.com/jenkins-x/jx-api/pkg/apis/jenkins.io/v1"
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

	//log.Printf("BuildPipelineView querying for PipelineActivity %s in namespace %s\n", name, ns)

	u, err := viewhelpers.GetResourceByName(ctx, client, "jenkins.io/v1", "PipelineActivity", name, ns)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return component.NewText("Error: pipeline not found"), nil
	}

	pa, err := viewhelpers.ToPipelineActivity(u)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// lets try find the pod for the pipeline
	pod, err := findPodForPipeline(ctx, client, pluginContext, pa)
	if err != nil {
		log.Println(err)
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
		breadcrumbs = append(breadcrumbs, viewhelpers.ToMarkdownLink("Pod", links.GetPodLink(ns, podName)))
		breadcrumbs = append(breadcrumbs, viewhelpers.ToMarkdownLink("Steps", plugin.GetPipelineContainersLink(ns, pa.Name, podName)))
		if terminalsWork {
			breadcrumbs = append(breadcrumbs, viewhelpers.ToMarkdownLink("Terminal", plugin.GetPipelineTerminalLink(ns, pa.Name, podName)))
		}
	}
	breadcrumbs = append(breadcrumbs, viewhelpers.ToMarkdownLink("Logs", plugin.GetPipelineLogLink(pa.Name)))
	header := component.NewMarkdownText(viewhelpers.ToBreadcrumbMarkdown(breadcrumbs...))

	detailSummarySections := []component.SummarySection{
		{"Status", ToPipelineLastStepStatus(pa, true, true)},
		{"Source", ToRepository(pa)},
	}
	statusSummarySections := []component.SummarySection{
		{"Started", component.NewMarkdownText(ToPipelineStartCompleteTimeMarkdown(pa))},
	}
	if pa.Spec.CompletedTimestamp != nil {
		statusSummarySections = append(statusSummarySections, component.SummarySection{"Duration", ToDuration(pa)})
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
	log.Printf("could not find pod for PipelineActivity %s with selector %s", pa.Name, selector.String())
	return nil, nil
}
