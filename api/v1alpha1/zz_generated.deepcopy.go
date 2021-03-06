//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2022 Alex.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AptPackage) DeepCopyInto(out *AptPackage) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AptPackage.
func (in *AptPackage) DeepCopy() *AptPackage {
	if in == nil {
		return nil
	}
	out := new(AptPackage)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DesiredPackages) DeepCopyInto(out *DesiredPackages) {
	*out = *in
	if in.Apt != nil {
		in, out := &in.Apt, &out.Apt
		*out = make([]AptPackage, len(*in))
		copy(*out, *in)
	}
	if in.Snap != nil {
		in, out := &in.Snap, &out.Snap
		*out = make([]SnapPackage, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DesiredPackages.
func (in *DesiredPackages) DeepCopy() *DesiredPackages {
	if in == nil {
		return nil
	}
	out := new(DesiredPackages)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Module) DeepCopyInto(out *Module) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Module.
func (in *Module) DeepCopy() *Module {
	if in == nil {
		return nil
	}
	out := new(Module)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Node) DeepCopyInto(out *Node) {
	*out = *in
	if in.Modules != nil {
		in, out := &in.Modules, &out.Modules
		*out = make([]Module, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Node.
func (in *Node) DeepCopy() *Node {
	if in == nil {
		return nil
	}
	out := new(Node)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SnapPackage) DeepCopyInto(out *SnapPackage) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SnapPackage.
func (in *SnapPackage) DeepCopy() *SnapPackage {
	if in == nil {
		return nil
	}
	out := new(SnapPackage)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UbuntuMachineConfiguration) DeepCopyInto(out *UbuntuMachineConfiguration) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UbuntuMachineConfiguration.
func (in *UbuntuMachineConfiguration) DeepCopy() *UbuntuMachineConfiguration {
	if in == nil {
		return nil
	}
	out := new(UbuntuMachineConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *UbuntuMachineConfiguration) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UbuntuMachineConfigurationList) DeepCopyInto(out *UbuntuMachineConfigurationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]UbuntuMachineConfiguration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UbuntuMachineConfigurationList.
func (in *UbuntuMachineConfigurationList) DeepCopy() *UbuntuMachineConfigurationList {
	if in == nil {
		return nil
	}
	out := new(UbuntuMachineConfigurationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *UbuntuMachineConfigurationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UbuntuMachineSpec) DeepCopyInto(out *UbuntuMachineSpec) {
	*out = *in
	if in.DesiredModules != nil {
		in, out := &in.DesiredModules, &out.DesiredModules
		*out = make([]Module, len(*in))
		copy(*out, *in)
	}
	in.DesiredPackages.DeepCopyInto(&out.DesiredPackages)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UbuntuMachineSpec.
func (in *UbuntuMachineSpec) DeepCopy() *UbuntuMachineSpec {
	if in == nil {
		return nil
	}
	out := new(UbuntuMachineSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UbuntuMachineStatus) DeepCopyInto(out *UbuntuMachineStatus) {
	*out = *in
	if in.Nodes != nil {
		in, out := &in.Nodes, &out.Nodes
		*out = make([]Node, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UbuntuMachineStatus.
func (in *UbuntuMachineStatus) DeepCopy() *UbuntuMachineStatus {
	if in == nil {
		return nil
	}
	out := new(UbuntuMachineStatus)
	in.DeepCopyInto(out)
	return out
}
