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
	v1 "github.com/urbanmedia/passbolt-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// log is for logging in this package.
var passboltsecretlog = logf.Log.WithName("passboltsecret-resource")

// ConvertTo converts this PassboltSecret to the Hub version (v1).
func (src *PassboltSecret) ConvertTo(dstRaw conversion.Hub) error {
	passboltsecretlog.V(100).Info("converting PassboltSecret v1alpha3 to v1")
	dst := dstRaw.(*v1.PassboltSecret)
	dst.ObjectMeta = src.ObjectMeta
	src.Spec.LeaveOnDelete = dst.Spec.LeaveOnDelete
	dst.Spec.SecretType = src.Spec.SecretType

	// migrate secrets of type Opaque
	if src.Spec.SecretType == corev1.SecretTypeOpaque {
		dst.Spec.PassboltSecrets = make(map[string]v1.PassboltSecretRef)
		for k, v := range src.Spec.PassboltSecrets {
			dst.Spec.PassboltSecrets[k] = v1.PassboltSecretRef{
				ID:    v.ID,
				Field: v1.FieldName(v.Field),
				Value: v.Value,
			}
		}
		dst.Spec.PlainTextFields = src.Spec.PlainTextFields
	}

	// migrate secrets of type kubernetes.io/dockerconfigjson
	if src.Spec.SecretType == corev1.SecretTypeDockerConfigJson {
		dst.Spec.PassboltSecretID = src.Spec.PassboltSecretID
	}

	dst.Status.LastSync = src.Status.LastSync
	dst.Status.SyncStatus = v1.SyncStatus(src.Status.SyncStatus)
	dst.Status.SyncErrors = make([]v1.SyncError, len(src.Status.SyncErrors))
	for _, v := range src.Status.SyncErrors {
		dst.Status.SyncErrors = append(dst.Status.SyncErrors, v1.SyncError{
			Message:          v.Message,
			SecretKey:        v.SecretKey,
			PassboltSecretID: v.PassboltSecretID,
			Time:             v.Time,
		})
	}
	return nil
}

// ConvertFrom converts from the Hub version (v1alpha2) to this version.
func (dst *PassboltSecret) ConvertFrom(srcRaw conversion.Hub) error {
	passboltsecretlog.V(100).Info("converting from PassboltSecret v1 to v1alpha2")
	src := srcRaw.(*v1.PassboltSecret)
	dst.ObjectMeta = src.ObjectMeta
	dst.Spec.LeaveOnDelete = src.Spec.LeaveOnDelete
	dst.Spec.SecretType = src.Spec.SecretType

	if src.Spec.SecretType == corev1.SecretTypeOpaque {
		dst.Spec.PassboltSecrets = make(map[string]PassboltSecretRef)
		for i, s := range src.Spec.PassboltSecrets {
			dst.Spec.PassboltSecrets[i] = PassboltSecretRef{
				ID:    s.ID,
				Field: FieldName(s.Field),
				Value: s.Value,
			}
		}
		dst.Spec.PlainTextFields = src.Spec.PlainTextFields
	}

	if src.Spec.SecretType == corev1.SecretTypeDockerConfigJson {
		dst.Spec.PassboltSecretID = src.Spec.PassboltSecretID
	}

	dst.Status.LastSync = src.Status.LastSync
	dst.Status.SyncStatus = SyncStatus(src.Status.SyncStatus)
	dst.Status.SyncErrors = make([]SyncError, len(src.Status.SyncErrors))
	for _, se := range src.Status.SyncErrors {
		dst.Status.SyncErrors = append(dst.Status.SyncErrors, SyncError{
			Message:          se.Message,
			PassboltSecretID: se.PassboltSecretID,
			SecretKey:        se.SecretKey,
			Time:             se.Time,
		})
	}
	return nil
}
