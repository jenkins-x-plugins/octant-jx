package settings // import "github.com/jenkins-x/octant-jx/pkg/plugin/settings"

import (
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/plugin"
)

func GetCapabilities() *plugin.Capabilities {
	return &plugin.Capabilities{
		ActionNames: []string{action.RequestSetNamespace},
		IsModule:    true,
	}
}
