package actioners

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/jenkins-x/jx-logging/v3/pkg/log"

	"github.com/pkg/errors"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
)

func HandleTriggerBootJob(request *service.ActionRequest) error {
	log.Logger().Infof("triggering boot job")

	ns, err := request.Payload.String("namespace")
	if err != nil {
		// Todo: Not sure what to do with error ...
		log.Logger().Info(err)
	}

	if ns == "" {
		// ToDO: not sure what is being done here, ns is not used at all
		ns = "jx" //nolint:ineffassign
	}

	return runCommand("jxl", "boot", "step", "redirect", "-s", "jx-boot-octant", "-c", "jxl boot run -b --no-tail")
}

func runCommand(name string, args ...string) error {
	tmpDir, err := ioutil.TempDir("", "jxl-boot-run-")
	if err != nil {
		return errors.Wrap(err, "failed to create temp dir for boot run")
	}

	cmd := exec.Command(name, args...)
	cmd.Dir = tmpDir
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	log.Logger().Infof("running %s %s", name, strings.Join(args, " "))
	err = cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "failed to execute: jxl boot run")
	}

	log.Logger().Infof("boot run completed!")
	return nil
}
