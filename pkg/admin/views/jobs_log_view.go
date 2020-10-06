package views // import "github.com/jenkins-x/octant-jx/pkg/plugin/views"

import (
	"fmt"

	"github.com/jenkins-x/jx-logging/v3/pkg/log"

	"github.com/jenkins-x/octant-jx/pkg/admin"
	"github.com/jenkins-x/octant-jx/pkg/common/links"
	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx"
	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	"k8s.io/apimachinery/pkg/labels"
)

func BuildJobsLogViewForPath(request service.Request, pluginContext pluginctx.Context, path, jobName string) (component.Component, error) {
	config := JobsViewConfigs[path]
	selector := config.Selector
	return BuildJobsViewLogsForPathAndSelector(request, pluginContext, path, jobName, config, selector)
}

func BuildJobsViewLogsForPathAndSelector(request service.Request, pluginContext pluginctx.Context, path, jobName string, config *JobViewConfig, selector labels.Set) (component.Component, error) {
	if config == nil {
		return component.NewText(fmt.Sprintf("No view configuration found for path %s", path)), nil
	}

	ctx := request.Context()
	client := request.DashboardClient()
	ns := pluginContext.Namespace

	if jobName != "" && selector["job-name"] == "" {
		selector["job-name"] = jobName
	}
	title := config.Title
	if title == "" {
		title = "Jobs"
	}
	parentLink := viewhelpers.ToMarkdownLink(title, admin.JobsViewLink(path))
	headerText := viewhelpers.ToBreadcrumbMarkdown(admin.RootBreadcrumb, parentLink)
	if jobName != "" {
		headerText = viewhelpers.ToBreadcrumbMarkdown(headerText, viewhelpers.ToMarkdownLink("Job", links.GetJobLink(ns, jobName)))
	}

	// lets try find the pod for the pipeline
	var logsView component.Component
	pod, err := viewhelpers.FindLatestPodForSelector(ctx, client, pluginContext.Namespace, selector)
	if err != nil {
		log.Logger().Info(err)
	}
	if pod != nil {
		podName := pod.GetName()
		if ns == "" {
			ns = pod.Namespace
		}
		logsView, err = viewhelpers.ViewPipelineLogs(ns, podName)
		if err != nil {
			log.Logger().Info(err)
			logsView = component.NewText(fmt.Sprintf("could not find pod: %s", err.Error()))
		}

		headerText = viewhelpers.ToBreadcrumbMarkdown(headerText, viewhelpers.ToMarkdownLink("Pod", links.GetPodLink(ns, podName)))

	} else {
		logsView = component.NewText(fmt.Sprintf("could not find pod for selector: %s", selector.String()))
	}
	headerText = viewhelpers.ToBreadcrumbMarkdown(headerText, "Logs")

	flexLayout := component.NewFlexLayout("")
	flexLayout.AddSections(component.FlexLayoutSection{
		{Width: component.WidthFull, View: component.NewMarkdownText(headerText)},
		{Width: component.WidthFull, View: logsView},
	})
	return flexLayout, nil
}
