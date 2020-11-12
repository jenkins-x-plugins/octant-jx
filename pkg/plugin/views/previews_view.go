package views // import "github.com/jenkins-x/octant-jx/pkg/plugin/views"

import (
	"fmt"
	"html"
	"strings"

	"github.com/jenkins-x/jx-logging/v3/pkg/log"

	"github.com/jenkins-x/jx-preview/pkg/apis/preview/v1alpha1"
	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx"
	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/jenkins-x/octant-jx/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type PreviewsViewConfig struct {
	OwnerFilter component.TableFilter
}

func (f *PreviewsViewConfig) TableFilters() []*component.TableFilter {
	return []*component.TableFilter{&f.OwnerFilter}
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

	header := viewhelpers.NewMarkdownText(viewhelpers.ToBreadcrumbMarkdown(plugin.RootBreadcrumb, "Previews"))

	config := &PreviewsViewConfig{}

	table := component.NewTableWithRows(
		"Previews", "There are no Previews!",
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

func toPreviewTableRow(u unstructured.Unstructured, config *PreviewsViewConfig) (*component.TableRow, error) {
	r := &v1alpha1.Preview{}
	err := viewhelpers.ToStructured(&u, r)
	if err != nil {
		return nil, err
	}

	prs := &r.Spec.PullRequest
	ownerLink, repoLink, prLink := ToOwnerRepoLinks(r)

	prName := fmt.Sprintf("#%d", prs.Number)
	previewLink := ""
	appURL := r.Spec.Resources.URL
	if appURL != "" {
		previewLink = fmt.Sprintf("<a href='%s' title='try the application' target='preview' class='badge badge-purple'>Preview</a>", appURL)
	}

	authorLink := ""
	username := prs.User.Username
	if username != "" {
		name := prs.User.Name
		if name == "" {
			name = username
		}
		authorLink = fmt.Sprintf("<a href='%s' title='%s' target='author>%s</a>", prs.User.LinkURL, name, username)
		if prs.User.ImageURL != "" {
			authorLink = fmt.Sprintf("<img src='%s'> %s", prs.User.ImageURL, authorLink)
		}
	}

	return &component.TableRow{
		"Owner":        viewhelpers.NewMarkdownText(viewhelpers.ToMarkdownLink(prs.Owner, ownerLink)),
		"Repository":   viewhelpers.NewMarkdownText(viewhelpers.ToMarkdownLink(prs.Repository, repoLink)),
		"Pull Request": viewhelpers.NewMarkdownText(viewhelpers.ToMarkdownLink(prName, prLink) + " " + prs.Title),
		"Preview":      viewhelpers.NewMarkdownText(previewLink),
		"Author":       viewhelpers.NewMarkdownText(authorLink),
	}, nil
}

func ToPreviewStatus(r *v1alpha1.Preview) component.Component {
	status := ""
	if r.Annotations != nil {
		value := strings.ToLower(r.Annotations["webhook.jenkins-x.io"])
		if value == "true" {
			status = `<clr-icon shape="check-circle" class="is-solid is-success" title="Webhook registered successfully"></clr-icon>`
		} else if value != "" {
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

func ToOwnerRepoLinks(r *v1alpha1.Preview) (ownerLink string, repoLink string, prLink string) {
	s := &r.Spec.PullRequest
	owner := s.Owner
	prLink = s.URL
	idx := strings.Index(prLink, owner)
	if idx > 0 {
		ownerLink = prLink[0:idx] + owner + "/"
		repoLink = ownerLink + "/" + s.Repository + "/"
	}
	return
}
