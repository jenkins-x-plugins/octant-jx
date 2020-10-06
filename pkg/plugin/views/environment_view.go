package views // import "github.com/jenkins-x/octant-jx/pkg/plugin/views"

import (
	"fmt"
	"strings"

	"github.com/jenkins-x/jx-logging/v3/pkg/log"

	v1 "github.com/jenkins-x/jx-api/v3/pkg/apis/jenkins.io/v1"
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
		log.Logger().Info(err)
		return nil, err
	}

	header := component.NewMarkdownText(viewhelpers.ToBreadcrumbMarkdown(plugin.RootBreadcrumb, viewhelpers.ToMarkdownLink("Environments", plugin.GetEnviromentsLink()), ToEnvironmentName(r)))

	summary := component.NewSummary("Summary",
		component.SummarySection{Header: "Name", Content: ToEnvironmentNameComponent(r)},
		component.SummarySection{Header: "Source", Content: ToEnvironmentSource(r)},
		component.SummarySection{Header: "Namespace", Content: ToEnvironmentNamespace(r)},
		component.SummarySection{Header: "Promote", Content: ToEnvironmentPromote(r)},
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
		log.Logger().Info(err)
		return nil, err
	}

	header := component.NewMarkdownText(fmt.Sprintf("## [Environments](%s) / %s", plugin.GetEnviromentsLink(), ToEnvironmentName(pa)))

	summary := component.NewSummary("Apps",
		component.SummarySection{Header: "Name", Content: ToEnvironmentNameComponent(pa)},
		component.SummarySection{Header: "Source", Content: ToEnvironmentSource(pa)},
		component.SummarySection{Header: "Namespace", Content: ToEnvironmentNamespace(pa)},
		component.SummarySection{Header: "Promote", Content: ToEnvironmentPromote(pa)},
	)

	flexLayout := component.NewFlexLayout("")
	flexLayout.AddSections(component.FlexLayoutSection{
		{Width: component.WidthFull, View: header},
		{Width: component.WidthHalf, View: summary},
	})
	return flexLayout, nil
}
