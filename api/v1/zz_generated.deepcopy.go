//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2024 Verlag der Tagesspiegel GmbH.

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

package v1

import (
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PassboltSecret) DeepCopyInto(out *PassboltSecret) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PassboltSecret.
func (in *PassboltSecret) DeepCopy() *PassboltSecret {
	if in == nil {
		return nil
	}
	out := new(PassboltSecret)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PassboltSecret) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PassboltSecretList) DeepCopyInto(out *PassboltSecretList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]PassboltSecret, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PassboltSecretList.
func (in *PassboltSecretList) DeepCopy() *PassboltSecretList {
	if in == nil {
		return nil
	}
	out := new(PassboltSecretList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PassboltSecretList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PassboltSecretRef) DeepCopyInto(out *PassboltSecretRef) {
	*out = *in
	if in.Value != nil {
		in, out := &in.Value, &out.Value
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PassboltSecretRef.
func (in *PassboltSecretRef) DeepCopy() *PassboltSecretRef {
	if in == nil {
		return nil
	}
	out := new(PassboltSecretRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PassboltSecretSpec) DeepCopyInto(out *PassboltSecretSpec) {
	*out = *in
	if in.PassboltSecretID != nil {
		in, out := &in.PassboltSecretID, &out.PassboltSecretID
		*out = new(string)
		**out = **in
	}
	if in.PassboltSecrets != nil {
		in, out := &in.PassboltSecrets, &out.PassboltSecrets
		*out = make(map[string]PassboltSecretRef, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
	if in.PlainTextFields != nil {
		in, out := &in.PlainTextFields, &out.PlainTextFields
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PassboltSecretSpec.
func (in *PassboltSecretSpec) DeepCopy() *PassboltSecretSpec {
	if in == nil {
		return nil
	}
	out := new(PassboltSecretSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PassboltSecretStatus) DeepCopyInto(out *PassboltSecretStatus) {
	*out = *in
	in.LastSync.DeepCopyInto(&out.LastSync)
	if in.SyncErrors != nil {
		in, out := &in.SyncErrors, &out.SyncErrors
		*out = make([]SyncError, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PassboltSecretStatus.
func (in *PassboltSecretStatus) DeepCopy() *PassboltSecretStatus {
	if in == nil {
		return nil
	}
	out := new(PassboltSecretStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SyncError) DeepCopyInto(out *SyncError) {
	*out = *in
	in.Time.DeepCopyInto(&out.Time)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SyncError.
func (in *SyncError) DeepCopy() *SyncError {
	if in == nil {
		return nil
	}
	out := new(SyncError)
	in.DeepCopyInto(out)
	return out
}
