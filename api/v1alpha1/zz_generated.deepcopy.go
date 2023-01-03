//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2022 @ Verlag Der Tagesspiegel GmbH

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
func (in *PassboltSecretSpec) DeepCopyInto(out *PassboltSecretSpec) {
	*out = *in
	if in.Secrets != nil {
		in, out := &in.Secrets, &out.Secrets
		*out = make([]SecretSpec, len(*in))
		copy(*out, *in)
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
func (in *PassboltSpec) DeepCopyInto(out *PassboltSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PassboltSpec.
func (in *PassboltSpec) DeepCopy() *PassboltSpec {
	if in == nil {
		return nil
	}
	out := new(PassboltSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretSpec) DeepCopyInto(out *SecretSpec) {
	*out = *in
	out.PassboltSecret = in.PassboltSecret
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretSpec.
func (in *SecretSpec) DeepCopy() *SecretSpec {
	if in == nil {
		return nil
	}
	out := new(SecretSpec)
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
