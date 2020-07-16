package workspaces

import (
	"fmt"
	"path/filepath"

	"github.com/jenkins-x/jx-helpers/pkg/yamls"
	"github.com/jenkins-x/octant-jx/pkg/common/files"
)

type Workspace struct {
	Name        string `json:"name"`
	GitURL      string `json:"gitURL"`
	Team        string `json:"team"`
	Environment string `json:"environment"`
	ConnectCLI  string `json:"connectCLI"`
	BrowserPath string `json:"browserPath"`
	Port        int    `json:"port"`
	Default     bool   `json:"default"`
}

type WorkspaceOctant struct {
	Workspace
	URL string
}

func ToWorkspaceOctants(workspaces []*Workspace, octants *Octants) []WorkspaceOctant {
	answer := []WorkspaceOctant{}
	for _, w := range workspaces {
		wo := WorkspaceOctant{
			Workspace: *w,
		}
		path := w.BrowserPath
		if path == "" {
			path = "/#/jx/apps"
		}
		if w.Port > 0 {
			wo.URL = fmt.Sprintf("http://localhost:%d%s", w.Port, path)
		}
		answer = append(answer, wo)
	}

	for _, o := range octants.Octants {
		for i, w := range answer {
			if w.Name == o.Name {
				if o.Port > 0 {
					path := w.BrowserPath
					if path == "" {
						path = "/#/jx/apps"
					}
					answer[i].URL = fmt.Sprintf("http://localhost:%d%s", o.Port, path)
				}
				break
			}
		}
	}
	return answer
}

func LoadWorkspaces() ([]*Workspace, error) {
	answer := []*Workspace{}
	fileName := filepath.Join(files.JXOPSHomeDir(), "workspaces.yaml")
	err := yamls.LoadFile(fileName, &answer)
	return answer, err

	/*
		return []*Workspace{
			{
				Name:        "bootv3-demo",
				Team:        "labs",
				Environment: "dev",
				GitURL:      "https://github.com/cb-kubecd/environment-bootv3-demo-dev",
			},
				{
					Name:        "bootv3-demo-staging",
					Team:        "labs",
					Environment: "staging",
					GitURL:      "https://github.com/cb-kubecd/environment-bootv3-cheese2",
				},
				{
					Name:        "saas-staging",
					Team:        "saas",
					Environment: "staging",
					GitURL:      "https://github.com/cb-kubecd/environment-bootv3-stage-saas",
				},
			{
				Name:        "labs-infra",
				Team:        "labs",
				Environment: "dev",
				GitURL:      "https://github.com/jenkins-x/environment-flash-dev",
			},
		}, nil
	*/
}
