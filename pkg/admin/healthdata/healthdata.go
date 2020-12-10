package healthdata

import (
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/jenkins-x/octant-jx/pkg/assets"
)

var (
	HealthInfo = map[string]string{}
)

func init() {
	data, err := assets.Asset("files/health.yaml")
	if err != nil {
		fmt.Printf("warning: could not find health yaml: %s", err.Error())
		return
	}

	m := map[string]string{}
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		fmt.Printf("warning: unmarshal health yaml %s", err.Error())
	}
	for k, v := range m {
		HealthInfo[k] = v
	}
}
