package views // import "github.com/jenkins-x/octant-jx/pkg/plugin/views"

import (
	"fmt"
	"strings"
	"time"

	"github.com/jenkins-x/jx-logging/v3/pkg/log"

	"github.com/jenkins-x/octant-jx/pkg/common/links"
	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx"
	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/jenkins-x/octant-jx/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	corev1 "k8s.io/api/core/v1"
)

const useTerminal = true

func BuildPipelineTerminalView(request service.Request, pluginContext pluginctx.Context) (component.Component, error) {
	paths := strings.Split(strings.TrimSuffix(request.Path(), "/"), "/")
	pipelineName := ""
	pl := len(paths)
	if pl > 1 {
		pipelineName = paths[pl-2]
	}
	name := paths[pl-1]
	ctx := request.Context()
	client := request.DashboardClient()
	ns := pluginContext.Namespace

	log.Logger().Infof("BuildPipelineTerminalView querying for Pipeline %s Pod %s in namespace %s\n", pipelineName, name, ns)

	u, err := viewhelpers.GetResourceByName(ctx, client, "v1", "Pod", name, ns)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return component.NewText(fmt.Sprintf("Error: Pod %s not found in namespace %s", name, ns)), nil
	}

	pod := &corev1.Pod{}
	err = viewhelpers.ToStructured(u, &pod)
	if err != nil {
		log.Logger().Info(err)
		return component.NewText(fmt.Sprintf("Error: failed to load Pod %s not found in namespace %s", name, ns)), nil
	}
	containers := pod.Spec.Containers
	if len(containers) == 0 {
		return component.NewText(fmt.Sprintf("Error: no containers for Pod %s found in namespace %s", name, ns)), nil
	}
	lastContainer := containers[len(containers)-1]

	breadcrumbs := []string{
		plugin.RootBreadcrumb,
		viewhelpers.ToMarkdownLink("Pipelines", plugin.GetPipelinesLink()),
	}
	if pipelineName != "" {
		breadcrumbs = append(breadcrumbs, viewhelpers.ToMarkdownLink("Pipeline", plugin.GetPipelineLink(pipelineName)))
	}
	breadcrumbs = append(breadcrumbs, viewhelpers.ToMarkdownLink("Pod", links.GetPodLink(ns, name)), "Terminal")
	header := viewhelpers.NewMarkdownText(viewhelpers.ToBreadcrumbMarkdown(breadcrumbs...))

	var terminal component.Component
	containerName := lastContainer.Name
	if useTerminal {
		details := component.TerminalDetails{
			Container: containerName,
			Command:   "/bin/sh",
			CreatedAt: time.Now(),
			Active:    true,
			//UUID:      name + "-" + containerName,
		}
		terminal = component.NewTerminal(ns, name, name, nil, details)
		return terminal, nil
	} else {
		terminal = viewhelpers.NewMarkdownText(fmt.Sprintf("this would be a terminal for pod %s container %s", name, containerName))
	}

	flexLayout := component.NewFlexLayout("")
	flexLayout.AddSections(component.FlexLayoutSection{
		{Width: component.WidthFull, View: header},
		{Width: component.WidthFull, View: terminal},
	})
	return flexLayout, nil
}
