package actioners

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
)

func HandleTriggerBootJob(request *service.ActionRequest) error {
	log.Printf("triggering boot job")

	ns, err := request.Payload.String("namespace")
	if err != nil || ns == "" {
		ns = "jx"
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

	log.Printf("running %s %s", name, strings.Join(args, " "))
	err = cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "failed to execute: jxl boot run")
	}

	log.Printf("boot run completed!")
	return nil
}
