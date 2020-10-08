package views // import "github.com/jenkins-x/octant-jx/pkg/plugin/views"

import (
	"fmt"
	"strings"

	"github.com/jenkins-x/jx-logging/v3/pkg/log"

	"github.com/jenkins-x/octant-jx/pkg/admin"
	"github.com/jenkins-x/octant-jx/pkg/common/actions"
	"github.com/jenkins-x/octant-jx/pkg/common/links"
	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx"
	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
)

type JobViewConfig struct {
	ButtonGroup        *component.ButtonGroup
	ActionName         string
	ActionResourceName string
	Title              string
	Selector           labels.Set
}

var (
	JobsViewConfigs = map[string]*JobViewConfig{
		admin.BootJobsPath: {
			ActionName:         actions.TriggerBootJob,
			ActionResourceName: "jx-boot",
			Title:              "Boot Jobs",
			Selector: labels.Set{
				"app": "jx-boot",
			},
		},
		admin.GCPipelineJobsPath: {
			Title: "GC Pipeline Jobs",
			Selector: labels.Set{
				"app": "gcactivities",
			},
		},
		admin.GCPodJobsPath: {
			Title: "GC Preview Jobs",
			Selector: labels.Set{
				"app": "gcpods",
			},
		},
		admin.GCPreviewJobsPath: {
			Title: "GC Preview Jobs",
			Selector: labels.Set{
				"app": "gcpreviews",
			},
		},
		admin.UpgradeJobsPath: {
			Title: "Upgrade Jobs",
			Selector: labels.Set{
				"app": "jenkins-x-upgrade-processor",
			},
		},
	}
)

func BuildJobsViewForPath(request service.Request, pluginContext pluginctx.Context, path string) (component.Component, error) {
	config := JobsViewConfigs[path]
	if config == nil {
		return component.NewText(fmt.Sprintf("No view configuration found for path %s", path)), nil
	}

	ctx := request.Context()
	client := request.DashboardClient()
	ns := pluginContext.Namespace

	selector := config.Selector
	title := config.Title
	if title == "" {
		title = "Jobs"
	}

	jobs := []*batchv1.Job{}
	jl, err := viewhelpers.ListResourcesBySelector(ctx, client, "batch/v1", "Job", ns, selector)
	if err != nil {
		log.Logger().Infof("failed to load Jobs: %s", err.Error())
	} else {
		if len(jl.Items) == 0 {
			log.Logger().Infof("could not find any Jobs in namespace %s for selector %#v", ns, selector)
		}
		for k := range jl.Items {
			j := &batchv1.Job{}
			err = viewhelpers.ToStructured(&jl.Items[k], j)
			if err != nil {
				log.Logger().Infof("failed to convert to Job: %s", err.Error())
			} else {
				jobs = append(jobs, j)
			}
		}
	}

	tableTitle := title
	table := component.NewTableWithRows(
		tableTitle, "There are no "+title,
		component.NewTableCols("Name", "Pods", "Age", ""),
		[]component.TableRow{})

	for _, r := range jobs {
		tr, err := toJobTableRow(r, path)
		if err != nil {
			log.Logger().Infof("failed to create Table Row: %s", err.Error())
			continue
		}
		if tr != nil {
			table.Add(*tr)
		}
	}

	table.Sort("Age", true)
	flexLayout := component.NewFlexLayout("")

	if !pluginContext.Composite {
		header := viewhelpers.NewMarkdownText(viewhelpers.ToBreadcrumbMarkdown(admin.RootBreadcrumb, title))

		buttonGroup := config.ButtonGroup
		if buttonGroup == nil {
			cronJobList, err := viewhelpers.ListResourcesBySelector(ctx, client, "batch/v1", "CronJob", ns, selector)
			if err != nil {
				log.Logger().Infof("failed to load CronJobs: %s", err.Error())
			}
			if cronJobList != nil {
				buttonGroup = createJobButtons(pluginContext, config, cronJobList, title, selector)
			}
		}
		if buttonGroup != nil {
			flexLayout.AddSections(component.FlexLayoutSection{
				{Width: component.WidthHalf, View: header},
				{Width: component.WidthHalf, View: buttonGroup},
			})
		} else {
			flexLayout.AddSections(component.FlexLayoutSection{
				{Width: component.WidthFull, View: header},
			})
		}
	}
	flexLayout.AddSections(component.FlexLayoutSection{
		{Width: component.WidthFull, View: table},
	})
	return flexLayout, nil
}

func toJobTableRow(r *batchv1.Job, path string) (*component.TableRow, error) {
	return &component.TableRow{
		"Name": ToJobName(r),
		"Pods": ToJobPods(r),
		"Age":  component.NewTimestamp(r.CreationTimestamp.Time),
		"":     viewhelpers.NewMarkdownText(fmt.Sprintf(`<a href='%s' class="btn btn-info-outline btn-sm">Logs</a>`, admin.JobsLogsViewLink(path, r.Name))),
	}, nil
}

func createJobButtons(pluginContext pluginctx.Context, config *JobViewConfig, list *unstructured.UnstructuredList, title string, selector labels.Set) *component.ButtonGroup {
	name := config.ActionResourceName
	if name == "" {
		if len(list.Items) != 1 {
			for i, u := range list.Items {
				log.Logger().Infof("ignored CronJob %d with name %s", i, u.GetName())
			}
			if len(list.Items) == 0 {
				log.Logger().Infof("no CronJobs found for selector %#v", selector)
			}
			// TODO try find the CronJob using labels/annotations?
			return nil
		}
		name = list.Items[0].GetName()
	}
	if name == "" {
		return nil
	}

	buttonGroup := component.NewButtonGroup()
	actionName := config.ActionName
	if actionName == "" {
		actionName = "action.octant.dev/cronJob"
	}
	buttonGroup.AddButton(component.NewButton("Manually Trigger", action.CreatePayload(actionName, action.Payload{
		"namespace":  pluginContext.Namespace,
		"apiVersion": "batch/v1",
		"kind":       "CronJob",
		"name":       name,
	}), JobConfirmation(title)))
	return buttonGroup
}

// JobConfirmation for title
func JobConfirmation(title string) component.ButtonOption {
	confirmationTitle := fmt.Sprintf("Trigger %s Job", title)
	confirmationBody := fmt.Sprintf("Are you sure you want to trigger the **%s**?", strings.TrimSuffix(title, "s"))
	return component.WithButtonConfirmation(confirmationTitle, confirmationBody)
}

func ToJobName(r *batchv1.Job) component.Component {
	name := r.Name
	ref := links.GetJobLink(r.Namespace, name)
	iconPrefix := ToJobIcon(r)
	if iconPrefix != "" {
		iconPrefix += "&nbsp;&nbsp;"
	}
	return viewhelpers.NewMarkdownText(fmt.Sprintf(`%s<a href="%s" title="Deployment %s">%s</a>`, iconPrefix, ref, name, name))
}

func ToJobPods(r *batchv1.Job) component.Component {
	s := r.Status
	b := strings.Builder{}
	if s.Active > 0 {
		if b.Len() > 0 {
			b.WriteString(" ")
		}
		b.WriteString(fmt.Sprintf(`<span class="badge badge-info" title="Running pods">%d</span>`, s.Active))
	}
	if s.Succeeded > 0 {
		if b.Len() > 0 {
			b.WriteString(" ")
		}
		b.WriteString(fmt.Sprintf(`<span class="badge badge-success" title="Pods succeeded">%d</span>`, s.Succeeded))
	}
	if s.Failed > 0 {
		if b.Len() > 0 {
			b.WriteString(" ")
		}
		b.WriteString(fmt.Sprintf(`<span class="badge badge-danger" title="Pods failed">%d</span>`, s.Failed))
	}
	return viewhelpers.NewMarkdownText(b.String())
}

func ToJobIcon(r *batchv1.Job) string {
	s := r.Status
	if s.Succeeded > 0 {
		return `<clr-icon shape="check-circle" class="is-solid is-success" title="Succeeded"></clr-icon>`
	}
	if s.Active > 0 {
		return `<span class="spinner spinner-inline" title="Running"></span>`
	}
	if s.Failed > 0 {
		return `<clr-icon shape="warning-standard" class="is-solid is-danger" title="Failed"></clr-icon>`
	}
	return `<clr-icon shape="clock" title="Pending"></clr-icon>`

}
