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

package v1alpha3

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
	LeaveOnDelete bool `json:"leaveOnDelete,omitempty"`
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

	// PassboltSecrets is a map of string (key in K8s secret) and struct that contains the reference to the secret in passbolt.
	// +kubebuilder:validation:Optional
	PassboltSecrets map[string]PassboltSecretRef `json:"passboltSecrets,omitempty"`

	// PlainTextFields is a map of string (key in K8s secret) and string (value in K8s secret).
	// +kubebuilder:validation:Optional
	PlainTextFields map[string]string `json:"plainTextFields,omitempty"`
}

type FieldName string

const (
	FieldNameUsername FieldName = "username"
	FieldNamePassword FieldName = "password"
	FieldNameUri      FieldName = "uri"
)

type PassboltSecretRef struct {
	// Name of the secret in passbolt
	// +kubebuilder:validation:Required
	ID string `json:"id"`
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
	// PassboltSecretID is the name of the secret that failed to sync.
	PassboltSecretID string `json:"passboltSecretID"`
	// SecretKey is the key of the secret that failed to sync.
	SecretKey string `json:"secretKey"`
	// Time is the time the error occurred.
	Time metav1.Time `json:"time"`
}

func (s SyncError) Error() string {
	return fmt.Sprintf("failed to sync secret %s/%s: %s", s.PassboltSecretID, s.SecretKey, s.Message)
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

// PassboltSecret is the Schema for the passboltsecrets API
type PassboltSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PassboltSecretSpec   `json:"spec,omitempty"`
	Status PassboltSecretStatus `json:"status,omitempty"`
}

// Hub marks this type as a conversion hub.
func (*PassboltSecret) Hub() {}

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
