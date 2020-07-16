package viewhelpers

import (
	"strconv"

	v1 "github.com/jenkins-x/jx-api/pkg/apis/jenkins.io/v1"
)

// PipelineBuildNumber  returns the build number for the pipeline
func PipelineBuildNumber(pa *v1.PipelineActivity) int {
	build := pa.Spec.Build
	if build == "" {
		return 0
	}
	n, err := strconv.Atoi(build)
	if err != nil {
		return -1
	}
	return n
}
