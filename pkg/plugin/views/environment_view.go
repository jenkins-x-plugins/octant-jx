package views // import "github.com/jenkins-x/octant-jx/pkg/plugin/views"

import (
	"fmt"
	"log"
	"strings"

	v1 "github.com/jenkins-x/jx-api/pkg/apis/jenkins.io/v1"
	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx"
	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/jenkins-x/octant-jx/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func BuildEnvironmentView(request service.Request, pluginContext pluginctx.Context) (component.Component, error) {
	name := strings.TrimPrefix(request.Path(), "/")
	name = strings.TrimPrefix(name, plugin.EnvironmentsPath+"/")
	ctx := request.Context()
	client := request.DashboardClient()

	u, err := viewhelpers.GetResourceByName(ctx, client, "jenkins.io/v1", "Environment", name, pluginContext.Namespace)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return component.NewText("Error: environment not found"), nil
	}

	r := &v1.Environment{}
	err = viewhelpers.ToStructured(u, r)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	header := component.NewMarkdownText(viewhelpers.ToBreadcrumbMarkdown(plugin.RootBreadcrumb, viewhelpers.ToMarkdownLink("Environments", plugin.GetEnviromentsLink()), ToEnvironmentName(r)))

	summary := component.NewSummary("Summary",
		component.SummarySection{"Name", ToEnvironmentNameComponent(r)},
		component.SummarySection{"Source", ToEnvironmentSource(r)},
		component.SummarySection{"Namespace", ToEnvironmentNamespace(r)},
		component.SummarySection{"Promote", ToEnvironmentPromote(r)},
	)

	flexLayout := component.NewFlexLayout("")
	flexLayout.AddSections(component.FlexLayoutSection{
		{Width: component.WidthFull, View: header},
		{Width: component.WidthHalf, View: summary},
	})
	return flexLayout, nil
}

func BuildEnvironmentAppsView(request service.Request, pluginContext pluginctx.Context) (component.Component, error) {
	name := strings.TrimPrefix(request.Path(), "/")
	name = strings.TrimPrefix(name, plugin.EnvironmentsPath+"/")
	ctx := request.Context()
	client := request.DashboardClient()

	u, err := viewhelpers.GetResourceByName(ctx, client, "jenkins.io/v1", "Environment", name, pluginContext.Namespace)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return component.NewText("Error: pipeline not found"), nil
	}

	pa := &v1.Environment{}
	err = viewhelpers.ToStructured(u, pa)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	header := component.NewMarkdownText(fmt.Sprintf("## [Environments](%s) / %s", plugin.GetEnviromentsLink(), ToEnvironmentName(pa)))

	summary := component.NewSummary("Apps",
		component.SummarySection{"Name", ToEnvironmentNameComponent(pa)},
		component.SummarySection{"Source", ToEnvironmentSource(pa)},
		component.SummarySection{"Namespace", ToEnvironmentNamespace(pa)},
		component.SummarySection{"Promote", ToEnvironmentPromote(pa)},
	)

	flexLayout := component.NewFlexLayout("")
	flexLayout.AddSections(component.FlexLayoutSection{
		{Width: component.WidthFull, View: header},
		{Width: component.WidthHalf, View: summary},
	})
	return flexLayout, nil
}
