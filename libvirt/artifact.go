package libvirt

import (
	"fmt"
	"os"
)

type Artifact struct {
	outputFile string
	VMName     string
	state      map[string]interface{}
}

func (Artifact) BuilderId() string {
	return BuilderId
}

func (a Artifact) Files() []string {
	return []string{a.outputFile}
}

func (a Artifact) Id() string {
	return a.VMName
}

func (a Artifact) String() string {
	return fmt.Sprintf("[%s]: %s", a.VMName, a.outputFile)
}

func (a Artifact) State(name string) interface{} {
	return a.state[name]
}

func (a Artifact) Destroy() error {
	return os.Remove(a.outputFile)
}
