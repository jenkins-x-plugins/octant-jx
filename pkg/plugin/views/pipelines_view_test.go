package views_test

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	v1 "github.com/jenkins-x/jx-api/v3/pkg/apis/jenkins.io/v1"
	"github.com/jenkins-x/octant-jx/pkg/plugin/views"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"
)

func TestPipelineLastStep(t *testing.T) {
	type testCase struct {
		Name     string
		Expected string
		LinkURL  string
	}

	testCases := []testCase{
		{
			Name:     "pa-preview.yaml",
			Expected: `<clr-icon shape="check-circle" class="is-solid is-success" title="Succeeded"></clr-icon> <span>Promote <a href="http://node-demo-js-15-jx-cb-kubecd-node-demo-js-15-pr-1.146.148.5.128.nip.io">Preview</a></span>`,
		},
		{
			Name:     "pa5.yaml",
			Expected: `<span class="spinner spinner-inline" title="Running"></span> <span>Promote to Staging <a href="https://github.com/cb-kubecd/environment-bootv3-demo-staging/pull/7">#7</a></span>`,
		},
		{
			Name:     "pa1.yaml",
			Expected: `<span class="spinner spinner-inline" title="Running"></span> Place Tools`,
		},
		{
			Name:     "pa2.yaml",
			Expected: `<span class="spinner spinner-inline" title="Running"></span> Create Tekton Crds`,
		},
		{
			Name:     "pa3.yaml",
			Expected: `<span class="spinner spinner-inline" title="Running"></span> Credential Initializer Bwl4v`,
		},
		{
			Name:     "pa4.yaml",
			Expected: `<span class="spinner spinner-inline" title="Running"></span> Build Container Build`,
		},
	}

	for _, tc := range testCases {
		fileName := filepath.Join("test_data", tc.Name)
		data, err := ioutil.ReadFile(fileName)
		require.NoError(t, err, "failed to load %s", fileName)
		pa := &v1.PipelineActivity{}
		err = yaml.Unmarshal(data, pa)
		require.NoError(t, err, "failed to unmarshal YAML %s", fileName)

		actual := views.ToPipelineLastStepStatus(pa, false, false)
		actualText := actual.String()

		t.Logf("ToPipelineLastStepStatus for %s generated %s", tc.Name, actualText)
		assert.Equal(t, tc.Expected, actualText, "ToPipelineLastStepStatus for test %s", tc.Name)
	}
}
