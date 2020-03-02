package libvirt

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer/builder/qemu"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/template/interpolate"
	"github.com/libvirt/libvirt-go"
	libvirtxml "github.com/libvirt/libvirt-go-xml"
)

const BuilderId = "libvirt"

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
	if p.config.ConnectionURI == "" {
		p.config.ConnectionURI = "qemu:///system"
	}
	if p.config.Memory == 0 {
		p.config.Memory = 1024
	}
	if p.config.Cores < 1 {
		p.config.Cores = 1
	}
	return nil
}

func getArch() string {
	switch runtime.GOARCH {
	case "amd64":
		return "x86_64"
	default:
		return runtime.GOARCH
	}
}

func (p *PostProcessor) PostProcess(ctx context.Context, ui packer.Ui, artifact packer.Artifact) (packer.Artifact, bool, bool, error) {
	if artifact.BuilderId() != qemu.BuilderId {
		err := fmt.Errorf(
			"Unknown artifact type: %s\nCan only import from qemu artifacts.",
			artifact.BuilderId())
		return nil, false, false, err
	}
	var drive uint = 0
	conn, err := libvirt.NewConnect(p.config.ConnectionURI)
	if err != nil {
		return nil, false, false, err
	}
	diskFilePath, err := filepath.Abs(artifact.Files()[0])
	_, err = os.Stat(diskFilePath)
	if os.IsNotExist(err) {
		err := fmt.Errorf("Disk file: %s does not exist.", diskFilePath)
		return nil, false, false, err
	}
	arch := getArch()
	vmName := artifact.State("diskName").(string)
	domainDefinition := &libvirtxml.Domain{
		Type:   artifact.State("domainType").(string),
		Name:   vmName,
		Memory: &libvirtxml.DomainMemory{Value: p.config.Memory, Unit: "MB", DumpCore: "on"},
		VCPU:   &libvirtxml.DomainVCPU{Value: p.config.Cores},
		CPU:    &libvirtxml.DomainCPU{Mode: "host-model"},
		Devices: &libvirtxml.DomainDeviceList{
			Disks: []libvirtxml.DomainDisk{
				{
					Source:  &libvirtxml.DomainDiskSource{File: &libvirtxml.DomainDiskSourceFile{File: diskFilePath}},
					Target:  &libvirtxml.DomainDiskTarget{Dev: "hda", Bus: "ide"},
					Alias:   &libvirtxml.DomainAlias{Name: "ide0-0-0"},
					Address: &libvirtxml.DomainAddress{Drive: &libvirtxml.DomainAddressDrive{Controller: &drive, Bus: &drive, Target: &drive, Unit: &drive}},
					Driver:  &libvirtxml.DomainDiskDriver{Name: "qemu", Type: "qcow2"},
				},
			},
			Graphics: []libvirtxml.DomainGraphic{
				{
					Spice: &libvirtxml.DomainGraphicSpice{Port: 5900, AutoPort: "yes", Listen: "127.0.0.1"},
				},
			},
			Interfaces: []libvirtxml.DomainInterface{
				{
					Model: &libvirtxml.DomainInterfaceModel{Type: "virtio"},
				},
			},
		},
		OS: &libvirtxml.DomainOS{
			Type: &libvirtxml.DomainOSType{
				Arch: arch,
				Type: "hvm",
			},
		},
	}
	xml, err := domainDefinition.Marshal()
	if err != nil {
		return nil, false, false, err
	}
	_, err = conn.DomainDefineXML(xml)
	if err != nil {
		return nil, false, false, err
	}
	newArtifact := &Artifact{
		outputFile: diskFilePath,
		VMName:     vmName,
		state:      make(map[string]interface{}),
	}
	newArtifact.state["arch"] = arch
	return newArtifact, true, true, nil
}
