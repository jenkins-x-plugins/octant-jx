package views // import "github.com/jenkins-x/octant-jx/pkg/plugin/views"

import (
	"fmt"
	"strings"

	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx"
	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/jenkins-x/octant-jx/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func BuildOverview(request service.Request, pluginContext pluginctx.Context) (component.Component, error) {
	b := strings.Builder{}
	b.WriteString(`<h1>Jenkins X</h1>
`)

	for _, nav := range plugin.Navigations {
		link := plugin.PathPrefix + "/" + nav.Path
		b.WriteString(fmt.Sprintf(`
<a href="%s"><h3>%s</h3></a>
`, link, nav.Title))
	}

	return viewhelpers.NewMarkdownText(b.String()), nil
}
