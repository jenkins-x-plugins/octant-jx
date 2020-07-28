package views

import (
	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

var (
	noSecretDocs = `The current namespace does not appear to contain the Jenkins X Git Operator used to install and upgrade Jenkins X 3.x.

Are you sure this namespace is a Jenkins X GitOps namespace?
`

	//nolint:gosec
	settingUpSecretsDoc = `

## Create a GitOps repository and install the operator

If you do not yet have a git repostory to manage this cluster via GitOps then:

*  use the [jx admin create](https://github.com/jenkins-x/jx-admin/blob/master/docs/cmd/jx-admin_create.md) command
* see how to [Create a DevOps git repository and install the git operator](https://jenkins-x.io/docs/v3/install-setup/getting-started/repository/)

## Setting up the Jenkins X Git Operator

If you already have a git repository to manage this cluster via gitOps then:

* install the Git Operator via the [jx admin operator](https://github.com/jenkins-x/jx-admin/blob/master/docs/cmd/jx-admin_operator.md) command
* see how to [Install the Jenkins X Git Operator](https://jenkins-x.io/docs/v3/install-setup/getting-started/operator/)
`
)

// BuildNoBootSecretView view that there is no boot secret
func BuildNoBootSecretView(request service.Request, pluginContext pluginctx.Context, cr *component.ContentResponse) error {
	card := component.NewCard(component.Title(component.NewMarkdownText("## No Jenkins X Git Operator found")))
	layout := component.NewFlexLayout("jenkins x git operator")

	text := noSecretDocs + settingUpSecretsDoc

	section := component.FlexLayoutSection{
		{
			Width: component.WidthFull,
			View:  component.NewMarkdownText(text),
		},
	}
	layout.AddSections(section)
	card.SetBody(layout)

	cr.Add(card)
	return nil
}
