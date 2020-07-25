package viewhelpers

import (
	"fmt"
	"strings"
)

var gitHubPrefix = "https://github.com/"

func ToMarkdownLinkFromURL(link string) string {
	title := strings.TrimPrefix(link, "https://")
	title = strings.TrimPrefix(title, "http://")
	return ToMarkdownLink(title, link)
}

func ToMarkdownLink(title, link string) string {
	return fmt.Sprintf("[%s](%s)", title, link)
}

func ToMarkdownExternalLink(title, target, link string) string {
	return fmt.Sprintf(`<a href="%s" target="%s">%s</a>`, link, target, title)
}

// ToBreadcrumbMarkdown creates a breadcrumb row
func ToBreadcrumbMarkdown(names ...string) string {
	b := strings.Builder{}
	for i, name := range names {
		if i > 0 {
			b.WriteString(` <clr-icon shape="angle right"></clr-icon> `)
		}
		b.WriteString(name)
	}
	return b.String()
}

func ToGitLinkMarkdown(repoLink string) string {
	if repoLink == "" {
		return ""
	}
	title := strings.TrimSuffix(repoLink, ".git")
	if strings.HasPrefix(title, gitHubPrefix) {
		title = strings.TrimPrefix(title, gitHubPrefix)
		paths := strings.SplitN(title, "/", 2)
		if len(paths) == 2 {
			return ToOwnerRepositoryLinkMarkdown(paths[0], paths[1], repoLink)
		}
	}
	title = strings.TrimPrefix(title, "https://")
	return ToMarkdownLink(title, repoLink)
}

// ToOwnerRepositoryLinkMarkdown converts the given owner/repository and link to markdown
func ToOwnerRepositoryLinkMarkdown(owner, repository, repoLink string) string {
	if repoLink == "" {
		return fmt.Sprintf("%s / %s", owner, repository)
	}
	ownerLink := strings.TrimSuffix(repoLink, ".git")
	ownerLink = strings.TrimSuffix(ownerLink, "/")
	i := strings.LastIndex(ownerLink, "/")
	if i > 0 {
		ownerLink = ownerLink[0:i]
	}
	return fmt.Sprintf("[%s](%s) / [%s](%s)", owner, ownerLink, repository, repoLink)
}
