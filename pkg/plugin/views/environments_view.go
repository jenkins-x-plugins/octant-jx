package views // import "github.com/jenkins-x/octant-jx/pkg/plugin/views"

import (
	"log"

	v1 "github.com/jenkins-x/jx-api/pkg/apis/jenkins.io/v1"
	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx"
	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/jenkins-x/octant-jx/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func BuildEnvironmentsView(request service.Request, pluginContext pluginctx.Context) (component.Component, error) {
	ctx := request.Context()
	client := request.DashboardClient()

	dl, err := client.List(ctx, store.Key{
		APIVersion: "jenkins.io/v1",
		Kind:       "Environment",
		Namespace:  pluginContext.Namespace,
	})

	if err != nil {
		log.Printf("failed: %s", err.Error())
		return nil, err
	}

	log.Printf("got list of Environment %d\n", len(dl.Items))

	header := component.NewMarkdownText(viewhelpers.ToBreadcrumbMarkdown(plugin.RootBreadcrumb, "Environments"))

	table := component.NewTableWithRows(
		"Environments", "There are no Environments!",
		component.NewTableCols("Name", "Namespace", "Promote", "Source"),
		[]component.TableRow{})

	for _, pa := range dl.Items {
		tr, err := toEnvironmentTableRow(pa)
		if err != nil {
			log.Printf("failed to create Table Row: %s", err.Error())
			continue
		}
		if tr != nil {
			table.Add(*tr)
		}
	}

	table.Sort("Name", false)
	flexLayout := component.NewFlexLayout("")
	flexLayout.AddSections(component.FlexLayoutSection{
		{Width: component.WidthFull, View: header},
		{Width: component.WidthFull, View: table},
	})

	return flexLayout, nil
}

func toEnvironmentTableRow(u unstructured.Unstructured) (*component.TableRow, error) {
	r := &v1.Environment{}
	err := viewhelpers.ToStructured(&u, r)
	if err != nil {
		return nil, err
	}

	name := r.Name
	if name == "" {
		name = u.GetName()
	}
	return &component.TableRow{
		"Name":      ToEnvironmentNameLink(r),
		"Source":    ToEnvironmentSource(r),
		"Namespace": ToEnvironmentNamespace(r),
		"Promote":   ToEnvironmentPromote(r),
	}, nil
}
