package views // import "github.com/jenkins-x/octant-jx/pkg/plugin/views"

import (
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"github.com/jenkins-x/octant-jx/pkg/admin/healthdata"

	"github.com/jenkins-x/octant-jx/pkg/admin"
	"github.com/pkg/errors"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx"
	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func HealthView(request service.Request, _ pluginctx.Context) (component.Component, error) {

	ctx := request.Context()
	client := request.DashboardClient()

	checks, err := client.List(ctx, store.Key{
		APIVersion: "comcast.github.io/v1",
		Kind:       "KuberhealthyCheck",
	})

	if err != nil {
		log.Logger().Infof("failed: %s", err.Error())
		return nil, err
	}

	h, err := client.List(ctx, store.Key{
		APIVersion: "comcast.github.io/v1",
		Kind:       "KuberhealthyState",
	})

	annotate(checks, h)

	if err != nil {
		log.Logger().Infof("failed: %s", err.Error())
		return nil, err
	}

	log.Logger().Infof("got list of KuberhealthyState %d\n", len(h.Items))

	header := viewhelpers.NewMarkdownText(viewhelpers.ToBreadcrumbMarkdown(admin.RootBreadcrumb, "Health"))

	table := component.NewTableWithRows(
		"Health", "There are no Health statuses!",
		component.NewTableCols("Name", "Namespace", "Healthy", "Errors"),
		[]component.TableRow{})

	for k := range h.Items {
		tr, err := toHealthTableRow(&h.Items[k])
		if err != nil {
			log.Logger().Infof("failed to create Table Row: %s", err.Error())
			continue
		}
		if tr != nil {
			table.Add(*tr)
		}
	}

	table.Sort("Name")
	flexLayout := component.NewFlexLayout("Health")
	flexLayout.AddSections(component.FlexLayoutSection{
		{Width: component.WidthFull, View: header},
		{Width: component.WidthFull, View: table},
	})

	return flexLayout, nil

}

func annotate(checks *unstructured.UnstructuredList, states *unstructured.UnstructuredList) {
	m := map[string]map[string]string{}
	for _, u := range checks.Items {
		name := u.GetName()
		ann := u.GetAnnotations()
		if ann != nil && name != "" {
			m[name] = ann
		}
	}
	for i := range states.Items {
		u := &states.Items[i]
		name := u.GetName()
		ann := m[name]
		if len(ann) == 0 {
			continue
		}
		curr := u.GetAnnotations()
		if curr == nil {
			curr = map[string]string{}
		}
		for k, v := range ann {
			if curr[k] == "" {
				curr[k] = v
			}
		}
		u.SetAnnotations(curr)
	}
}

func toHealthTableRow(u *unstructured.Unstructured) (*component.TableRow, error) {
	ann := u.GetAnnotations()
	if ann == nil {
		ann = map[string]string{}
	}
	name := u.GetName()
	namespace, _, err := unstructured.NestedString(u.Object, "metadata", "namespace")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get OK from spec %s", name)
	}

	status, _, err := unstructured.NestedBool(u.Object, "spec", "OK")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get OK from spec %s", name)
	}

	statusComment := ""
	if status {
		statusComment = `<clr-icon shape="check-circle" class="is-solid is-success" title="True"></clr-icon> True`
	} else {
		statusComment = `<clr-icon shape="warning-standard" class="is-solid is-danger" title="False"></clr-icon> False`
	}

	healthErrorList, _, err := unstructured.NestedStringSlice(u.Object, "spec", "Errors")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get Errors from spec %s", name)
	}
	healthErrorMessage := ""
	for _, healthError := range healthErrorList {
		healthErrorMessage = healthErrorMessage + healthError + "\n"
	}

	nameComponent := component.NewText(name)

	docLink := ann["docs.jenkins-x.io"]
	if docLink == "" {
		docLink = healthdata.HealthInfo[name]
	}
	if docLink != "" {
		nameComponent = viewhelpers.NewMarkdownText(viewhelpers.ToMarkdownExternalLink(name, "docs", docLink))
	}

	return &component.TableRow{
		"Name":      nameComponent,
		"Namespace": component.NewText(namespace),
		"Healthy":   viewhelpers.NewMarkdownText(statusComment),
		"Errors":    component.NewText(healthErrorMessage),
	}, nil
}
