package views // import "github.com/jenkins-x/octant-jx/pkg/plugin/views"

import (
	"strings"

	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"github.com/jenkins-x/octant-jx/pkg/common/links"

	"github.com/jenkins-x/jx-secret/pkg/apis/external/v1"
	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx"
	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/jenkins-x/octant-jx/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type SecretsViewConfig struct {
}

func BuildSecretsView(request service.Request, pluginContext pluginctx.Context) (component.Component, error) {
	ctx := request.Context()
	client := request.DashboardClient()

	ns := pluginContext.Namespace
	dl, err := client.List(ctx, store.Key{
		APIVersion: "kubernetes-client.io/v1",
		Kind:       "ExternalSecret",
		Namespace:  ns,
	})

	if err != nil {
		log.Logger().Infof("failed: %s", err.Error())
		return nil, err
	}

	secretMap := map[string]*corev1.Secret{}

	sl, err := client.List(ctx, store.Key{
		APIVersion: "v1",
		Kind:       "Secret",
		Namespace:  ns,
	})
	if err == nil {
		for i := range sl.Items {
			u := &sl.Items[i]
			s := &corev1.Secret{}
			err := viewhelpers.ToStructured(u, s)
			if err == nil {
				secretMap[s.Name] = s
			}
		}
	}

	log.Logger().Infof("got list of Preview %d\n", len(dl.Items))

	header := viewhelpers.NewMarkdownText(viewhelpers.ToBreadcrumbMarkdown(plugin.RootBreadcrumb, "Secrets"))

	config := &SecretsViewConfig{}

	table := component.NewTableWithRows(
		"Secrets", "There are no Secrets!",
		component.NewTableCols("Name", "Status"),
		[]component.TableRow{})

	for _, es := range dl.Items {
		tr, err := toSecretTableRow(ns, es, secretMap, config)
		if err != nil {
			log.Logger().Infof("failed to create Table Row: %s", err.Error())
			continue
		}
		if tr != nil {
			table.Add(*tr)
		}
	}

	table.Sort("Sort", false)

	flexLayout := component.NewFlexLayout("")
	flexLayout.AddSections(component.FlexLayoutSection{
		{Width: component.WidthFull, View: header},
		{Width: component.WidthFull, View: table},
	})

	return flexLayout, nil
}

func toSecretTableRow(ns string, u unstructured.Unstructured, secretMap map[string]*corev1.Secret, filters *SecretsViewConfig) (*component.TableRow, error) {
	r := &v1.ExternalSecret{}
	err := viewhelpers.ToStructured(&u, r)
	if err != nil {
		return nil, err
	}

	name := r.Name
	nameView := component.NewText(name)
	secret := secretMap[name]
	if secret != nil {
		nameView = viewhelpers.NewMarkdownText(viewhelpers.ToMarkdownLink(name, links.GetSecretLink(ns, name)))
	}

	return &component.TableRow{
		"Sort":   component.NewText(name),
		"Name":   nameView,
		"Status": ToSecretStatus(r, secretMap),
	}, nil
}

func ToSecretStatus(r *v1.ExternalSecret, secretMap map[string]*corev1.Secret) component.Component {
	status := ""
	name := r.Name
	switch r.Status.Status {
	case "SUCCESS":
		status = `<clr-icon shape="check-circle" class="is-solid is-success" title="Secret is populated"></clr-icon>`
	default:
		status = `<clr-icon shape="warning-standard" class="is-solid is-danger" title="secret not populated"></clr-icon>`

		var secretData map[string][]byte
		secret := secretMap[name]
		if secret != nil {
			secretData = secret.Data
		}
		if secretData == nil {
			secretData = map[string][]byte{}
		}

		var missingKeys []string
		for _, d := range r.Spec.Data {
			key := d.Key
			v := secretData[key]
			if len(v) == 0 {
				missingKeys = append(missingKeys, key)
			}
		}
		if len(missingKeys) > 0 {
			status += " missing: " + strings.Join(missingKeys, ", ")
		}
	}
	return viewhelpers.NewMarkdownText(status)
}
