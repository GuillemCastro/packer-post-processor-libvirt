//go:generate mapstructure-to-hcl2 -type Config
package libvirt

import (
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/template/interpolate"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`

	ConnectionURI string `mapstructure:"connection_uri"`
	Name          string `mapstructure:"name"`
	Memory        uint   `mapstructure:"memory"`
	Cores         int    `mapstructure:"cores"`

	ctx interpolate.Context
}
