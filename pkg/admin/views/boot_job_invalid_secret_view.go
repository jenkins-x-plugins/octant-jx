package views

import (
	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

var (
	invalidBootSecretsDocs = `You need to edit the secrets to populate missing values.

`
)

// BuildBootInvalidSecretView view that the secrets are not yet valid
func BuildBootInvalidSecretView(request service.Request, pluginContext pluginctx.Context, cr *component.ContentResponse, gitURL string) error {
	card := component.NewCard(component.Title(component.NewMarkdownText("## Invalid Jenkins X GitOps Secrets")))

	text := invalidBootSecretsDocs + settingUpSecretsDoc

	layout := component.NewFlexLayout("starting boot")
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
