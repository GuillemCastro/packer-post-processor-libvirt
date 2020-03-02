# packer-post-processor-libvirt ![Go](https://github.com/GuillemCastro/packer-post-processor-libvirt/workflows/Go/badge.svg?branch=master)

A Packer post-processor that automatically imports your built image to your preferred virtualization platform using libvirt.

## Status

- [x] Import a image built with the qemu builder
- [x] Configure memory and number of cores
- [ ] Configure devices like network interfaces and graphic adapters

## Installation and usage

For installing the post-processor, just execute,

```bash
go get -v github.com/GuillemCastro/packer-post-processor-libvirt
```

Then in `~/.packerconfig` add the libvirt post-processor. It should look like this,

```json
{
    "post-processors": {
        "libvirt": "packer-post-processor-libvirt"
    }
}
```

Then you can add it to your image configuration for packer,

```json
{
  "builders": [
    {
      "type": "qemu", //right now only qemu is supported
      "format": "qcow2", //with qcow2 images
      "accelerator": "kvm", //and kvm as the accelerator
      "vm_name": "ubuntu1804", //this will be the name of your VM
    }
  ],
  "post-processors": [
      {
          "type": "libvirt",
          "connection_uri": "qemu:///system", // QEMU connection URI (optional) by default "qemu:///system"
          "memory": 1024, // Ammount of RAM (optional) by default 1024MB
          "cores": 1, // Ammount of CPU cores (optional) by default 1
      }
  ]
}
```

## License

```
MIT License

Copyright (c) 2020 Guillem Castro

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
