package views // import "github.com/jenkins-x/octant-jx/pkg/plugin/views"

import (
	"fmt"
	"io"
	"strings"

	"github.com/jenkins-x/jx-logging/v3/pkg/log"

	"github.com/jenkins-x/octant-jx/pkg/common/links"
	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx"
	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/jenkins-x/octant-jx/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	corev1 "k8s.io/api/core/v1"
)

const (
	tektonCommandSeparator        = " -c "
	tektonInitialCommandSeparator = "-entrypoint "
)

type containersViewContext struct {
	Request      service.Request
	Namespace    string
	PipelineName string
	PodName      string
}

func (c *containersViewContext) ContainerLink(containerName string) string {
	return plugin.GetPipelineContainerLink(c.Namespace, c.PipelineName, c.PodName, containerName)
}

func BuildPipelineContainersView(request service.Request, pluginContext pluginctx.Context) (component.Component, error) {
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
	breadcrumbs = append(breadcrumbs, viewhelpers.ToMarkdownLink("Pod", links.GetPodLink(ns, name)), "Steps")
	header := component.NewMarkdownText(viewhelpers.ToBreadcrumbMarkdown(breadcrumbs...))
	notesCard := component.NewCard(nil)
	vc := containersViewContext{
		Request:      request,
		Namespace:    ns,
		PipelineName: pipelineName,
		PodName:      name,
	}
	notesCard.SetBody(ToPipelinePodContainersView(vc, pod))

	flexLayout := component.NewFlexLayout("")
	flexLayout.AddSections(component.FlexLayoutSection{
		{Width: component.WidthFull, View: header},
		{Width: component.WidthFull, View: notesCard},
	})
	return flexLayout, nil
}

func ToPipelinePodContainersView(vc containersViewContext, pod *corev1.Pod) component.Component {
	b := &strings.Builder{}
	containers := pod.Spec.Containers
	for k := range containers {
		appendPipelineContainer(vc, b, k, &containers[k])
	}
	return component.NewMarkdownText(b.String())
}

func appendPipelineContainer(vc containersViewContext, w io.StringWriter, index int, c *corev1.Container) {
	name := ToStepName(vc, c)
	n, err := w.WriteString(fmt.Sprintf("* %s\n", name))
	if err != nil {
		log.Logger().Debug(err)
	}
	log.Logger().Debugf("wrote %d bytes\n", n)
	image := ToImage(c)
	commandLine := ToCommandLine(index, c)
	n, err = w.WriteString(fmt.Sprintf("  * **%s**\n", commandLine))
	if err != nil {
		log.Logger().Debug(err)
	}
	log.Logger().Debugf("wrote %d bytes\n", n)
	n, err = w.WriteString(fmt.Sprintf("  * %s\n", image))
	if err != nil {
		log.Logger().Debug(err)
	}
	log.Logger().Debugf("wrote %d bytes\n", n)
}

func ToStepName(vc containersViewContext, c *corev1.Container) string {
	title := strings.TrimPrefix(c.Name, "step-")

	return viewhelpers.ToMarkdownLink(title, vc.ContainerLink(c.Name))
}

func ToImage(c *corev1.Container) string {
	image := c.Image
	link := links.GetImageLink(image)
	if link != "" {
		image = viewhelpers.ToMarkdownLink(image, link)
	} else {
		image = "**" + image + "**"
	}
	return image
}

func ToCommandLine(index int, c *corev1.Container) string {
	args := append(c.Command, c.Args...)
	commandLine := strings.Join(args, " ")

	// lets strip the tekton entry point CLI
	separator := tektonCommandSeparator
	if index == 0 {
		separator = tektonInitialCommandSeparator
	}
	idx := strings.Index(commandLine, separator)
	if idx > 0 {
		commandLine = commandLine[idx+len(separator):]
		if index == 0 {
			commandLine = strings.Replace(commandLine, "-- ", "", 1)
		}
	} else {
		// lets see if its a "... -entrypoint foo -- " string
		separator = " -entrypoint "
		idx = strings.Index(commandLine, separator)
		if idx > 0 {
			commandLine = commandLine[idx+len(separator):]
			commandLine = strings.Replace(commandLine, "-- ", "", 1)
		}
	}
	return commandLine
}
