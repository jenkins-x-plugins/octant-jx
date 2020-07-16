package main // import "github.com/jenkins-x/octant-jx/cmd/octant-jx

import (
	"fmt"
	"log"
	"os"

	"github.com/jenkins-x/octant-jx/pkg/admin/router"
	"github.com/jenkins-x/octant-jx/pkg/admin/settings"
	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
)

// Default variables overridden by ldflags
var (
	version   = "(dev-version)"
	gitCommit = "(dev-commit)"
	buildTime = "(dev-buildtime)"
)

func main() {
	args := os.Args
	if len(args) == 2 {
		switch args[1] {
		case "version":
			fmt.Println(version)
			return
		}
	}

	name := settings.GetName()
	description := settings.GetDescription()
	capabilities := settings.GetCapabilities()

	pluginContext := pluginctx.Context{
		Namespace: "jx",
	}

	h := &router.Handlers{
		Context: &pluginContext,
	}
	err := h.Load()
	if err != nil {
		panic(err)
	}

	options := settings.GetOptions(h)

	log.Printf("starting the Jenkins X plugin")

	plugin, err := service.Register(name, description, capabilities, options...)
	if err != nil {
		panic(err)
	}
	plugin.Serve()
}
