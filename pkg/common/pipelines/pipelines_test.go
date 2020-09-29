package pipelines

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	v1 "github.com/jenkins-x/jx-api/pkg/apis/jenkins.io/v1"
	"github.com/jenkins-x/jx-helpers/pkg/testhelpers"
	"github.com/jenkins-x/jx-helpers/pkg/yamls"
	"github.com/stretchr/testify/require"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"sigs.k8s.io/yaml"
)

func TestInitialPipelineActivity(t *testing.T) {
	prFile := filepath.Join("test_data", "initial", "pipelinerun.yaml")
	require.FileExists(t, prFile)

	tmpDir, err := ioutil.TempDir("", "")
	require.NoError(t, err, "failed to create temp dir")

	data, err := ioutil.ReadFile(prFile)
	require.NoError(t, err, "failed to load %s", prFile)

	pr := &v1beta1.PipelineRun{}
	err = yaml.Unmarshal(data, pr)
	require.NoError(t, err, "failed to unmarshal %s", prFile)

	pa := &v1.PipelineActivity{}
	ToPipelineActivity(pr, pa)

	ClearTimestamps(pa)

	paFile := filepath.Join(tmpDir, "pa.yaml")
	err = yamls.SaveFile(pa, paFile)
	require.NoError(t, err, "failed to save %s", paFile)

	t.Logf("created PipelineActivity %s\n", paFile)

	testhelpers.AssertTextFilesEqual(t, filepath.Join("test_data", "initial", "expected.yaml"), paFile, "generated git credentials file")
}

func TestCreatePipelineActivity(t *testing.T) {
	prFile := filepath.Join("test_data", "create", "pipelinerun.yaml")
	require.FileExists(t, prFile)

	tmpDir, err := ioutil.TempDir("", "")
	require.NoError(t, err, "failed to create temp dir")

	data, err := ioutil.ReadFile(prFile)
	require.NoError(t, err, "failed to load %s", prFile)

	pr := &v1beta1.PipelineRun{}
	err = yaml.Unmarshal(data, pr)
	require.NoError(t, err, "failed to unmarshal %s", prFile)

	pa := &v1.PipelineActivity{}
	ToPipelineActivity(pr, pa)

	ClearTimestamps(pa)

	paFile := filepath.Join(tmpDir, "pa.yaml")
	err = yamls.SaveFile(pa, paFile)
	require.NoError(t, err, "failed to save %s", paFile)

	t.Logf("created PipelineActivity %s\n", paFile)

	testhelpers.AssertTextFilesEqual(t, filepath.Join("test_data", "create", "expected.yaml"), paFile, "generated git credentials file")
}

func TestMergePipelineActivity(t *testing.T) {
	prFile := filepath.Join("test_data", "merge", "pipelinerun.yaml")
	require.FileExists(t, prFile)

	paFile := filepath.Join("test_data", "merge", "pa.yaml")
	require.FileExists(t, prFile)

	tmpDir, err := ioutil.TempDir("", "")
	require.NoError(t, err, "failed to create temp dir")

	pr := &v1beta1.PipelineRun{}
	err = yamls.LoadFile(prFile, pr)
	require.NoError(t, err, "failed to load %s", prFile)

	pa := &v1.PipelineActivity{}
	err = yamls.LoadFile(paFile, pa)
	require.NoError(t, err, "failed to load %s", paFile)

	ToPipelineActivity(pr, pa)

	ClearTimestamps(pa)

	paFile = filepath.Join(tmpDir, "pa.yaml")
	err = yamls.SaveFile(pa, paFile)
	require.NoError(t, err, "failed to save %s", paFile)

	t.Logf("created PipelineActivity %s\n", paFile)

	testhelpers.AssertTextFilesEqual(t, filepath.Join("test_data", "merge", "expected.yaml"), paFile, "generated git credentials file")
}

func ClearTimestamps(pa *v1.PipelineActivity) {
	pa.Spec.StartedTimestamp = nil
	pa.Spec.CompletedTimestamp = nil
	for i := range pa.Spec.Steps {
		step := &pa.Spec.Steps[i]
		if step.Stage != nil {
			step.Stage.StartedTimestamp = nil
			step.Stage.CompletedTimestamp = nil

			for j := range step.Stage.Steps {
				s2 := &step.Stage.Steps[j]
				s2.StartedTimestamp = nil
				s2.CompletedTimestamp = nil
			}
		}
		if step.Promote != nil {
			step.Promote.StartedTimestamp = nil
			step.Promote.CompletedTimestamp = nil
		}
		if step.Preview != nil {
			step.Preview.StartedTimestamp = nil
			step.Preview.CompletedTimestamp = nil
		}
	}
}
