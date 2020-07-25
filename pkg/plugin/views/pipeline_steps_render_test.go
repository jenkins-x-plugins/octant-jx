package views_test

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/jenkins-x/octant-jx/pkg/plugin/views"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

func TestPipelineStepContainerName(t *testing.T) {
	podFile := filepath.Join("test_data", "pod2.yaml")
	data, err := ioutil.ReadFile(podFile)
	require.NoError(t, err, "failed to load %s", podFile)
	pod := &corev1.Pod{}
	err = yaml.Unmarshal(data, pod)
	require.NoError(t, err, "failed to unmarshal YAML %s", podFile)

	pipelineActivityNames := []string{"Build Container Build"}
	actual := views.FindContainerNameForStepName(pod, pipelineActivityNames)
	assert.Equal(t, "step-build-container-build", actual, "for pod %s and PipelineActivity strings %s", podFile, pipelineActivityNames)
}
