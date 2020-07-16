package views

import (
	"fmt"
	"log"
	"strings"

	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx"
	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	corev1 "k8s.io/api/core/v1"
)

// BootJobExtraView lets render the output of any previous commands
func BootJobExtraView(request service.Request, pluginContext pluginctx.Context, cr *component.ContentResponse, view component.Component) error {
	ctx := request.Context()
	client := request.DashboardClient()

	name := "jx-boot-octant"
	ns := pluginContext.Namespace
	u, err := viewhelpers.GetResourceByName(ctx, client, "v1", "Secret", name, ns)
	if err != nil {
		return nil
	}
	if u == nil {
		return nil
	}

	secret := &corev1.Secret{}
	err = viewhelpers.ToStructured(u, secret)
	if err != nil {
		log.Println(err)
		return nil
	}
	data := secret.Data
	if len(data) == 0 {
		return nil
	}
	b := strings.Builder{}
	for i := 0; i < len(data); i++ {
		k := fmt.Sprintf("l%d", i)
		line := data[k]
		if line == nil {
			log.Printf("Secret %s in namespace %s does not have key %s", name, ns, k)
			return nil
		}
		b.WriteString(string(line))
		b.WriteString("\n")
	}

	flexLayout, ok := view.(*component.FlexLayout)
	if !ok {
		log.Printf("view is not a FlexLayout - was %#v", view)
		return nil
	}
	title := RunCommandTitleMarkdown(secret)
	card := component.NewCard(component.Title(component.NewMarkdownText(title)))
	layout := component.NewFlexLayout("starting boot")

	section := component.FlexLayoutSection{
		{
			Width: component.WidthFull,
			View:  component.NewMarkdownText(b.String()),
		},
	}
	layout.AddSections(section)
	card.SetBody(layout)

	flexLayout.AddSections(component.FlexLayoutSection{
		{Width: component.WidthFull, View: card},
	})
	return nil
}

func RunCommandTitleMarkdown(secret *corev1.Secret) string {
	annotations := secret.Annotations
	if annotations == nil {
		annotations = map[string]string{}
	}
	started := annotations["jenkins.io/started"]
	completed := annotations["jenkins.io/completed"]
	failed := annotations["jenkins.io/failed"]

	message := "Triggering the Boot Job"
	icon := `<clr-icon shape="clock" title="Pending"></clr-icon>`

	at := started
	if completed != "" {
		icon = `<clr-icon shape="check-circle" class="is-solid is-success" title="Succeeded"></clr-icon>`
		message = "Triggered the Boot Job"
		at = completed
	} else if failed != "" {
		icon = `<clr-icon shape="warning-standard" class="is-solid is-danger" title="Failed"></clr-icon>`
		message = "Failed to trigger the Boot Job"
		at = completed
	} else if started != "" {
		icon = `<span class="spinner spinner-inline" title="Running"></span>`
	}
	suffix := ""
	if at != "" {
		suffix = " at " + at
	}
	return fmt.Sprintf("%s&nbsp;**%s**%s", icon, message, suffix)
}
