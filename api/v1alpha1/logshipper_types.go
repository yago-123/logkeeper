/*
Copyright 2025.

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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// LogShipperSpec defines the desired state of LogShipper.
type LogShipperSpec struct {
	// Image is the container image for the log shipper pod.
	// +kubebuilder:validation:MinLength=1
	Image string `json:"image"`

	// NodeSelector is used to select nodes where the log shipper pods should run.
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// Resources define the resource requests and limits for the log shipper pods.
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
}

// LogShipperStatus defines the observed state of LogShipper.
type LogShipperStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// LogShipper is the Schema for the logshippers API.
type LogShipper struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LogShipperSpec   `json:"spec,omitempty"`
	Status LogShipperStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// LogShipperList contains a list of LogShipper.
type LogShipperList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LogShipper `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LogShipper{}, &LogShipperList{})
}
