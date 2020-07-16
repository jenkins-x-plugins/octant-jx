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
	"log"
	"strings"
	"time"

	"github.com/jenkins-x/octant-jx/pkg/common/helm"
	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx"
	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/jenkins-x/octant-jx/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/yaml"
)

// see https://github.com/vmware-tanzu/octant/issues/919
const yamlSupportedInPlugin = false

func BuildHelmReleaseView(request service.Request, pluginContext pluginctx.Context) (component.Component, error) {
	paths := strings.Split(strings.TrimPrefix(request.Path(), "/"), "/")
	releaseName := paths[len(paths)-1]

	ctx := request.Context()
	client := request.DashboardClient()
	ns := pluginContext.Namespace

	//log.Printf("BuildHelmReleasesView querying for secrets owned by helm with release %s\n", releaseName)

	ul, err := client.List(ctx, store.Key{
		APIVersion: "v1",
		Kind:       "Secret",
		Namespace:  ns,
		Selector: &labels.Set{
			"owner":  "helm",
			"status": "deployed",
			"name":   releaseName,
		},
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Printf("BuildHelmReleaseView looking for release %s in namespace %s got list of secrets %d\n", releaseName, ns, len(ul.Items))

	helmReleases := helm.UnstructuredListToAnyHelmReleaseList(ul)
	if len(helmReleases) == 0 {
		return component.NewText(fmt.Sprintf("Error: release %s not found in namespace %s", releaseName, ns)), nil
	}
	r := helmReleases[0]

	header := component.NewMarkdownText(viewhelpers.ToBreadcrumbMarkdown(
		plugin.RootBreadcrumb,
		viewhelpers.ToMarkdownLink("Helm", plugin.GetHelmLink()),
		releaseName))

	statusSummarySections := []component.SummarySection{
		{"Name", ToHelmName(r)},
		{"Status", ToHelmStatus(r)},
		{"Last Deployed", component.NewText(r.Info.LastDeployed.Format(time.ANSIC))},
		{"Revision", component.NewText(fmt.Sprintf("%d", r.Version))},
	}

	statusSummary := component.NewSummary("Status", statusSummarySections...)
	notesCard := viewhelpers.NewMarkdownCard("Notes", fmt.Sprintf("\n%s\n", strings.TrimSpace(r.Info.Notes)))

	flexLayout := component.NewFlexLayout("")
	flexLayout.AddSections(component.FlexLayoutSection{
		{Width: component.WidthFull, View: header},
		{Width: component.WidthHalf, View: statusSummary},
		{Width: component.WidthHalf, View: notesCard},
	})

	if r.Config != nil {
		values := r.Config
		if len(values) > 0 {
			data, err := yaml.Marshal(values)
			if err != nil {
				log.Printf("failed to marshal helm values to YAML: %s with values %#v", err.Error(), values)
			} else {
				var yamlView component.Component

				if yamlSupportedInPlugin {
					yamlView = component.NewYAML(component.TitleFromString("Values"), string(data))
				} else {
					yamlView = viewhelpers.NewMarkdownCard("Values", YAMLToMarkdown(string(data)))
				}
				flexLayout.AddSections(component.FlexLayoutSection{
					{Width: component.WidthFull, View: yamlView},
				})
			}
		}
	}
	return flexLayout, nil
}

// YAMLToMarkdown convert the YAML source into markdown
func YAMLToMarkdown(text string) string {
	listSeparator := "- "

	b := strings.Builder{}
	lines := strings.Split(strings.TrimSpace(text), "\n")
	for _, line := range lines {
		// find first non whitespace character
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		idx := strings.Index(line, trimmed)
		if idx < 0 {
			idx = 0
		}

		// lets remove list indents
		separator := "* "
		remaining := line[idx:]
		if strings.HasPrefix(remaining, listSeparator) {
			remaining = strings.TrimPrefix(line[idx:], listSeparator)
			separator = "  " + separator
		}
		line = line[0:idx] + separator + remaining
		b.WriteString(line)
		b.WriteString("\n")
	}
	return b.String()
}
