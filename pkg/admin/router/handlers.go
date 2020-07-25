package router // import "github.com/jenkins-x/octant-jx/pkg/plugin/router"

import (
	"strings"

	"github.com/jenkins-x/jx-logging/pkg/log"

	"github.com/jenkins-x/jx-helpers/pkg/gitclient/giturl"
	"github.com/jenkins-x/octant-jx/pkg/admin"
	"github.com/jenkins-x/octant-jx/pkg/admin/workspaces"
	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx"
	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/pkg/errors"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/jenkins-x/octant-jx/pkg/admin/views"
)

type Handlers struct {
	Context *pluginctx.Context

	Workspaces []*workspaces.Workspace
	Octants    *workspaces.Octants
}

func (h *Handlers) Load() error {
	var err error
	h.Workspaces, err = workspaces.LoadWorkspaces()
	if err != nil {
		return errors.Wrap(err, "failed to load workspaces")
	}

	h.Octants, err = workspaces.NewOctants()
	if err != nil {
		return errors.Wrap(err, "failed to load octants")
	}
	return nil
}

func (h *Handlers) GetWorkspaces() []workspaces.WorkspaceOctant {
	return workspaces.ToWorkspaceOctants(h.Workspaces, h.Octants)
}

func (h *Handlers) InitRoutes(router *service.Router) {
	h.handleView(router, admin.OverviewPath, views.BuildOverview)
	h.handleView(router, admin.BootPipelinesPath, views.BuildBootPipelinesView)
	h.handleView(router, admin.FailedPipelinesPath, views.BuildFailedReleasePipelinesView)
	h.handleView(router, admin.HealthPath, views.HealthView)

	h.handleJobsPath(router, admin.GCPipelineJobsPath)
	h.handleJobsPath(router, admin.GCPodJobsPath)
	h.handleJobsPath(router, admin.GCPreviewJobsPath)
	h.handleJobsPath(router, admin.UpgradeJobsPath)

	router.HandleFunc("/"+admin.BootJobsPath, h.handleBootJobsView)
	h.addBootJobLogsHandlers(router)

	router.HandleFunc("/"+admin.WorkspacesPath, h.handleWorkspacesPath)
}

func (h *Handlers) handleView(router *service.Router, path string, fn func(request service.Request, pluginContext pluginctx.Context) (component.Component, error)) {
	router.HandleFunc("/"+path, func(request service.Request) (component.ContentResponse, error) {
		view, err := fn(request, *h.Context)
		if err != nil {
			return component.EmptyContentResponse, err
		}
		response := component.NewContentResponse(nil)
		response.Add(view)
		return *response, nil
	})
}

func (h *Handlers) handleJobsPathEnrich(router *service.Router, path string, enrichFn func(request service.Request, pluginContext pluginctx.Context, cr *component.ContentResponse, view component.Component) error) {
	router.HandleFunc("/"+path, func(request service.Request) (component.ContentResponse, error) {
		view, err := views.BuildJobsViewForPath(request, *h.Context, path)
		if err != nil {
			return component.EmptyContentResponse, err
		}
		response := component.NewContentResponse(nil)
		response.Add(view)
		if enrichFn != nil {
			err := enrichFn(request, *h.Context, response, view)
			if err != nil {
				log.Logger().Debug(err)
			}
		}
		return *response, nil
	})
	h.addJobLogsHandlers(router, path)
}

func (h *Handlers) addBootJobLogsHandlers(router *service.Router) {
	path := admin.BootJobsPath
	router.HandleFunc("/"+path+"/logs/*", func(request service.Request) (component.ContentResponse, error) {
		config := views.JobsViewConfigs[path]
		view, err := views.BuildJobsViewLogsForPathAndSelector(request, *h.Context, path, "jx-boot", config, labels.Set{})
		if err != nil {
			return component.EmptyContentResponse, err
		}
		response := component.NewContentResponse(nil)
		response.Add(view)
		return *response, nil
	})
}

func (h *Handlers) addJobLogsHandlers(router *service.Router, path string) {
	router.HandleFunc("/"+path+"/logs", func(request service.Request) (component.ContentResponse, error) {
		view, err := views.BuildJobsLogViewForPath(request, *h.Context, path, "")
		if err != nil {
			return component.EmptyContentResponse, err
		}
		response := component.NewContentResponse(nil)
		response.Add(view)
		return *response, nil
	})

	router.HandleFunc("/"+path+"/logs/*", func(request service.Request) (component.ContentResponse, error) {
		paths := strings.Split(strings.TrimSuffix(request.Path(), "/"), "/")
		name := paths[len(paths)-1]

		view, err := views.BuildJobsLogViewForPath(request, *h.Context, path, name)
		if err != nil {
			return component.EmptyContentResponse, err
		}
		response := component.NewContentResponse(nil)
		response.Add(view)
		return *response, nil
	})
}

func (h *Handlers) handleJobsPath(router *service.Router, path string) {
	h.handleJobsPathEnrich(router, path, nil)
}

func (h *Handlers) handleBootJobsView(request service.Request) (component.ContentResponse, error) {
	response := component.NewContentResponse(nil)
	pluginContext := *h.Context
	ctx := request.Context()
	client := request.DashboardClient()
	ns := pluginContext.Namespace

	name := "jx-boot-git-url"
	u, err := viewhelpers.GetResourceByName(ctx, client, "v1", "Secret", name, ns)
	if err != nil || u == nil {
		if err != nil {
			log.Logger().Infof("failed to load Secret %s in namespace %s: %s", name, ns, err.Error())
		}
		err = views.BuildNoBootSecretView(request, pluginContext, response)
		if err != nil {
			log.Logger().Info(err)
			return component.EmptyContentResponse, err
		}
		return *response, nil
	}

	// validate the secrets have been validated...
	validSecrets := false
	gitURL := ""
	secret := &corev1.Secret{}
	err = viewhelpers.ToStructured(u, secret)
	if err != nil {
		log.Logger().Info(err)
	} else if secret.Data != nil {
		d := secret.Data["secrets-verify"]
		if d != nil {
			text := strings.ToLower(string(d))
			if text == "true" || text == "valid" {
				validSecrets = true
			}
		}

		d = secret.Data["git-url"]
		if d != nil {
			text := string(d)
			if text != "" {
				repo, _ := giturl.ParseGitURL(text)
				if repo != nil && repo.Name != "" {
					gitURL = repo.URLWithoutUser()
				}
			}
		}
	}
	if !validSecrets {
		err = views.BuildBootInvalidSecretView(request, pluginContext, response, gitURL)
		if err != nil {
			log.Logger().Info(err)
			return component.EmptyContentResponse, err
		}
		return *response, nil
	}

	view, err := views.BuildJobsViewForPath(request, pluginContext, admin.BootJobsPath)
	if err != nil {
		return component.EmptyContentResponse, err
	}
	response.Add(view)
	err = views.BootJobExtraView(request, pluginContext, response, view)
	if err != nil {
		return component.EmptyContentResponse, err
	}
	return *response, nil
}

func (h *Handlers) handleWorkspacesPath(request service.Request) (component.ContentResponse, error) {
	ws := h.GetWorkspaces()
	view, err := views.BuildWorkspacesView(request, ws)
	if err != nil {
		return component.EmptyContentResponse, err
	}
	response := component.NewContentResponse(nil)
	response.Add(view)
	return *response, nil
}
