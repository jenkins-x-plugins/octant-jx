package views

import (
	v1 "github.com/jenkins-x/jx-api/pkg/apis/jenkins.io/v1"
	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func ToEnvironmentNameLink(r *v1.Environment) component.Component {
	name := ToEnvironmentName(r)
	ref := r.Name
	return component.NewLink(name, name, ref)
}

func ToEnvironmentNameComponent(r *v1.Environment) component.Component {
	return component.NewText(ToEnvironmentName(r))
}

func ToEnvironmentName(r *v1.Environment) string {
	s := &r.Spec
	l := s.Label
	if l == "" {
		l = r.Name
	}
	return l
}

func ToEnvironmentSource(r *v1.Environment) component.Component {
	return component.NewMarkdownText(viewhelpers.ToGitLinkMarkdown(r.Spec.Source.URL))
}

func ToEnvironmentNamespace(r *v1.Environment) component.Component {
	spec := &r.Spec
	prefix := ""
	if r.Spec.RemoteCluster {
		prefix = "Remote "
	}
	return component.NewText(prefix + spec.Namespace)
}

func ToEnvironmentRemote(r *v1.Environment) component.Component {
	text := ""
	if r.Spec.RemoteCluster {
		// TODO switch to checkbox when we can use html/markdown views https://github.com/vmware-tanzu/octant/issues/882
		text = "yes"
	}
	return component.NewText(text)
}

func ToEnvironmentPromote(r *v1.Environment) component.Component {
	return component.NewText(string(r.Spec.PromotionStrategy))
}
