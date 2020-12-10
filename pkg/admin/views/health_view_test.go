package views

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"

	"github.com/vmware-tanzu/octant/pkg/view/component"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Test_toHealthTableRow(t *testing.T) {

	tests := []struct {
		name    string
		want    *component.TableRow
		wantErr bool
	}{
		{
			name: "kuberhealthy1.yaml",
			want: &component.TableRow{
				"Name":      viewhelpers.NewMarkdownText(`<a href="https://github.com/Comcast/kuberhealthy/blob/230c4f1/cmd/dns-resolution-check/README.md" target="docs">dns-status-internal</a>`),
				"Namespace": component.NewText("kuberhealthy"),
				"Errors":    component.NewText(""),
				"Healthy":   viewhelpers.NewMarkdownText(`<clr-icon shape="check-circle" class="is-solid is-success" title="True"></clr-icon> True`),
			},
			wantErr: false,
		},
		{
			name: "kuberhealthy2.yaml",
			want: &component.TableRow{
				"Name":      viewhelpers.NewMarkdownText(`<a href="https://github.com/Comcast/kuberhealthy/blob/230c4f1/cmd/dns-resolution-check/README.md" target="docs">dns-status-internal</a>`),
				"Namespace": component.NewText("kuberhealthy"),
				"Errors":    component.NewText("foo\nbar\n"),
				"Healthy":   viewhelpers.NewMarkdownText(`<clr-icon shape="check-circle" class="is-solid is-success" title="True"></clr-icon> True`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fileName := filepath.Join("test_data", tt.name)
			data, err := ioutil.ReadFile(fileName)
			require.NoError(t, err, "failed to load %s", fileName)

			u := &unstructured.Unstructured{}
			err = yaml.Unmarshal(data, u)
			require.NoError(t, err, "failed to unmarshal YAML %s", fileName)

			got, err := toHealthTableRow(u)
			if (err != nil) != tt.wantErr {
				t.Errorf("toHealthTableRow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toHealthTableRow() got = %v, want %v", got, tt.want)
			}
		})
	}
}
