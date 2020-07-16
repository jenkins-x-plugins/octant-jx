package views

import (
	"fmt"

	"github.com/jenkins-x/jx-helpers/pkg/gitclient/giturl"
	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

var (
	noSecretDocs = `The current namespace does not appear to contain the Jenkins X Boot Secret used to indicate the secrets have been setup for use with helm 3.

Are you sure this namespace is a Jenkins X boot namespace?
`
	exampleGitURL = "https://github.com/myorg/envronment-mycluster.git"

	settingUpSecretsDoc = `

## Setting up the Secrets

To setup secrets try the following:

    git clone %s
	cd %s
    jxl boot secrets edit 

That should setup the boot secret properly.
`
)

// BuildNoBootSecretView view that there is no boot secret
func BuildNoBootSecretView(request service.Request, pluginContext pluginctx.Context, cr *component.ContentResponse) error {
	card := component.NewCard(component.Title(component.NewMarkdownText("## No Boot Secret found")))
	layout := component.NewFlexLayout("starting boot")

	text := noSecretDocs + GetBootEditSecretsMarkdown(exampleGitURL)

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

func GetBootEditSecretsMarkdown(gitURL string) string {
	if gitURL == "" {
		gitURL = exampleGitURL
	}

	// lets find the folder
	repo, _ := giturl.ParseGitURL(gitURL)
	dir := "myrepo"
	if repo != nil && repo.Name != "" {
		dir = repo.Name
	}
	return fmt.Sprintf(settingUpSecretsDoc, gitURL, dir)
}
