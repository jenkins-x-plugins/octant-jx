package views // import "github.com/jenkins-x/octant-jx/pkg/plugin/views"

import (
	"github.com/jenkins-x/jx-logging/v3/pkg/log"

	"github.com/jenkins-x/octant-jx/pkg/admin"
	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx"
	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/jenkins-x/octant-jx/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

const (
	overviewUseCards = false
)

// ViewInfo represents the information for views
type ViewInfo struct {
	Title   string
	Path    string
	Factory func(request service.Request, pluginContext pluginctx.Context) (component.Component, error)
	Width   int
}

func BuildOverview(request service.Request, pluginContext pluginctx.Context) (component.Component, error) {
	views := []ViewInfo{
		{
			Title: "Boot Jobs",
			Path:  admin.BootJobsPath,
		},
		{
			Title: "Upgrade Jobs",
			Path:  admin.UpgradeJobsPath,
		},
		{
			Title:   "Boot Pipelines",
			Path:    admin.BootPipelinesPath,
			Factory: BuildBootPipelinesView,
			Width:   component.WidthFull,
		},
		{
			Title: "GC Pipeline Jobs",
			Path:  admin.GCPipelineJobsPath,
		},
		{
			Title: "GC Pod Jobs",
			Path:  admin.GCPodJobsPath,
		},
		{
			Title: "GC Preview Jobs",
			Path:  admin.GCPreviewJobsPath,
		},
	}
	layout := component.NewFlexLayout("Overview")
	section := component.FlexLayoutSection{}

	header := component.NewMarkdownText(viewhelpers.ToBreadcrumbMarkdown(admin.RootBreadcrumb, "Overview"))
	section = append(section, component.FlexLayoutItem{
		Width: component.WidthFull,
		View:  header,
	})

	overviewContext := pluginContext
	overviewContext.Composite = true
	for _, v := range views {
		if v.Factory == nil {
			v.Factory = func(request service.Request, pluginContext pluginctx.Context) (component.Component, error) {
				return BuildJobsViewForPath(request, overviewContext, v.Path)
			}
		}
		view, err := v.Factory(request, overviewContext)
		if err != nil {
			log.Logger().Infof("failed to view %s %s", v.Title, err.Error())
			return layout, err
		}

		if overviewUseCards {
			view = createCard(v.Title, v.Path, view)
		}
		width := v.Width
		if width == 0 {
			width = component.WidthHalf
		}
		section = append(section, component.FlexLayoutItem{
			Width: width,
			View:  view,
		})
	}
	layout.AddSections(section)
	return layout, nil
}

func createCard(title, path string, view component.Component) component.Component {
	section := component.FlexLayoutSection{
		{
			Width: component.WidthFull,
			View:  view,
		},
	}
	link := plugin.PathPrefix + admin.PluginName + "/" + path
	cardTitle := component.NewMarkdownText(viewhelpers.ToMarkdownLink(title, link))
	card := component.NewCard([]component.TitleComponent{cardTitle})

	layout := component.NewFlexLayout(title)
	layout.AddSections(section)
	card.SetBody(layout)
	return card
}
