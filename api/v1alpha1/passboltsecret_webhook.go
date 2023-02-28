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

package v1alpha1

import (
	"github.com/urbanmedia/passbolt-operator/api/v1alpha2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// log is for logging in this package.
var passboltsecretlog = logf.Log.WithName("passboltsecret-resource")

// ConvertTo converts this CronJob to the Hub version (v1alpha2).
func (src *PassboltSecret) ConvertTo(dstRaw conversion.Hub) error {
	passboltsecretlog.V(100).Info("converting PassboltSecret v1alpha1 to v1alpha2")
	dst := dstRaw.(*v1alpha2.PassboltSecret)
	src.Spec.LeaveOnDelete = dst.Spec.LeaveOnDelete
	dst.Spec.Secrets = make([]v1alpha2.SecretSpec, len(src.Spec.Secrets))
	for i, s := range src.Spec.Secrets {
		dst.Spec.Secrets[i] = v1alpha2.SecretSpec{
			PassboltSecret: v1alpha2.PassboltSpec{
				Name:  s.PassboltSecret.Name,
				Field: v1alpha2.FieldName(s.PassboltSecret.Field),
			},
			KubernetesSecretKey: s.KubernetesSecretKey,
		}
	}
	return nil
}

// ConvertFrom converts from the Hub version (v1alpha2) to this version.
func (dst *PassboltSecret) ConvertFrom(srcRaw conversion.Hub) error {
	passboltsecretlog.V(100).Info("converting from PassboltSecret v1alpha2 to v1alpha1")
	src := srcRaw.(*v1alpha2.PassboltSecret)
	dst.Spec.LeaveOnDelete = src.Spec.LeaveOnDelete
	dst.Spec.Secrets = make([]SecretSpec, len(src.Spec.Secrets))
	for i, s := range src.Spec.Secrets {
		dst.Spec.Secrets[i] = SecretSpec{
			PassboltSecret: PassboltSpec{
				Name:  s.PassboltSecret.Name,
				Field: FieldName(s.PassboltSecret.Field),
			},
			KubernetesSecretKey: s.KubernetesSecretKey,
		}
	}
	return nil
}

func (r *PassboltSecret) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
