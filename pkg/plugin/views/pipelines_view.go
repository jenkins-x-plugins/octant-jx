package views // import "github.com/jenkins-x/octant-jx/pkg/plugin/views"

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/jenkins-x/jx-logging/pkg/log"

	v1 "github.com/jenkins-x/jx-api/pkg/apis/jenkins.io/v1"
	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx"
	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/jenkins-x/octant-jx/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

type PipelinesViewConfig struct {
	Context          pluginctx.Context
	Columns          []string
	Title            string
	Header           string
	Filter           func(*v1.PipelineActivity, []*v1.PipelineActivity) bool
	OwnerFilter      component.TableFilter
	RepositoryFilter component.TableFilter
	BranchFilter     component.TableFilter
	StatusFilter     component.TableFilter
}

func BuildPipelinesViewDefault(request service.Request, pluginContext pluginctx.Context) (component.Component, error) {
	return BuildPipelinesView(request, &PipelinesViewConfig{Context: pluginContext})
}

func BuildPipelinesViewRecent(request service.Request, pluginContext pluginctx.Context) (component.Component, error) {
	recentTime := time.Now().Add(-1 * time.Hour)

	config := &PipelinesViewConfig{Context: pluginContext}
	config.Title = "Recent Pipelines"
	config.Filter = func(pa *v1.PipelineActivity, all []*v1.PipelineActivity) bool {
		completed := pa.Spec.CompletedTimestamp
		if completed != nil {
			if completed.Time.Before(recentTime) {
				return false
			}
		}
		return true
	}
	return BuildPipelinesView(request, config)
}

func BuildPipelinesView(request service.Request, config *PipelinesViewConfig) (component.Component, error) {
	ctx := request.Context()
	client := request.DashboardClient()

	dl, err := client.List(ctx, store.Key{
		APIVersion: "jenkins.io/v1",
		Kind:       "PipelineActivity",
		Namespace:  config.Context.Namespace,
	})

	if err != nil {
		log.Logger().Infof("failed: %s", err.Error())
	}

	title := config.Title
	if title == "" {
		title = "Pipelines"
	}
	colNames := config.Columns
	if len(colNames) == 0 {
		colNames = []string{"Owner", "Repository", "Branch", "Build", "Status", "Message"}
	}
	table := component.NewTableWithRows(
		title, "There are no "+title,
		component.NewTableCols(colNames...),
		[]component.TableRow{})

	paList := []*v1.PipelineActivity{}
	if dl != nil {
		for k, v := range dl.Items {
			pa, err := viewhelpers.ToPipelineActivity(&dl.Items[k])
			if err != nil {
				log.Logger().Infof("failed to convert to PipelineActivity for %s: %s", v.GetName(), err.Error())
				continue
			}
			if pa != nil {
				paList = append(paList, pa)
			}
		}
	}
	if config.Filter != nil {
		allList := paList
		paList = []*v1.PipelineActivity{}
		for _, r := range allList {
			if config.Filter(r, allList) {
				paList = append(paList, r)
			}
		}
	}

	// default statuses
	config.StatusFilter.Values = []string{"Succeeded", "Running", "Failed"}

	for _, pa := range paList {
		tr, err := toPipelineTableRow(pa, config)
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
	table.AddFilter("Branch", config.BranchFilter)
	table.AddFilter("Status", config.StatusFilter)

	flexLayout := component.NewFlexLayout("")
	if !config.Context.Composite {
		headerText := config.Header
		if headerText == "" {
			headerText = viewhelpers.ToBreadcrumbMarkdown(plugin.RootBreadcrumb, title)
		}
		header := component.NewMarkdownText(headerText)

		flexLayout.AddSections(component.FlexLayoutSection{
			{Width: component.WidthFull, View: header},
		})
	}
	flexLayout.AddSections(component.FlexLayoutSection{
		{Width: component.WidthFull, View: table},
	})
	return flexLayout, nil
}

func (f *PipelinesViewConfig) TableFilters() []*component.TableFilter {
	return []*component.TableFilter{&f.OwnerFilter, &f.RepositoryFilter, &f.BranchFilter, &f.StatusFilter}
}

func toPipelineTableRow(pa *v1.PipelineActivity, filters *PipelinesViewConfig) (*component.TableRow, error) {
	s := &pa.Spec
	owner := s.GitOwner
	repository := s.GitRepository
	branch := s.GitBranch
	status := string(s.Status)
	viewhelpers.AddFilterValue(&filters.OwnerFilter, owner)
	viewhelpers.AddFilterValue(&filters.RepositoryFilter, repository)
	viewhelpers.AddFilterValue(&filters.BranchFilter, branch)
	viewhelpers.AddFilterValue(&filters.StatusFilter, status)
	return &component.TableRow{
		"Owner":      component.NewText(owner),
		"Repository": component.NewText(repository),
		"Branch":     component.NewText(branch),
		"Build":      ToPipelineName(pa),
		"Status":     component.NewText(status),
		"Message":    ToPipelineLastStepStatus(pa, false, true),
		"Sort":       ToOrder(pa),
	}, nil
}

func ToOrder(pa *v1.PipelineActivity) component.Component {
	s := &pa.Spec
	n := math.MaxInt64 - viewhelpers.PipelineBuildNumber(pa)

	// lets sort in most recent PR first
	lower := strings.ToLower(s.GitBranch)
	if strings.HasPrefix(lower, "pr-") {
		prNumber := strings.TrimPrefix(lower, "pr-")
		if prNumber != "" {
			pr, err := strconv.Atoi(prNumber)
			if err == nil {
				lower = fmt.Sprintf("pr-%019d", math.MaxInt64-pr)
			}
		}
	}
	order := fmt.Sprintf("%s/%s/%s/%019d/%s", s.GitOwner, s.GitRepository, lower, n, s.Context)
	return component.NewText(order)

}
