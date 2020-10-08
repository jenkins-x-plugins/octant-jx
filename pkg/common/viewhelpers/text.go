package viewhelpers

import "github.com/vmware-tanzu/octant/pkg/view/component"

// NewMarkdownText creates a trusted markdown object
func NewMarkdownText(text string) *component.Text {
	answer := component.NewMarkdownText(text)
	answer.EnableTrustedContent()
	return answer
}
