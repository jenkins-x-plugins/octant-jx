package views_test

import (
	"testing"

	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/stretchr/testify/assert"
)

func TestToRepositoryLinkMarkdown(t *testing.T) {
	testCases := []struct {
		Link     string
		Expected string
	}{
		{
			Link:     "https://github.com/jenkins-x/jx",
			Expected: "[jenkins-x](https://github.com/jenkins-x) / [jx](https://github.com/jenkins-x/jx)",
		},
		{
			Link:     "https://github.com/jenkins-x/jx.git",
			Expected: "[jenkins-x](https://github.com/jenkins-x) / [jx](https://github.com/jenkins-x/jx.git)",
		},
		{
			Link:     "https://gitlab.com/jenkins-x/jx.git",
			Expected: "[gitlab.com/jenkins-x/jx](https://gitlab.com/jenkins-x/jx.git)",
		},
	}

	for _, tc := range testCases {
		actual := viewhelpers.ToGitLinkMarkdown(tc.Link)
		assert.Equal(t, tc.Expected, actual, "TestToRepositoryLinkMarkdown for link %s", tc.Link)
	}
}

func TestToOwnerRepositoryLinkMarkdown(t *testing.T) {
	testCases := []struct {
		Owner      string
		Repository string
		Link       string
		Expected   string
	}{
		{
			Owner:      "jenkins-x",
			Repository: "jx",
			Link:       "https://github.com/jenkins-x/jx",
			Expected:   "[jenkins-x](https://github.com/jenkins-x) / [jx](https://github.com/jenkins-x/jx)",
		},
		{
			Owner:      "jenkins-x",
			Repository: "jx",
			Link:       "https://github.com/jenkins-x/jx.git",
			Expected:   "[jenkins-x](https://github.com/jenkins-x) / [jx](https://github.com/jenkins-x/jx.git)",
		},
		{
			Owner:      "jenkins-x",
			Repository: "jx",
			Link:       "https://gitlab.com/jenkins-x/jx.git",
			Expected:   "[jenkins-x](https://gitlab.com/jenkins-x) / [jx](https://gitlab.com/jenkins-x/jx.git)",
		},
	}

	for _, tc := range testCases {
		actual := viewhelpers.ToOwnerRepositoryLinkMarkdown(tc.Owner, tc.Repository, tc.Link)
		assert.Equal(t, tc.Expected, actual, "TestToRepositoryLinkMarkdown for link %s", tc.Link)
	}
}
