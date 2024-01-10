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

package v1alpha2

import (
	"fmt"

	"github.com/urbanmedia/passbolt-operator/api/v1alpha3"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

var (
	// GetByID is a function that returns the name of a secret by its ID.
	GetSecretID func(name string) (string, error) = nil
	// GetBySecretName is a function that returns the ID of a secret by its name.
	GetSecretName func(id string) (string, error) = nil
)

// ConvertTo converts this CronJob to the Hub version (v1alpha2).
func (src *PassboltSecret) ConvertTo(dstRaw conversion.Hub) error {
	passboltsecretlog.V(100).Info("converting PassboltSecret v1alpha2 to v1alpha3")
	dst := dstRaw.(*v1alpha3.PassboltSecret)
	dst.ObjectMeta = src.ObjectMeta
	src.Spec.LeaveOnDelete = dst.Spec.LeaveOnDelete
	dst.Spec.SecretType = src.Spec.SecretType

	// migrate secrets of type Opaque
	if src.Spec.SecretType == corev1.SecretTypeOpaque {
		dst.Spec.PassboltSecrets = make(map[string]v1alpha3.PassboltSecretRef)
		for i, s := range src.Spec.Secrets {
			pbID, err := GetSecretID(s.PassboltSecret.Name)
			if err != nil {
				return fmt.Errorf("error migrating secret %s at index %d: %w", s.PassboltSecret.Name, i, err)
			}
			dst.Spec.PassboltSecrets[s.KubernetesSecretKey] = v1alpha3.PassboltSecretRef{
				ID:    pbID,
				Field: v1alpha3.FieldName(s.PassboltSecret.Field),
				Value: s.PassboltSecret.Value,
			}
		}
	}

	// migrate secrets of type kubernetes.io/dockerconfigjson
	if src.Spec.SecretType == corev1.SecretTypeDockerConfigJson {
		pbID, err := GetSecretID(*src.Spec.PassboltSecretName)
		if err != nil {
			return fmt.Errorf("error migrating secret %s in namespace %s: %w", src.GetName(), src.GetNamespace(), err)
		}
		dst.Spec.PassboltSecretID = &pbID
	}

	dst.Status.LastSync = src.Status.LastSync
	dst.Status.SyncStatus = v1alpha3.SyncStatus(src.Status.SyncStatus)
	dst.Status.SyncErrors = make([]v1alpha3.SyncError, len(src.Status.SyncErrors))
	for _, se := range src.Status.SyncErrors {
		dst.Status.SyncErrors = append(dst.Status.SyncErrors, v1alpha3.SyncError{
			Message:          se.Message,
			PassboltSecretID: se.SecretName,
			SecretKey:        se.SecretKey,
			Time:             se.Time,
		})
	}
	return nil
}

// ConvertFrom converts from the Hub version (v1alpha2) to this version.
func (dst *PassboltSecret) ConvertFrom(srcRaw conversion.Hub) error {
	passboltsecretlog.V(100).Info("converting from PassboltSecret v1alpha3 to v1alpha2")
	src := srcRaw.(*v1alpha3.PassboltSecret)
	dst.ObjectMeta = src.ObjectMeta
	dst.Spec.LeaveOnDelete = src.Spec.LeaveOnDelete
	dst.Spec.SecretType = src.Spec.SecretType

	if src.Spec.SecretType == corev1.SecretTypeOpaque {
		dst.Spec.Secrets = []SecretSpec{}
		for i, s := range src.Spec.PassboltSecrets {
			id, err := GetSecretName(s.ID)
			if err != nil {
				return fmt.Errorf("error migrating secret %s at index %s: %w", s.ID, i, err)
			}
			dst.Spec.Secrets = append(dst.Spec.Secrets, SecretSpec{
				KubernetesSecretKey: i,
				PassboltSecret: PassboltSpec{
					Name:  id,
					Field: FieldName(s.Field),
					Value: s.Value,
				},
			})
		}
		for i, s := range src.Spec.PlainTextFields {
			dst.Spec.Secrets = append(dst.Spec.Secrets, SecretSpec{
				KubernetesSecretKey: i,
				PassboltSecret: PassboltSpec{
					Name:  i,
					Value: &s,
				},
			})
		}
	}

	if src.Spec.SecretType == corev1.SecretTypeDockerConfigJson {
		name, err := GetSecretName(*src.Spec.PassboltSecretID)
		if err != nil {
			return fmt.Errorf("error migrating secret %s: %w", *src.Spec.PassboltSecretID, err)
		}
		dst.Spec.PassboltSecretName = &name
	}

	dst.Status.LastSync = src.Status.LastSync
	dst.Status.SyncStatus = SyncStatus(src.Status.SyncStatus)
	dst.Status.SyncErrors = make([]SyncError, len(src.Status.SyncErrors))
	for _, se := range src.Status.SyncErrors {
		dst.Status.SyncErrors = append(dst.Status.SyncErrors, SyncError{
			Message:    se.Message,
			SecretName: se.PassboltSecretID,
			SecretKey:  se.SecretKey,
			Time:       se.Time,
		})
	}
	return nil
}
