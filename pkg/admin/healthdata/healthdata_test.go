package healthdata_test

import (
	"testing"

	"github.com/jenkins-x/octant-jx/pkg/admin/healthdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthData(t *testing.T) {
	require.NotEmpty(t, healthdata.HealthInfo, "no health data found")

	testCases := []struct {
		key      string
		expected string
	}{
		{
			key:      "deployment",
			expected: "https://github.com/Comcast/kuberhealthy/blob/230c4f1/cmd/deployment-check/README.md",
		},
	}

	for _, tc := range testCases {
		actual := healthdata.HealthInfo[tc.key]
		assert.Equal(t, tc.expected, actual, "for key %s", tc.key)
	}
}
