package pipelines

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/jenkins-x/jx-helpers/pkg/testhelpers"
	"github.com/jenkins-x/jx-helpers/pkg/yamls"
	"github.com/stretchr/testify/require"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"sigs.k8s.io/yaml"
)

func TestToPipelineActivity(t *testing.T) {
	prFile := filepath.Join("test_data", "pipelinerun.yaml")
	require.FileExists(t, prFile)

	tmpDir, err := ioutil.TempDir("", "")
	require.NoError(t, err, "failed to create temp dir")

	data, err := ioutil.ReadFile(prFile)
	require.NoError(t, err, "failed to load %s", prFile)

	pr := &v1beta1.PipelineRun{}
	err = yaml.Unmarshal(data, pr)
	require.NoError(t, err, "failed to unmarshal %s", prFile)

	pa := ToPipelineActivity(pr)

	paFile := filepath.Join(tmpDir, "pa.yaml")
	err = yamls.SaveFile(pa, paFile)
	require.NoError(t, err, "failed to save %s", paFile)

	t.Logf("created PipelineActivity %s\n", paFile)

	testhelpers.AssertTextFilesEqual(t, filepath.Join("test_data", "expected.yaml"), paFile, "generated git credentials file")
}
