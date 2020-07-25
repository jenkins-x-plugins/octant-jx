package viewhelpers

import "github.com/vmware-tanzu/octant/pkg/view/component"

// NewTextCard helper function to create a card with a text body
func NewTextCard(title, bodyText string) *component.Card {
	notesCard := component.NewCard(component.TitleFromString(title))
	notesBody := component.NewText(bodyText)
	notesCard.SetBody(notesBody)
	return notesCard
}

// NewMarkdownCard helper function to create a card with a text body
func NewMarkdownCard(title, bodyText string) *component.Card {
	notesCard := component.NewCard(component.TitleFromString(title))
	notesBody := component.NewMarkdownText(bodyText)
	notesCard.SetBody(notesBody)
	return notesCard
}
