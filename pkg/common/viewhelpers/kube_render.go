package viewhelpers

import "github.com/vmware-tanzu/octant/pkg/view/component"

func ViewPipelineLogs(ns, podName string, containers ...string) (*component.Logs, error) {
	logsView := component.NewLogs(ns, podName, containers...)
	return logsView, nil
}
