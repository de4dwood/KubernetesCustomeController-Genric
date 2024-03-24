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

package v1alpha1

import (
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AppStackSpec defines the desired state of AppStack

type WorkLoad struct {
	Replicas              *int32                 `json:"replicas"`
	Image                 string                 `json:"image"`
	PodAnnotations        map[string]string      `json:"podannotations,omitempty"`
	DeploymentAnnotations map[string]string      `json:"deploymentannotations,omitempty"`
	Volumes               []corev1.Volume        `json:"volumes,omitempty"`
	VolumeMounts          []corev1.VolumeMount   `json:"volumemounts,omitempty"`
	EnvFrom               []corev1.EnvFromSource `json:"envfrom,omitempty"`
	Env                   []corev1.EnvVar        `json:"env,omitempty"`
}

type Service struct {
	Type        string               `json:"type,omitempty"`
	Ports       []corev1.ServicePort `json:"ports,omitempty"`
	Annotations map[string]string    `json:"serviceannotations,omitempty"`
}

type HPA struct {
	Min     *int32                   `json:"min,omitempty"`
	Max     *int32                   `json:"max,omitempty"`
	Metrics autoscalingv2.MetricSpec `json:"metrics,omitempty"`
}

type Ingress struct {
	Rules        networkingv1.IngressRule `json:"rules"`
	TLS          networkingv1.IngressTLS  `json:"tls,omitempty"`
	IngressCLass *string                  `json:"ingressclass,omitempty"`
	Annotations  map[string]string        `json:"annotations,omitempty"`
}

type AppStackSpec struct {
	Workload WorkLoad `json:"workload"`
	Service  Service  `json:"service,omitempty"`
	HPA      HPA      `json:"hpa,omitempty"`
	Ingress  Ingress  `json:"ingress,omitempty"`
}

// AppStackStatus defines the observed state of AppStack
type AppStackStatus struct {
	DeploymentStatus appsv1.DeploymentStatus                     `json:"deploymentStatus"`
	ServiceStatus    corev1.ServiceStatus                        `json:"serviceStatus,omitempty"`
	IngressStatus    networkingv1.IngressStatus                  `json:"ingressStatus,omitempty"`
	HPAStatus        autoscalingv2.HorizontalPodAutoscalerStatus `json:"hpaStatus,omitempty"`

	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// AppStack is the Schema for the appstacks API
type AppStack struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AppStackSpec   `json:"spec,omitempty"`
	Status AppStackStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AppStackList contains a list of AppStack
type AppStackList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AppStack `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AppStack{}, &AppStackList{})
}
