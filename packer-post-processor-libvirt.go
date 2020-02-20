//go:generate mapstructure-to-hcl2 -type Config
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/digitalocean/go-libvirt"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer/builder/qemu"
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/packer/plugin"
	"github.com/hashicorp/packer/template/interpolate"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`

	SocketType    string `mapstructure:"socket_type"`
	ConnectionURI string `mapstructure:"connection_uri"`

	ctx interpolate.Context
}

type PostProcessor struct {
	config Config
}

func (p *PostProcessor) ConfigSpec() hcldec.ObjectSpec {
	return p.config.FlatMapstructure().HCL2Spec()
}

func (p *PostProcessor) Configure(raws ...interface{}) error {
	err := config.Decode(&p.config, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &p.config.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{},
		},
	}, raws...)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostProcessor) PostProcess(ctx context.Context, ui packer.Ui, artifact packer.Artifact) (packer.Artifact, bool, bool, error) {
	if artifact.BuilderId() != qemu.BuilderId {
		err := fmt.Errorf(
			"Unknown artifact type: %s\nCan only import from qemu artifacts.",
			artifact.BuilderId())
		return nil, false, false, err
	}
	log.Println(artifact.String())
	c, err := net.DialTimeout(p.config.SocketType, p.config.ConnectionURI, 2*time.Second)
	if err != nil {
		return nil, false, false, err
	}
	l := libvirt.New(c)
	if err := l.Connect(); err != nil {
		log.Fatalf("failed to connect: %v", err)
		return nil, false, false, err
	}
	v, err := l.Version()
	if err != nil {
		log.Fatalf("failed to retrieve libvirt version: %v", err)
		return nil, false, false, err
	}
	log.Println("Version:", v)
	return nil, false, false, nil
}

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterPostProcessor(new(PostProcessor))
	server.Serve()
}
