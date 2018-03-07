package abiquo

import (
	"fmt"

	"github.com/abiquo/ojal/core"
)

// VirtualMachine represents an Abiquo API VMcore.Resource
type VirtualMachine struct {
	Backups   []BackupPolicy    `json:"backupPolicies,omitempty"`
	CPU       int               `json:"cpu,omitempty"`
	ID        int               `json:"id,omitempty"`
	Label     string            `json:"label,omitempty"`
	Metadata  string            `json:"metadata,omitempty"`
	Monitored bool              `json:"monitored"`
	Name      string            `json:"name,omitempty"`
	RAM       int               `json:"ram,omitempty"`
	State     string            `json:"state,omitempty"`
	Variables map[string]string `json:"variables,omitempty"`
	UUID      string            `json:"uuid,omitempty"`
	core.DTO
}

type VirtualMachineMetadata struct {
	Metadata VirtualMachineMetadataFields `json:"metadata"`
	core.DTO
}

type VirtualMachineMetadataFields struct {
	StartupScript string `json:"startup-script,omitempty"`
}

// Reconfigure reconfigures v
func (v *VirtualMachine) Reconfigure() (err error) {
	_, err = core.Rest(v, core.Put(v.URL(), "acceptedrequest", v.Media(), v))
	return
}

// Deploy deploys v
func (v *VirtualMachine) Deploy() (err error) {
	return NewTask(core.Post(
		v.Rel("deploy").Href,
		"acceptedrequest",
		"virtualmachinetask",
		v,
	))
}

// Undeploy undeploys v
func (v *VirtualMachine) Undeploy() (err error) {
	return NewTask(core.Post(
		v.Rel("undeploy").Href,
		"acceptedrequest",
		"virtualmachinetask",
		v,
	))
}

// Off power offs the VM
func (v *VirtualMachine) Off() (err error) {
	return NewTask(core.Put(
		v.Rel("state").Href,
		"acceptedrequest",
		"virtualmachinestate",
		map[string]interface{}{"state": "OFF"},
	))
}

// SetMetadata sets the VM metadata as requested
func (v *VirtualMachine) SetMetadata(metadata *VirtualMachineMetadata) error {
	return core.Update(v.Rel("metadata"), metadata)
}

// Disks returns an slice with the VM disk links
func (v *VirtualMachine) Disks() []*core.Link {
	return v.Links.Filter(func(l *core.Link) bool {
		return l.IsMedia("harddisk") || l.IsMedia("volume")
	})
}

// NICs returns an slice with the VM NIC links
func (v *VirtualMachine) NICs() (nics []*core.Link) {
	return v.Links.Filter(func(l *core.Link) bool {
		return l.IsMedia("privateip") || l.IsMedia("externalip") || l.IsMedia("publicip")
	})
}

// AttachDisk add a disk link to the *VirtualMachine in the last position
func (v *VirtualMachine) AttachDisk(hd *HardDisk) error {
	diskLink := hd.Link().SetRel(fmt.Sprintf("disk%v", len(v.Disks())))
	diskLink.DiskControllerType = hd.DiskControllerType
	diskLink.DiskController = hd.DiskController
	v.Add(diskLink)
	return nil
}

// AttachNIC add a nic link to the *VirtualMachine in the last position
func (v *VirtualMachine) AttachNIC(nic *core.Link) error {
	nicLink := nic.SetRel(fmt.Sprintf("nic%v", len(v.NICs())))
	v.Add(nicLink)
	return nil
}

func (v *VirtualMachine) Delete() error {
	return core.Remove(v)
}
