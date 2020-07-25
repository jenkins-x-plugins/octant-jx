package views // import "github.com/jenkins-x/octant-jx/pkg/plugin/views"

import (
	"fmt"
	"strings"

	"github.com/jenkins-x/jx-logging/pkg/log"

	"github.com/jenkins-x/octant-jx/pkg/common/links"
	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx"
	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/jenkins-x/octant-jx/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func BuildAppsView(request service.Request, pluginContext pluginctx.Context) (component.Component, error) {
	ctx := request.Context()
	client := request.DashboardClient()
	ns := pluginContext.Namespace

	dl, err := client.List(ctx, store.Key{
		APIVersion: "apps/v1",
		Kind:       "Deployment",
		Namespace:  ns,
	})

	if err != nil {
		log.Logger().Infof("failed to load Deployments: %s", err.Error())
		return nil, err
	}

	ingList, err := client.List(ctx, store.Key{
		APIVersion: "extensions/v1beta1",
		Kind:       "Ingress",
		Namespace:  ns,
	})

	ingresses := []*v1beta1.Ingress{}
	if err != nil {
		log.Logger().Infof("failed to load Ingress: %s", err.Error())
	} else {
		for k := range ingList.Items {
			ing := &v1beta1.Ingress{}
			err = viewhelpers.ToStructured(&ingList.Items[k], ing)
			if err != nil {
				log.Logger().Infof("failed to convert to Ingress: %s", err.Error())
			} else {
				ingresses = append(ingresses, ing)
			}
		}
	}

	log.Logger().Debugf("got list of Deployment %d\n", len(dl.Items))

	header := component.NewMarkdownText(viewhelpers.ToBreadcrumbMarkdown(plugin.RootBreadcrumb, "Apps"))

	table := component.NewTableWithRows(
		"Apps", "There are no Apps!",
		component.NewTableCols("Name", "Version", "Pods", "URL"),
		[]component.TableRow{})

	icons := map[string]string{}
	deployments := []*appsv1.Deployment{}
	for k := range dl.Items {
		r := &appsv1.Deployment{}
		err := viewhelpers.ToStructured(&dl.Items[k], r)
		if err != nil {
			return nil, err
		}
		deployments = append(deployments, r)
		icon := ToDeploymentIcon(r)
		if icon != "" {
			icons[r.Name] = icon
		}
	}

	for _, r := range deployments {
		icon := icons[r.Name]
		if icon == "" && len(icons) > 0 {
			icon = DefaultIcon
		}
		tr, err := toAppTableRow(r, icon, ingresses)
		if err != nil {
			log.Logger().Infof("failed to create Table Row: %s", err.Error())
			continue
		}
		if tr != nil {
			table.Add(*tr)
		}
	}

	table.Sort("Name", false)
	flexLayout := component.NewFlexLayout("")
	flexLayout.AddSections(component.FlexLayoutSection{
		{Width: component.WidthFull, View: header},
		{Width: component.WidthFull, View: table},
	})

	return flexLayout, nil
}

func toAppTableRow(r *appsv1.Deployment, icon string, ingresses []*v1beta1.Ingress) (*component.TableRow, error) {
	name, version := ToResourceNameVersion(&r.ObjectMeta)

	return &component.TableRow{
		"Name":    ToDeploymentLink(r, name, icon),
		"Version": ToDeploymentVersion(r, version),
		"Pods":    ToDeploymentStatus(r),
		"URL":     component.NewMarkdownText(viewhelpers.ToMarkdownLinkFromURL(FindIngressLinkByAppName(name, ingresses))),
	}, nil
}

func ToDeploymentVersion(r *appsv1.Deployment, version string) component.Component {
	if r != nil {
		ann := r.Annotations
		if ann != nil {
			home := ann[AnnotationHome]
			if home != "" {
				return component.NewMarkdownText(viewhelpers.ToMarkdownLink(version, home))
			}
		}
	}
	return component.NewText(version)
}

func ToDeploymentLink(r *appsv1.Deployment, name, icon string) component.Component {
	ref := links.GetDeploymentLink(r.Namespace, r.Name)
	iconPrefix := ""
	if icon != "" {
		iconPrefix = fmt.Sprintf(`<img src="%s" width="24" height="24">&nbsp;`, icon)
	}
	return component.NewMarkdownText(fmt.Sprintf(`%s<a href="%s" title="Deployment %s">%s</a>`, iconPrefix, ref, name, name))
}

func ToDeploymentIcon(r *appsv1.Deployment) string {
	icon := ""
	if r.Annotations != nil {
		icon = viewhelpers.ToApplicationIcon(r.Annotations[AnnotationIcon])
	}
	return icon
}

func ToDeploymentStatus(r *appsv1.Deployment) component.Component {
	if r == nil {
		return component.NewText("")
	}
	s := &r.Status
	replicas := s.Replicas
	availableReplicas := s.AvailableReplicas
	clazz := "badge-info"
	if replicas > 0 {
		if availableReplicas <= 0 {
			clazz = "badge-danger"
		} else if availableReplicas < replicas {
			clazz = "badge-warning"
		}
	}
	return component.NewMarkdownText(fmt.Sprintf(`<span class="badge badge-info" title="Replicas">%d</span>/&nbsp;<span class="badge %s" title="Available Pods">%d</span>`, replicas, clazz, availableReplicas))
}

func FindIngressLinkByAppName(name string, ingresses []*v1beta1.Ingress) string {
	for _, ing := range ingresses {
		ingName, _ := ToResourceNameVersion(&ing.ObjectMeta)
		if ingName == name {
			link := ToIngressLink(ing)
			if link != "" {
				return link
			}
		}
	}
	return ""
}

func ToIngressLink(ing *v1beta1.Ingress) string {
	if ing != nil {
		if len(ing.Spec.Rules) > 0 {
			rule := ing.Spec.Rules[0]
			hostname := rule.Host
			for _, tls := range ing.Spec.TLS {
				for _, h := range tls.Hosts {
					if h != "" {
						url := "https://" + h
						return url
					}
				}
			}
			ann := ing.Annotations
			if hostname == "" && ann != nil {
				hostname = ann[AnnotationHost]
			}
			if hostname != "" {
				url := "http://" + hostname
				if rule.HTTP != nil {
					if len(rule.HTTP.Paths) > 0 {
						p := rule.HTTP.Paths[0].Path
						if p != "" {
							url += p
						}
					}
				}
				return url
			}
		}
	}
	return ""
}

func ToResourceNameVersion(r *metav1.ObjectMeta) (string, string) {
	name := ""
	version := ""
	labels := r.Labels
	if labels != nil {
		name = labels[LabelAppName]
		version = labels[LabelAppVersion]
		app := labels["app"]
		if app != "" {
			lastIdx := strings.LastIndex(app, "-")
			if lastIdx > 0 {
				if name == "" {
					name = app[lastIdx+1:]
				}
			}
		}
		chart := labels[LabelHelmChart]
		if chart == "" {
			chart = labels[LabelJXChart]
		}
		if chart != "" {
			lastIdx := strings.LastIndex(chart, "-")
			if lastIdx > 0 {
				if name == "" {
					name = chart[0:lastIdx]
				}
				if version == "" {
					version = chart[lastIdx+1:]
				}
			}
		}
	}
	if name == "" {
		name = r.Name
	}
	return name, version
}
