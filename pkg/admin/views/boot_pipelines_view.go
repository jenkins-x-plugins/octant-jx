package views // import "github.com/jenkins-x/octant-jx/pkg/plugin/views"

import (
	"fmt"

	"github.com/jenkins-x/jx-logging/v3/pkg/log"

	v1 "github.com/jenkins-x/jx-api/v4/pkg/apis/jenkins.io/v1"
	"github.com/jenkins-x/jx-helpers/v3/pkg/gitclient/giturl"
	"github.com/jenkins-x/octant-jx/pkg/admin"
	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx"
	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/jenkins-x/octant-jx/pkg/plugin/views"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func BuildBootPipelinesView(request service.Request, pluginContext pluginctx.Context) (component.Component, error) {
	ctx := request.Context()
	client := request.DashboardClient()

	ns := pluginContext.Namespace
	u, err := viewhelpers.GetResourceByName(ctx, client, "jenkins.io/v1", "Environment", "dev", ns)
	if err != nil {
		log.Logger().Infof("failed to find dev Environment in namespace %s", ns)
	}
	gitURL := ""
	owner := ""
	repository := ""
	if u != nil {
		r := &v1.Environment{}
		err = viewhelpers.ToStructured(u, r)
		if err != nil {
			log.Logger().Info(err)
		}
		gitURL = r.Spec.Source.URL
		if gitURL != "" {
			repo, err := giturl.ParseGitURL(gitURL)
			if err != nil {
				return viewhelpers.NewMarkdownText(fmt.Sprintf("could not parse the dev Environment git source URL `%s` due to: %s", gitURL, err.Error())), nil
			}
			if repo != nil {
				owner = repo.Organisation
				repository = repo.Name
			}
		}
	}
	if owner == "" || repository == "" {
		log.Logger().Infof("No _dev_ **Environment** resource is available with a link to the git repository in the %s namespace for git URL %s", ns, gitURL)
	}

	header := viewhelpers.ToBreadcrumbMarkdown(admin.RootBreadcrumb, "Boot Pipelines")
	if gitURL != "" {
		header = viewhelpers.ToBreadcrumbMarkdown(admin.RootBreadcrumb, "Boot Pipelines", viewhelpers.ToMarkdownLink("Source", gitURL))
	}
	viewConfig := views.PipelinesViewConfig{
		Context: pluginContext,
		Title:   "Boot Pipelines",
		Header:  header,
		Columns: []string{"Branch", "Build", "Status", "Message"},
		Filter: func(pa *v1.PipelineActivity, _ []v1.PipelineActivity) bool {
			return pa.Spec.GitOwner == owner && pa.Spec.GitRepository == repository
		},
	}
	return views.BuildPipelinesView(request, &viewConfig)
}
