//go:build !ignore_autogenerated

/*
Copyright 2024.

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
	"k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppStack) DeepCopyInto(out *AppStack) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppStack.
func (in *AppStack) DeepCopy() *AppStack {
	if in == nil {
		return nil
	}
	out := new(AppStack)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AppStack) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppStackList) DeepCopyInto(out *AppStackList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]AppStack, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppStackList.
func (in *AppStackList) DeepCopy() *AppStackList {
	if in == nil {
		return nil
	}
	out := new(AppStackList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AppStackList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppStackSpec) DeepCopyInto(out *AppStackSpec) {
	*out = *in
	in.Workload.DeepCopyInto(&out.Workload)
	in.Service.DeepCopyInto(&out.Service)
	in.HPA.DeepCopyInto(&out.HPA)
	in.Ingress.DeepCopyInto(&out.Ingress)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppStackSpec.
func (in *AppStackSpec) DeepCopy() *AppStackSpec {
	if in == nil {
		return nil
	}
	out := new(AppStackSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppStackStatus) DeepCopyInto(out *AppStackStatus) {
	*out = *in
	in.DeploymentStatus.DeepCopyInto(&out.DeploymentStatus)
	in.ServiceStatus.DeepCopyInto(&out.ServiceStatus)
	in.IngressStatus.DeepCopyInto(&out.IngressStatus)
	in.HPAStatus.DeepCopyInto(&out.HPAStatus)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppStackStatus.
func (in *AppStackStatus) DeepCopy() *AppStackStatus {
	if in == nil {
		return nil
	}
	out := new(AppStackStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HPA) DeepCopyInto(out *HPA) {
	*out = *in
	if in.Min != nil {
		in, out := &in.Min, &out.Min
		*out = new(int32)
		**out = **in
	}
	if in.Max != nil {
		in, out := &in.Max, &out.Max
		*out = new(int32)
		**out = **in
	}
	in.Metrics.DeepCopyInto(&out.Metrics)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HPA.
func (in *HPA) DeepCopy() *HPA {
	if in == nil {
		return nil
	}
	out := new(HPA)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Ingress) DeepCopyInto(out *Ingress) {
	*out = *in
	if in.Rules != nil {
		in, out := &in.Rules, &out.Rules
		*out = make([]IngressRule, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.TLS != nil {
		in, out := &in.TLS, &out.TLS
		*out = make([]networkingv1.IngressTLS, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.IngressCLass != nil {
		in, out := &in.IngressCLass, &out.IngressCLass
		*out = new(string)
		**out = **in
	}
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Ingress.
func (in *Ingress) DeepCopy() *Ingress {
	if in == nil {
		return nil
	}
	out := new(Ingress)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngressRule) DeepCopyInto(out *IngressRule) {
	*out = *in
	if in.Paths != nil {
		in, out := &in.Paths, &out.Paths
		*out = make([]Path, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressRule.
func (in *IngressRule) DeepCopy() *IngressRule {
	if in == nil {
		return nil
	}
	out := new(IngressRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Path) DeepCopyInto(out *Path) {
	*out = *in
	if in.PathType != nil {
		in, out := &in.PathType, &out.PathType
		*out = new(networkingv1.PathType)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Path.
func (in *Path) DeepCopy() *Path {
	if in == nil {
		return nil
	}
	out := new(Path)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Service) DeepCopyInto(out *Service) {
	*out = *in
	if in.Ports != nil {
		in, out := &in.Ports, &out.Ports
		*out = make([]v1.ServicePort, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Service.
func (in *Service) DeepCopy() *Service {
	if in == nil {
		return nil
	}
	out := new(Service)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkLoad) DeepCopyInto(out *WorkLoad) {
	*out = *in
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	if in.PodAnnotations != nil {
		in, out := &in.PodAnnotations, &out.PodAnnotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.DeploymentAnnotations != nil {
		in, out := &in.DeploymentAnnotations, &out.DeploymentAnnotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make([]v1.Volume, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.VolumeMounts != nil {
		in, out := &in.VolumeMounts, &out.VolumeMounts
		*out = make([]v1.VolumeMount, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.EnvFrom != nil {
		in, out := &in.EnvFrom, &out.EnvFrom
		*out = make([]v1.EnvFromSource, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkLoad.
func (in *WorkLoad) DeepCopy() *WorkLoad {
	if in == nil {
		return nil
	}
	out := new(WorkLoad)
	in.DeepCopyInto(out)
	return out
}
