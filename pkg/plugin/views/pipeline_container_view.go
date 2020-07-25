package views // import "github.com/jenkins-x/octant-jx/pkg/plugin/views"

import (
	"fmt"
	"strings"

	"github.com/jenkins-x/jx-logging/pkg/log"

	"github.com/jenkins-x/octant-jx/pkg/common/links"
	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx"
	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/jenkins-x/octant-jx/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	corev1 "k8s.io/api/core/v1"
)

func BuildPipelineContainerView(request service.Request, pluginContext pluginctx.Context) (component.Component, error) {
	paths := strings.Split(strings.TrimSuffix(request.Path(), "/"), "/")
	pipelineName := ""
	pl := len(paths)
	if len(paths) < 3 {
		return component.NewText("not enough values in the path"), nil

	}
	pipelineName = paths[pl-3]
	name := paths[pl-2]
	step := paths[pl-1]
	ctx := request.Context()
	client := request.DashboardClient()
	ns := pluginContext.Namespace

	log.Logger().Debugf("BuildPipelineContainersView querying for Pod %s in namespace %s\n", name, ns)

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

	breadcrumbs := []string{
		plugin.RootBreadcrumb,
		viewhelpers.ToMarkdownLink("Pipelines", plugin.GetPipelinesLink()),
	}
	if pipelineName != "" {
		breadcrumbs = append(breadcrumbs, viewhelpers.ToMarkdownLink("Pipeline", plugin.GetPipelineLink(pipelineName)))
	}
	breadcrumbs = append(breadcrumbs, viewhelpers.ToMarkdownLink("Pod", links.GetPodLink(ns, name)))
	if pipelineName != "" {
		breadcrumbs = append(breadcrumbs, viewhelpers.ToMarkdownLink("Steps", plugin.GetPipelineContainersLink(ns, pipelineName, name)))
	}
	breadcrumbs = append(breadcrumbs, step)

	header := component.NewMarkdownText(viewhelpers.ToBreadcrumbMarkdown(breadcrumbs...))

	notesCard := component.NewCard(nil)

	vc := containersViewContext{
		Request:      request,
		Namespace:    ns,
		PipelineName: pipelineName,
		PodName:      name,
	}

	found := false
	containers := pod.Spec.Containers
	//Todo: Need to evaluate the logic
	for k := range containers {
		if containers[k].Name == step {
			return ToPipelinePodContainerView(header, vc, pod, k, &containers[k]), nil
			// ToDo: unreachable code ...
			//nolint
			found = true
			break
		}

	}
	if !found {
		notesCard.SetBody(component.NewMarkdownText(fmt.Sprintf("Pod %s does not have a container called %s", pod, step)))
	}

	flexLayout := component.NewFlexLayout("")
	flexLayout.AddSections(component.FlexLayoutSection{
		{Width: component.WidthFull, View: header},
		{Width: component.WidthFull, View: notesCard},
	})
	return flexLayout, nil
}

func ToPipelinePodContainerView(header component.Component, vc containersViewContext, pod *corev1.Pod, index int, c *corev1.Container) component.Component {
	image := ToImage(c)
	commandLine := ToCommandLine(index, c)

	statusSummarySections := []component.SummarySection{
		{Header: "Name", Content: component.NewText(c.Name)},
		{Header: "Image", Content: component.NewMarkdownText(image)},
		{Header: "Working Dir", Content: component.NewText(c.WorkingDir)},
		{Header: "Command", Content: component.NewMarkdownText(fmt.Sprintf("```%s```", commandLine))},
	}
	statusSummary := component.NewSummary("Container", statusSummarySections...)

	volumesSections := []component.SummarySection{}
	if len(c.VolumeMounts) > 0 {
		vm := c.VolumeMounts
		for k := range vm {
			volumesSections = append(volumesSections, ToVolumeMountSection(pod, c, &vm[k]))
		}
	}
	viewhelpers.SortSummarySection(volumesSections)
	volumesSummary := component.NewSummary("Volume Mounts", volumesSections...)

	envSections := []component.SummarySection{}
	for _, e := range c.Env {
		envSections = append(envSections, ToEnvVarSection(e))
	}
	viewhelpers.SortSummarySection(envSections)
	envSummary := component.NewSummary("Environment Variables", envSections...)

	flexLayout := component.NewFlexLayout("")
	flexLayout.AddSections(component.FlexLayoutSection{
		{Width: component.WidthFull, View: header},
		{Width: component.WidthHalf, View: statusSummary},
		{Width: component.WidthHalf, View: volumesSummary},
		{Width: component.WidthFull, View: envSummary},
	})
	return flexLayout
}

func ToEnvVarSection(e corev1.EnvVar) component.SummarySection {
	name := e.Name
	value := e.Value
	f := e.ValueFrom
	if value == "" && f != nil {
		cm := f.ConfigMapKeyRef
		if cm != nil {
			value = fmt.Sprintf("from ConfigMap %s %s", cm.Name, cm.Key)
		}
		sec := f.SecretKeyRef
		if sec != nil {
			value = fmt.Sprintf("from Secret %s %s", sec.Name, sec.Key)
		}
	}
	return component.SummarySection{
		Header:  name,
		Content: component.NewMarkdownText(value),
	}
}

func ToVolumeMountSection(pod *corev1.Pod, c *corev1.Container, v *corev1.VolumeMount) component.SummarySection {
	mountPath := v.MountPath
	subPath := v.SubPath
	if subPath != "" {
		mountPath += subPath
	}
	return component.SummarySection{
		Header:  v.Name,
		Content: component.NewMarkdownText(mountPath),
	}
}
