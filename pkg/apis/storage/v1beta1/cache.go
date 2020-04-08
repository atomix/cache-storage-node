// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CacheStorageGroup cache storage group
const CacheStorageGroup = "storage.cloud.atomix.io"

// CacheStorageVersion cache storage version
const CacheStorageVersion = "v1beta1"

// CacheStorageKind cache storage kind
const CacheStorageKind = "CacheStorage"

// CacheStorageSpec is the k8s spec for a CacheStorage resource
type CacheStorageSpec struct {
	// Image is the image to run
	Image string `json:"image,omitempty"`

	// ImagePullPolicy is the pull policy to apply
	ImagePullPolicy corev1.PullPolicy `json:"pullPolicy,omitempty"`
}

// CacheStorageStatus defines the observed state of CacheStorage
type CacheStorageStatus struct {
	Ready bool `json:"ready,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CacheStorage is the Schema for the CacheStorage API
// +k8s:openapi-gen=true
type CacheStorage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the CacheStorage specification
	Spec CacheStorageSpec `json:"spec,omitempty"`

	// Status if the current status of the CacheStorage
	Status CacheStorageStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CacheStorageList contains a list of CacheStorage
type CacheStorageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is the set of items in the list
	Items []CacheStorage `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CacheStorage{}, &CacheStorageList{})
}
