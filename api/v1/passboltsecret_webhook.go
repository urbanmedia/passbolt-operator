/*
Copyright 2024 Verlag der Tagesspiegel GmbH.

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

package v1

import (
	"errors"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var (
	ErrInvalidSecretType              = errors.New("invalid secret type")
	ErrPassboltSecretNameIsRequired   = errors.New("passboltSecretName is required for secret type")
	ErrSecretsAreNotAllowed           = errors.New("secrets are not allowed")
	ErrFieldAndValueAreNotAllowed     = errors.New("field and value are not allowed")
	ErrFieldOrValueIsRequired         = errors.New("field or value is required")
	ErrSecretsAreRequired             = errors.New("secrets are required")
	ErrPassboltSecretNameIsNotAllowed = errors.New("passboltSecretName is not allowed")
)

// log is for logging in this package.
var passboltsecretlog = logf.Log.WithName("passboltsecret-resource")

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *PassboltSecret) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-passbolt-tagesspiegel-de-v1-passboltsecret,mutating=true,failurePolicy=fail,sideEffects=None,groups=passbolt.tagesspiegel.de,resources=passboltsecrets,verbs=create;update,versions=v1,name=mpassboltsecret.tagesspiegel.de,admissionReviewVersions=v1

var _ webhook.Defaulter = &PassboltSecret{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *PassboltSecret) Default() {
	passboltsecretlog.Info("default", "name", r.Name)
	if r.Spec.SecretType == "" || (r.Spec.SecretType != corev1.SecretTypeOpaque && r.Spec.SecretType != corev1.SecretTypeDockerConfigJson) {
		r.Spec.SecretType = corev1.SecretTypeOpaque
	}
}

//+kubebuilder:webhook:path=/validate-passbolt-tagesspiegel-de-v1-passboltsecret,mutating=false,failurePolicy=fail,sideEffects=None,groups=passbolt.tagesspiegel.de,resources=passboltsecrets,verbs=create;update,versions=v1,name=vpassboltsecret.tagesspiegel.de,admissionReviewVersions=v1

var _ webhook.Validator = &PassboltSecret{}

func (r *PassboltSecret) validatePassboltSecret() error {
	switch r.Spec.SecretType {
	case corev1.SecretTypeOpaque:
		if r.Spec.PassboltSecretID != nil {
			return fmt.Errorf("%w for secret %s.%s type %s", ErrPassboltSecretNameIsNotAllowed, r.GetName(), r.GetNamespace(), r.Spec.SecretType)
		}
		if len(r.Spec.PassboltSecrets) == 0 {
			return fmt.Errorf("%w for secret %s.%s type %s", ErrSecretsAreRequired, r.GetName(), r.GetNamespace(), r.Spec.SecretType)
		}
		// check if only FieldName or Value is set
		for _, secret := range r.Spec.PassboltSecrets {
			if secret.Field == "" && secret.Value == nil {
				return fmt.Errorf("%w for secret %s.%s and field %v", ErrFieldOrValueIsRequired, r.GetName(), r.GetNamespace(), secret)
			}
			if secret.Field != "" && secret.Value != nil {
				return fmt.Errorf("%w for secret %s.%s and field %v", ErrFieldAndValueAreNotAllowed, r.GetName(), r.GetNamespace(), secret)
			}
		}
		return nil
	case corev1.SecretTypeDockerConfigJson:
		if r.Spec.PassboltSecretID == nil {
			return fmt.Errorf("%w for secret %s.%s: %s", ErrPassboltSecretNameIsRequired, r.GetName(), r.GetNamespace(), r.Spec.SecretType)
		}
		if *r.Spec.PassboltSecretID == "" {
			return fmt.Errorf("%w for secret %s.%s: %s", ErrPassboltSecretNameIsRequired, r.GetName(), r.GetNamespace(), r.Spec.SecretType)
		}
		if len(r.Spec.PassboltSecrets) > 0 {
			return fmt.Errorf("%w for secret %s.%s type %s", ErrSecretsAreNotAllowed, r.GetName(), r.GetNamespace(), r.Spec.SecretType)
		}
		return nil
	default:
		return fmt.Errorf("%w %s.%s: %s", ErrInvalidSecretType, r.GetName(), r.GetNamespace(), r.Spec.SecretType)
	}
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *PassboltSecret) ValidateCreate() (admission.Warnings, error) {
	passboltsecretlog.Info("validate create", "name", r.Name)
	if err := r.validatePassboltSecret(); err != nil {
		return nil, err
	}
	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *PassboltSecret) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	passboltsecretlog.Info("validate update", "name", r.Name)
	if err := r.validatePassboltSecret(); err != nil {
		return nil, err
	}
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *PassboltSecret) ValidateDelete() (admission.Warnings, error) {
	passboltsecretlog.Info("validate delete", "name", r.Name)
	if err := r.validatePassboltSecret(); err != nil {
		return nil, err
	}
	return nil, nil
}
