package settings // import "github.com/jenkins-x/octant-jx/pkg/plugin/settings"
import "github.com/jenkins-x/octant-jx/pkg/plugin"

const (
	description = "Jenkins X support"
	rootNavIcon = "ci-cd" // See https://clarity.design/icons for all options
)

func GetName() string {
	return plugin.Name
}

func GetDescription() string {
	return description
}
