package workspaces

import (
	"path/filepath"

	"github.com/jenkins-x/jx-helpers/pkg/yamls"
	"github.com/jenkins-x/octant-jx/pkg/common/files"
)

type Octants struct {
	Octants  []Octant
	fileName string
}

type Octant struct {
	Name           string `json:"name"`
	Dir            string `json:"dir"`
	KubeConfigPath string `json:"kubeConfig"`
	Port           int    `json:"port"`
	Admin          bool   `json:"admin"`
}

func NewOctants() (*Octants, error) {
	o := &Octants{
		fileName: filepath.Join(files.JXOPSHomeDir(), "octants.yaml"),
	}
	err := yamls.LoadFile(o.fileName, &o.Octants)
	if err != nil {
		return nil, err
	}
	return o, nil
}

// Get gets the octant values
func (o *Octants) Get(workspace *Workspace) Octant {
	name := workspace.Name
	var answer Octant
	for _, r := range o.Octants {
		if r.Name == name {
			answer = r
			break
		}
	}
	answer.Name = name
	if workspace.Default {
		answer.Admin = true
	}
	if workspace.Port > 0 && answer.Port <= 0 {
		answer.Port = workspace.Port
	}
	return answer
}

// Set updates the octant values returning true if its a new octant
func (o *Octants) Set(values Octant) bool {
	for i, r := range o.Octants {
		if r.Name == values.Name {
			o.Octants[i] = values
			return false
		}
	}
	o.Octants = append(o.Octants, values)
	return true
}

func (o *Octants) Save() error {
	return SaveOctants(o.Octants)
}

func LoadOctants() ([]Octant, error) {
	answer := []Octant{}
	fileName := filepath.Join(files.JXOPSHomeDir(), "octants.yaml")
	err := yamls.LoadFile(fileName, &answer)
	return answer, err
}

func SaveOctants(octants []Octant) error {
	fileName := filepath.Join(files.JXOPSHomeDir(), "octants.yaml")
	err := yamls.SaveFile(octants, fileName)
	return err
}
