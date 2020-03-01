package main

import (
	"github.com/GuillemCastro/packer-post-processor-libvirt/libvirt"
	"github.com/hashicorp/packer/packer/plugin"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterPostProcessor(new(libvirt.PostProcessor))
	server.Serve()
}
