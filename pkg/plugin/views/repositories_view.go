package views // import "github.com/jenkins-x/octant-jx/pkg/plugin/views"

import (
	"html"
	"strings"

	"github.com/jenkins-x/jx-logging/v3/pkg/log"

	v1 "github.com/jenkins-x/jx-api/v3/pkg/apis/jenkins.io/v1"
	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx"
	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/jenkins-x/octant-jx/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type RepositoriesViewConfig struct {
	OwnerFilter component.TableFilter
}

func (f *RepositoriesViewConfig) TableFilters() []*component.TableFilter {
	return []*component.TableFilter{&f.OwnerFilter}
}

func BuildRepositoriesView(request service.Request, pluginContext pluginctx.Context) (component.Component, error) {
	ctx := request.Context()
	client := request.DashboardClient()

	dl, err := client.List(ctx, store.Key{
		APIVersion: "jenkins.io/v1",
		Kind:       "SourceRepository",
		Namespace:  pluginContext.Namespace,
	})

	if err != nil {
		log.Logger().Infof("failed: %s", err.Error())
		return nil, err
	}

	log.Logger().Infof("got list of SourceRepository %d\n", len(dl.Items))

	header := viewhelpers.NewMarkdownText(viewhelpers.ToBreadcrumbMarkdown(plugin.RootBreadcrumb, "Repositories"))

	config := &RepositoriesViewConfig{}

	table := component.NewTableWithRows(
		"Repositories", "There are no Repositories!",
		component.NewTableCols("Owner", "Name", "Status"),
		[]component.TableRow{})

	for _, pa := range dl.Items {
		tr, err := toRepositoryTableRow(pa, config)
		if err != nil {
			log.Logger().Infof("failed to create Table Row: %s", err.Error())
			continue
		}
		if tr != nil {
			table.Add(*tr)
		}
	}

	table.Sort("Name", false)

	viewhelpers.InitTableFilters(config.TableFilters())

	table.AddFilter("Owner", config.OwnerFilter)

	flexLayout := component.NewFlexLayout("")
	flexLayout.AddSections(component.FlexLayoutSection{
		{Width: component.WidthFull, View: header},
		{Width: component.WidthFull, View: table},
	})

	return flexLayout, nil
}

func toRepositoryTableRow(u unstructured.Unstructured, config *RepositoriesViewConfig) (*component.TableRow, error) {
	r := &v1.SourceRepository{}
	err := viewhelpers.ToStructured(&u, r)
	if err != nil {
		return nil, err
	}

	owner := r.Spec.Org
	viewhelpers.AddFilterValue(&config.OwnerFilter, owner)

	return &component.TableRow{
		"Owner":  component.NewText(owner),
		"Name":   ToRepositoryName(r),
		"Status": ToRepositoryStatus(r),
	}, nil
}

func ToRepositoryStatus(r *v1.SourceRepository) component.Component {
	status := ""
	if r.Annotations != nil {
		value := strings.ToLower(r.Annotations["webhook.jenkins-x.io"])
		if value == "true" {
			status = `<clr-icon shape="check-circle" class="is-solid is-success" title="Webhook registered successfully"></clr-icon>`
		}
		if value != "" {
			if strings.HasPrefix(value, "creat") {
				status = `<span class="spinner spinner-inline" title="Registering webhook..."></span>`
			} else {
				text := "Failed to register Webook"
				message := r.Annotations["webhook.jenkins-x.io/error"]
				if message != "" {
					text += ": " + html.EscapeString(message)
				}
				status = `<clr-icon shape="warning-standard" class="is-solid is-danger" title="` + text + `"></clr-icon>`
			}
		}
	}
	return viewhelpers.NewMarkdownText(status)
}

func ToRepositoryName(r *v1.SourceRepository) component.Component {
	s := &r.Spec
	u := s.URL
	if u == "" {
		u = s.HTTPCloneURL
	}
	if u == "" {
		if s.Org != "" && s.Repo != "" {
			u = s.Org + "/" + s.Repo
			if s.Provider != "" {
				u = s.Provider + "/" + u
			}
		}
	}
	return viewhelpers.NewMarkdownText(viewhelpers.ToMarkdownLink(s.Repo, u))
}
