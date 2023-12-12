/*
Copyright 2023 Verlag der Tagesspiegel GmbH.

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

package v1alpha2

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PassboltSecretSpec defines the desired state of PassboltSecret
type PassboltSecretSpec struct {
	// LeaveOnDelete defines if the secret should be deleted from Kubernetes when the PassboltSecret is deleted.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=true
	LeaveOnDelete bool `json:"leaveOnDelete"`
	// SecretType is the type of the secret. Defaults to Opaque.
	// If set to kubernetes.io/dockerconfigjson, the secret will be created as a docker config secret.
	// We also expect the PassboltSecretName to be set in this case.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=Opaque
	// +kubebuilder:validation:Enum=Opaque;kubernetes.io/dockerconfigjson
	SecretType corev1.SecretType `json:"secretType,omitempty"`
	// PassboltSecretName is the name of the passbolt secret name to be used as a docker config secret.
	// +kubebuilder:validation:Optional
	PassboltSecretName *string `json:"passboltSecretName,omitempty"`
	// Secrets is a list of secrets to be fetched from passbolt.
	// +kubebuilder:validation:Optional
	Secrets []SecretSpec `json:"secrets,omitempty"`
}

// SecretSpec defines the secret mapping between passbolt and kubernetes.
type SecretSpec struct {
	// Name of the secret in passbolt
	// +kubebuilder:validation:Required
	PassboltSecret PassboltSpec `json:"passboltSecret"`
	// KubernetesSecretKey is the key in the kubernetes secret where the passbolt secret will be stored.
	// +kubebuilder:validation:Required
	KubernetesSecretKey string `json:"kubernetesSecretKey"`
}

type FieldName string

const (
	FieldNameUsername FieldName = "username"
	FieldNamePassword FieldName = "password"
	FieldNameUri      FieldName = "uri"
)

type PassboltSpec struct {
	// Name of the secret in passbolt
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// Field is the field in the passbolt secret to be read.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=username;password;uri
	Field FieldName `json:"field,omitempty"`
	// Value is the plain text value of the secret.
	// This field allows to set a static value or using go templating to generate the value.
	// Valid template variables are:
	//   - Password
	//   - Username
	//   - URI
	// +kubebuilder:validation:Optional
	Value *string `json:"value,omitempty"`
}

type SyncStatus string

const (
	SyncStatusSuccess SyncStatus = "Success"
	SyncStatusError   SyncStatus = "Error"
	SyncStatusUnknown SyncStatus = "Unknown"
)

type SyncError struct {
	// Message is the error message.
	Message string `json:"message"`
	// SecretName is the name of the secret that failed to sync.
	SecretName string `json:"secretName"`
	// SecretKey is the key of the secret that failed to sync.
	SecretKey string `json:"secretKey"`
	// Time is the time the error occurred.
	Time metav1.Time `json:"time"`
}

func (s SyncError) Error() string {
	return fmt.Sprintf("failed to sync secret %s/%s: %s", s.SecretName, s.SecretKey, s.Message)
}

// PassboltSecretStatus defines the observed state of PassboltSecret
type PassboltSecretStatus struct {
	// SyncStatus is the status of the last sync.
	// +kubebuilder:validation:Enum=Success;Error;Unknown
	// +kubebuilder:default=Unknown
	SyncStatus SyncStatus `json:"syncStatus"`
	// LastSync is the last time the secret was synced from passbolt.
	// +kubebuilder:validation:Optional
	LastSync metav1.Time `json:"lastSync"`
	// SyncErrors is a list of errors that occurred during the last sync.
	SyncErrors []SyncError `json:"syncErrors,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Sync Status",type=string,JSONPath=`.status.syncStatus`
// +kubebuilder:printcolumn:name="Last Sync",type=string,JSONPath=`.status.lastSync`

// PassboltSecret is the Schema for the passboltsecrets API
type PassboltSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PassboltSecretSpec   `json:"spec,omitempty"`
	Status PassboltSecretStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PassboltSecretList contains a list of PassboltSecret
type PassboltSecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PassboltSecret `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PassboltSecret{}, &PassboltSecretList{})
}
