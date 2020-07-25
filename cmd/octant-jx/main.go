package main // import "github.com/jenkins-x/octant-jx/cmd/octant-jx

import (
	"fmt"
	"os"

	"github.com/jenkins-x/jx-logging/pkg/log"

	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"

	"github.com/jenkins-x/octant-jx/pkg/plugin/settings"
)

// Default variables overridden by ldflags
var (
	version = "(dev-version)"
)

func main() {
	args := os.Args
	if len(args) == 2 && args[1] == "version" {
		fmt.Println(version)
		return
	}

	name := settings.GetName()
	description := settings.GetDescription()
	capabilities := settings.GetCapabilities()

	pluginContext := pluginctx.Context{
		Namespace: "jx",
	}

	options := settings.GetOptions(&pluginContext)

	log.Logger().Infof("starting the Jenkins X plugin")

	plugin, err := service.Register(name, description, capabilities, options...)
	if err != nil {
		panic(err)
	}
	plugin.Serve()
}
