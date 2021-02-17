package views // import "github.com/jenkins-x/octant-jx/pkg/plugin/views"

import (
	"github.com/jenkins-x/jx-helpers/v3/pkg/kube/services"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx"
	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func BuildJenkinsView(request service.Request, pluginContext pluginctx.Context) (component.Component, error) {
	ctx := request.Context()
	client := request.DashboardClient()

	selector := &metav1.LabelSelector{
		MatchLabels: map[string]string{
			"app.kubernetes.io/name": "jenkins",
		},
	}

	ssetMap := map[string]*v1.StatefulSet{}
	ssl, err := client.List(ctx, store.Key{
		APIVersion:    "apps/v1",
		Kind:          "StatefulSet",
		LabelSelector: selector,
	})
	if err == nil {
		for _, u := range ssl.Items {
			sset := &v1.StatefulSet{}
			err = viewhelpers.ToStructured(&u, sset)
			if err == nil {
				key := ToKey(&sset.ObjectMeta)
				ssetMap[key] = sset
			}
		}
	}

	dl, err := client.List(ctx, store.Key{
		APIVersion:    "extensions/v1beta1",
		Kind:          "Ingress",
		LabelSelector: selector,
	})

	if err != nil {
		log.Logger().Infof("failed: %s", err.Error())
		return nil, err
	}

	log.Logger().Infof("got list of SourceRepository %d\n", len(dl.Items))

	table := component.NewTableWithRows(
		viewhelpers.TableTitle("Jenkins"), "There are no Jenkins Servers!",
		component.NewTableCols("Name"),
		[]component.TableRow{})

	for _, r := range dl.Items {
		tr, err := toJenkinsTableRow(r, ssetMap)
		if err != nil {
			log.Logger().Infof("failed to create Table Row: %s", err.Error())
			continue
		}
		if tr != nil {
			table.Add(*tr)
		}
	}

	table.Sort("Name")

	flexLayout := component.NewFlexLayout("")
	flexLayout.AddSections(component.FlexLayoutSection{
		{Width: component.WidthFull, View: table},
	})

	return flexLayout, nil
}

func ToKey(m *metav1.ObjectMeta) string {
	return m.Namespace + "/" + m.Name
}

func toJenkinsTableRow(u unstructured.Unstructured, ssetMap map[string]*v1.StatefulSet) (*component.TableRow, error) {
	r := &v1beta1.Ingress{}
	err := viewhelpers.ToStructured(&u, r)
	if err != nil {
		return nil, err
	}

	name := r.Namespace
	if name == "" {
		name = "jenkins"
	}
	labels := r.Labels
	if labels == nil || labels["app.kubernetes.io/name"] != "jenkins" {
		return nil, nil
	}

	url := services.IngressURL(r)
	md := JenkinsIconMarkdown(r, ssetMap) + " " + viewhelpers.ToMarkdownExternalLink(name, "jenkins", url)

	tr := &component.TableRow{
		"Name": viewhelpers.NewMarkdownText(md),
	}

	/* TODO add a grid action to view the credentials?
	tr.AddAction(component.GridAction{
			Name:         "Credentials",
			//ActionPath:   actionPath,
			//Payload:      payload,
			//Confirmation: confirmation,
			Type:         component.GridActionPrimary,
		})
	*/
	return tr, nil
}

func JenkinsIconMarkdown(r *v1beta1.Ingress, ssetMap map[string]*v1.StatefulSet) string {
	key := ToKey(&r.ObjectMeta)
	ssset := ssetMap[key]
	if ssset == nil || ssset.Status.Replicas <= 0 {
		return `<clr-icon shape="clock" title="Pending"></clr-icon>`
	}
	if ssset.Status.ReadyReplicas > 0 {
		return `<clr-icon shape="check-circle" class="is-solid is-success" title="Running"></clr-icon>`
	}
	return `<span class="spinner spinner-inline" title="Starting"></span>`
}
