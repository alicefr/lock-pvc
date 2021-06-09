// +build !ignore_autogenerated

/*
Copyright 2021.

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

package v1beta1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LockPVC) DeepCopyInto(out *LockPVC) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LockPVC.
func (in *LockPVC) DeepCopy() *LockPVC {
	if in == nil {
		return nil
	}
	out := new(LockPVC)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LockPVC) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LockPVCList) DeepCopyInto(out *LockPVCList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]LockPVC, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LockPVCList.
func (in *LockPVCList) DeepCopy() *LockPVCList {
	if in == nil {
		return nil
	}
	out := new(LockPVCList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LockPVCList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LockPVCOwner) DeepCopyInto(out *LockPVCOwner) {
	*out = *in
	out.LockPVCOwnerType = in.LockPVCOwnerType
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LockPVCOwner.
func (in *LockPVCOwner) DeepCopy() *LockPVCOwner {
	if in == nil {
		return nil
	}
	out := new(LockPVCOwner)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LockPVCSpec) DeepCopyInto(out *LockPVCSpec) {
	*out = *in
	out.Owner = in.Owner
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LockPVCSpec.
func (in *LockPVCSpec) DeepCopy() *LockPVCSpec {
	if in == nil {
		return nil
	}
	out := new(LockPVCSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LockPVCStatus) DeepCopyInto(out *LockPVCStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LockPVCStatus.
func (in *LockPVCStatus) DeepCopy() *LockPVCStatus {
	if in == nil {
		return nil
	}
	out := new(LockPVCStatus)
	in.DeepCopyInto(out)
	return out
}
