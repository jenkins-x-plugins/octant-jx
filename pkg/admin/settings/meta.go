package settings // import "github.com/jenkins-x/octant-jx/pkg/plugin/settings"
import "github.com/jenkins-x/octant-jx/pkg/admin"

const (
	description = "JX OPS"
	rootNavIcon = "administrator" // See https://clarity.design/icons for all options
)

func GetName() string {
	return admin.PluginName
}

func GetDescription() string {
	return description
}
