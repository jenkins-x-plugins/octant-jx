/*
Copyright 2019 Blood Orange

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package views

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jenkins-x/jx-logging/pkg/log"

	"github.com/jenkins-x/octant-jx/pkg/common/helm"
	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx"
	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/jenkins-x/octant-jx/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	rspb "helm.sh/helm/v3/pkg/release"
	"k8s.io/apimachinery/pkg/labels"
)

func BuildHelmReleasesView(request service.Request, pluginContext pluginctx.Context) (component.Component, error) {
	ctx := request.Context()
	client := request.DashboardClient()

	ul, err := client.List(ctx, store.Key{
		APIVersion: "v1",
		Kind:       "Secret",
		Namespace:  pluginContext.Namespace,
		Selector: &labels.Set{
			"owner":  "helm",
			"status": "deployed",
		},
	})

	if err != nil {
		log.Logger().Infof("failed: %s", err.Error())
		return nil, err
	}

	log.Logger().Debugf("BuildHelmReleasesView got list of secrets %d\n", len(ul.Items))

	helmReleases := helm.UnstructuredListToHelmReleaseList(ul)

	header := component.NewMarkdownText(viewhelpers.ToBreadcrumbMarkdown(plugin.RootBreadcrumb, "Helm"))

	table := component.NewTableWithRows(
		"Releases", "There are no Helm releases!",
		component.NewTableCols("Name", "Status", "Updated", "Revision", "Chart", "App Version"),
		[]component.TableRow{})

	for _, r := range helmReleases {
		tr := component.TableRow{
			"Name":     ToHelmName(r),
			"Revision": component.NewText(strconv.Itoa(r.Version)),
			"Status":   ToHelmStatus(r),
			"Updated":  ToHelmUpdated(r),
			"Chart": component.NewText(
				fmt.Sprintf("%s-%s", r.Chart.Metadata.Name, r.Chart.Metadata.Version)),
			"App Version": component.NewText(r.Chart.Metadata.AppVersion),
			"Sort":        component.NewText(r.Name),
		}
		table.Add(tr)
	}

	table.Sort("Sort", false)

	flexLayout := component.NewFlexLayout("")
	flexLayout.AddSections(component.FlexLayoutSection{
		{Width: component.WidthFull, View: header},
		{Width: component.WidthFull, View: table},
	})

	return flexLayout, nil
}

func ToHelmStatus(r *rspb.Release) component.Component {
	statusText := strings.Title(r.Info.Status.String())
	icon := ToHelmStatusIcon(r.Info.Status)
	if icon != "" {
		return component.NewMarkdownText(fmt.Sprintf(`%s %s`, icon, statusText))
	}
	return component.NewText(statusText)
}

func ToHelmStatusIcon(s rspb.Status) string {
	switch s {
	case rspb.StatusDeployed:
		return `<clr-icon shape="check-circle" class="is-solid is-success" title="Deployed"></clr-icon>`

	case rspb.StatusFailed:
		return `<clr-icon shape="warning-standard" class="is-solid is-danger" title="Failed"></clr-icon>`

	case rspb.StatusPendingInstall, rspb.StatusPendingUpgrade, rspb.StatusPendingRollback:
		return `<clr-icon shape="clock" title="Pending"></clr-icon>`

	case rspb.StatusUninstalling:
		return `<span class="spinner spinner-inline" title="Uninstalling"></span>`

	case rspb.StatusSuperseded, rspb.StatusUninstalled:
		return `<clr-icon shape="trash"></clr-icon>`

	case rspb.StatusUnknown:
		return `<clr-icon shape="unknown-status"></clr-icon>`

	default:
		return `<clr-icon shape="unknown-status"></clr-icon>`
	}
}

func ToHelmName(r *rspb.Release) component.Component {
	icon := ""
	if r.Chart != nil && r.Chart.Metadata != nil {
		icon = viewhelpers.ToApplicationIcon(r.Chart.Metadata.Icon)
	}
	if icon == "" {
		icon = DefaultIcon
	}
	iconPrefix := fmt.Sprintf(`<img src="%s" width="24" height="24">&nbsp;`, icon)
	name := r.Name
	ref := plugin.GetHelmReleaseLink(r.Name)
	return component.NewMarkdownText(fmt.Sprintf(`%s<a href="%s" title="Helm Release %s">%s</a>`, iconPrefix, ref, name, name))
}

func ToHelmUpdated(r *rspb.Release) component.Component {
	if t := r.Info.LastDeployed; !t.IsZero() {
		return component.NewTimestamp(t.Time)
	}
	return component.NewText("")

}
