package views // import "github.com/jenkins-x/octant-jx/pkg/plugin/views"

import (
	v1 "github.com/jenkins-x/jx-api/pkg/apis/jenkins.io/v1"
	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/jenkins-x/octant-jx/pkg/admin"
	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx"
	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/jenkins-x/octant-jx/pkg/plugin/views"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func BuildFailedReleasePipelinesView(request service.Request, pluginContext pluginctx.Context) (component.Component, error) {
	viewConfig := views.PipelinesViewConfig{
		Context: pluginContext,
		Title:   "Failed Last Release Pipelines",
		Header:  viewhelpers.ToBreadcrumbMarkdown(admin.RootBreadcrumb, "Failed Last Release Pipelines"),
		Filter: func(pa *v1.PipelineActivity, all []*v1.PipelineActivity) bool {
			if pa.Spec.GitBranch == "master" {
				switch pa.Spec.Status {
				case v1.ActivityStatusTypeAborted, v1.ActivityStatusTypeError, v1.ActivityStatusTypeFailed:
					// lets check if there is a newer pipeline that succeeded for this repo/branch
					// note we don't compare context as sometimes its empty
					s := &pa.Spec
					for _, r := range all {
						sr := &r.Spec
						if s.GitOwner == sr.GitOwner && s.GitRepository == sr.GitRepository &&
							s.GitBranch == sr.GitBranch &&
							sr.Status == v1.ActivityStatusTypeSucceeded {

							// lets see if this build is newer
							if viewhelpers.PipelineBuildNumber(r) > viewhelpers.PipelineBuildNumber(pa) {
								log.Logger().Debugf("failed pipeline %s excluded as build %s Succeeded and is newer", pa.Name, r.Spec.Build)
								return false
							}
						}
					}
					return true
				}
			}
			return false
		},
	}
	return views.BuildPipelinesView(request, &viewConfig)
}
