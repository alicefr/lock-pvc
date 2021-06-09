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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// LockPVCSpec defines the desired state of LockPVC
type LockPVCSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of LockPVC. Edit lockpvc_types.go to remove/update
	Owner LockPVCOwner `json:"pcvOwner,omitempty"`
	PVC   string       `json:"pvc,omitempty"`
}

// LockPVCStatus defines the observed state of LockPVC
type LockPVCStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Acquired bool   `json:"aquired,false"`
	Message  string `json:"message,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// LockPVC is the Schema for the lockpvcs API
type LockPVC struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LockPVCSpec   `json:"spec,omitempty"`
	Status LockPVCStatus `json:"status,omitempty"`
}

// LockPVCOwner defines the resource that owns the LockPVC
type LockPVCOwner struct {
	LockPVCOwnerType metav1.TypeMeta `json:",inline"`
	LockPVCOwnerName string          `json:"owner,omitempty"`
}

//+kubebuilder:object:root=true

// LockPVCList contains a list of LockPVC
type LockPVCList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LockPVC `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LockPVC{}, &LockPVCList{})
}
