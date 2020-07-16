package settings // import "github.com/jenkins-x/octant-jx/pkg/plugin/settings"

import (
	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx"
	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx/actioners"
	"github.com/jenkins-x/octant-jx/pkg/plugin"
	"github.com/jenkins-x/octant-jx/pkg/plugin/router"
	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
)

func GetOptions(pluginContext *pluginctx.Context) []service.PluginOption {
	h := router.Handlers{
		Context: pluginContext,
	}
	return []service.PluginOption{
		service.WithNavigation(
			func(_ *service.NavigationRequest) (navigation.Navigation, error) {
				return navigation.Navigation{
					Title:    "Jenkins X",
					Path:     plugin.Name + "/" + plugin.PipelinesPath,
					IconName: rootNavIcon,
					Children: plugin.Navigations,
				}, nil
			},
			h.InitRoutes,
		),
		service.WithActionHandler(actioners.CreateHandler(pluginContext)),
	}
}
