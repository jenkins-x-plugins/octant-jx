package views // import "github.com/jenkins-x/octant-jx/pkg/plugin/views"

import (
	"fmt"

	"github.com/jenkins-x/jx-logging/v3/pkg/log"

	"github.com/jenkins-x/jx-preview/pkg/apis/preview/v1alpha1"
	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx"
	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type PreviewsViewConfig struct {
	OwnerFilter      component.TableFilter
	RepositoryFilter component.TableFilter
}

func (f *PreviewsViewConfig) TableFilters() []*component.TableFilter {
	return []*component.TableFilter{&f.OwnerFilter, &f.RepositoryFilter}
}

func BuildPreviewsView(request service.Request, pluginContext pluginctx.Context) (component.Component, error) {
	ctx := request.Context()
	client := request.DashboardClient()

	dl, err := client.List(ctx, store.Key{
		APIVersion: "preview.jenkins.io/v1alpha1",
		Kind:       "Preview",
		Namespace:  pluginContext.Namespace,
	})

	if err != nil {
		log.Logger().Infof("failed: %s", err.Error())
		return nil, err
	}

	log.Logger().Infof("got list of Preview %d\n", len(dl.Items))

	config := &PreviewsViewConfig{}

	table := component.NewTableWithRows(
		viewhelpers.TableTitle("Previews"), "There are no Previews!",
		component.NewTableCols("Owner", "Repository", "Pull Request", "Preview", "Author"),
		[]component.TableRow{})

	for _, pa := range dl.Items {
		tr, err := toPreviewTableRow(pa, config)
		if err != nil {
			log.Logger().Infof("failed to create Table Row: %s", err.Error())
			continue
		}
		if tr != nil {
			table.Add(*tr)
		}
	}

	table.Sort("Sort", false)

	viewhelpers.InitTableFilters(config.TableFilters())

	table.AddFilter("Owner", config.OwnerFilter)
	table.AddFilter("Repository", config.RepositoryFilter)

	flexLayout := component.NewFlexLayout("")
	flexLayout.AddSections(component.FlexLayoutSection{
		{Width: component.WidthFull, View: table},
	})

	return flexLayout, nil
}

func toPreviewTableRow(u unstructured.Unstructured, filters *PreviewsViewConfig) (*component.TableRow, error) {
	r := &v1alpha1.Preview{}
	err := viewhelpers.ToStructured(&u, r)
	if err != nil {
		return nil, err
	}

	prs := &r.Spec.PullRequest
	prLink := prs.URL

	prName := fmt.Sprintf("#%d", prs.Number)
	previewLink := ""
	appURL := r.Spec.Resources.URL
	if appURL != "" {
		previewLink = fmt.Sprintf("<a href='%s' title='try the application' target='preview' class='badge badge-purple'>Try Me&nbsp;<clr-icon shape='link'></clr-icon></a>", appURL)
	}

	authorLink := ""
	username := prs.User.Username
	if username != "" {
		name := prs.User.Name
		if name == "" {
			name = username
		}
		authorLink = fmt.Sprintf("<a href='%s' title='%s' target='author'>%s</a>", prs.User.LinkURL, name, username)
		if prs.User.ImageURL != "" {
			authorLink = fmt.Sprintf("<img height='24' width='24' src='%s'> %s", prs.User.ImageURL, authorLink)
		}
	}

	owner := prs.Owner
	repository := prs.Repository
	fullName := fmt.Sprintf("%s/%s/%010d", owner, repository, prs.Number)

	viewhelpers.AddFilterValue(&filters.OwnerFilter, owner)
	viewhelpers.AddFilterValue(&filters.RepositoryFilter, repository)

	return &component.TableRow{
		"Sort":         component.NewText(fullName),
		"Owner":        component.NewText(owner),
		"Repository":   component.NewText(repository),
		"Pull Request": viewhelpers.NewMarkdownText(viewhelpers.ToMarkdownLink(prName, prLink) + " " + prs.Title),
		"Preview":      viewhelpers.NewMarkdownText(previewLink),
		"Author":       viewhelpers.NewMarkdownText(authorLink),
	}, nil
}
