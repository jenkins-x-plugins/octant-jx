package links

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetImageLink(t *testing.T) {
	t.Parallel()

	testCases := map[string]string{
		"":            "",
		"maven":       "https://hub.docker.com/_/maven",
		"maven:1.2.3": "https://hub.docker.com/_/maven:1.2.3",
		"index.docker.io/jenkinsxio/builder-go:1.2.3":                          "https://hub.docker.com/r/jenkinsxio/builder-go:1.2.3",
		"unknown.foo/bar/whatnot":                                              "",
		"gcr.io/jenkinsxio/builder-go:2.0.1099-435":                            "https://gcr.io/jenkinsxio/builder-go:2.0.1099-435",
		"gcr.io/abayer-pipeline-crd/tekton-for-jx/git-init:v20200414-2b72e7c6": "https://gcr.io/abayer-pipeline-crd/tekton-for-jx/git-init:v20200414-2b72e7c6",
	}

	for image, expected := range testCases {
		actual := GetImageLink(image)

		assert.Equal(t, expected, actual, "for image %s", image)
	}
}
