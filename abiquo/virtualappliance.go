package abiquo

import "github.com/abiquo/ojal/core"

// VirtualAppliance represents a VAPP dto
type VirtualAppliance struct {
	Name string `json:"name"`
	core.DTO
}

// CreateVM creates a VM inside v
func (v *VirtualAppliance) CreateVM(vm *VirtualMachine) error {
	return core.Create(v.Rel("virtualmachines").SetType("virtualmachine"), vm)
}
