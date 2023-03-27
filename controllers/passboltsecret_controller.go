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

package controllers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"text/template"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	corev1 "k8s.io/api/core/v1"

	"github.com/Masterminds/sprig/v3"
	passboltv1alpha2 "github.com/urbanmedia/passbolt-operator/api/v1alpha2"
	"github.com/urbanmedia/passbolt-operator/pkg/passbolt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PassboltSecretReconciler reconciles a PassboltSecret object
type PassboltSecretReconciler struct {
	client.Client
	Scheme         *runtime.Scheme
	PassboltClient *passbolt.Client
}

//+kubebuilder:rbac:groups=passbolt.tagesspiegel.de,resources=passboltsecrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=passbolt.tagesspiegel.de,resources=passboltsecrets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=passbolt.tagesspiegel.de,resources=passboltsecrets/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;create;update;delete;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PassboltSecret object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *PassboltSecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logr := log.FromContext(ctx)

	logr.Info("reconciling PassboltSecret", "name", req.NamespacedName)

	// get passbolt secret resource from Kubernetes
	secret := &passboltv1alpha2.PassboltSecret{}
	err := r.Client.Get(ctx, req.NamespacedName, secret)
	if err != nil {
		if err = client.IgnoreNotFound(err); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}
	// cleanup status
	secret.Status.SyncErrors = []passboltv1alpha2.SyncError{}

	// create status update function
	updateStatus := func(ctx context.Context, passboltSecret *passboltv1alpha2.PassboltSecret) error {
		err := r.Client.Status().Update(ctx, passboltSecret)
		if err != nil {
			logr.Error(err, "unable to update PassboltSecret status")
			return err
		}
		return nil
	}

	// make sure that the secret type is supported
	if secret.Spec.SecretType != corev1.SecretTypeOpaque && secret.Spec.SecretType != corev1.SecretTypeDockerConfigJson {
		logr.Info("unsupported secret type", "type", secret.Spec.SecretType)
		secret.Status.SyncStatus = passboltv1alpha2.SyncStatusError
		secret.Status.SyncErrors = append(secret.Status.SyncErrors, passboltv1alpha2.SyncError{
			Message: fmt.Sprintf("unsupported secret type %q", secret.Spec.SecretType),
			Time:    metav1.Now(),
		})
		if err := updateStatus(ctx, secret); err != nil {
			logr.Error(err, "unable to update status")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// retrieve secrets from passbolt and store them in Kubernetes secrets
	k8sSecret := &corev1.Secret{
		ObjectMeta: ctrl.ObjectMeta{
			Name:        secret.Name,
			Namespace:   secret.Namespace,
			Labels:      secret.Labels,
			Annotations: secret.Annotations,
		},
		StringData: map[string]string{},
	}

	if secret.Spec.SecretType == corev1.SecretTypeDockerConfigJson {
		logr.Info("updating docker config json secret")
		// get secret from passbolt
		secretData, err := r.PassboltClient.GetSecret(ctx, *secret.Spec.PassboltSecretName, "")
		if err != nil {
			secret.Status.SyncStatus = passboltv1alpha2.SyncStatusError
			secret.Status.SyncErrors = append(secret.Status.SyncErrors, passboltv1alpha2.SyncError{
				Message:    fmt.Sprintf("unable to GET secret %q.%q from passbolt: %s", scrt.PassboltSecret.Name, scrt.PassboltSecret.Field, err.Error()),
				Time:       metav1.Now(),
				SecretName: scrt.PassboltSecret.Name,
				SecretKey:  scrt.KubernetesSecretKey,
			})
			if err := updateStatus(ctx, secret); err != nil {
				logr.Error(err, "unable to update status")
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}

		// create docker config json
		myConfig := map[string]any{
			"auths": map[string]any{
				secretData.URI: map[string]string{
					"auth": base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", secretData.Username, secretData.Password))),
				},
			},
		}
		// parse map to json
		bts, err := json.Marshal(myConfig)
		if err != nil {
			secret.Status.SyncStatus = passboltv1alpha2.SyncStatusError
			secret.Status.SyncErrors = append(secret.Status.SyncErrors, passboltv1alpha2.SyncError{
				Message: fmt.Sprintf("unable to marshal docker config json: %s", err),
				Time:    metav1.Now(),
			})
			if err := updateStatus(ctx, secret); err != nil {
				logr.Error(err, "unable to update status")
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}
		// add docker config json to k8s secret
		k8sSecret.StringData[corev1.DockerConfigJsonKey] = string(bts)
		k8sSecret.Type = corev1.SecretTypeDockerConfigJson
	} else {
		logr.Info("updating opaque secret")
		// add fields to K8s secret
		for _, scrt := range secret.Spec.Secrets {
			secretData, err := r.PassboltClient.GetSecret(ctx, scrt.PassboltSecret.Name, scrt.PassboltSecret.Field)
			if err != nil {
				secret.Status.SyncStatus = passboltv1alpha2.SyncStatusError
				secret.Status.SyncErrors = append(secret.Status.SyncErrors, passboltv1alpha2.SyncError{
					Message:    fmt.Sprintf("unable to GET secret %q.%q from passbolt: %s", scrt.PassboltSecret.Name, scrt.PassboltSecret.Field, err.Error()),
					Time:       metav1.Now(),
					SecretName: scrt.PassboltSecret.Name,
					SecretKey:  scrt.KubernetesSecretKey,
				})
				err2 := updateStatus(ctx, secret)
				if err2 != nil {
					// if the status update fails, we do not want to return the error
					logr.Error(err2, "unable to update PassboltSecret status")
					return ctrl.Result{}, nil
				}
				// if the secret cannot be retrieved, we do not want to return the error
				return ctrl.Result{}, nil
			}
			// if the field field is set, we expect a field to be defined
			if scrt.PassboltSecret.Field != "" {
				k8sSecret.StringData[scrt.KubernetesSecretKey] = secretData.FieldValue(scrt.PassboltSecret.Field)
				continue
			}
			// if the value field is set, we expect a template to be defined
			if scrt.PassboltSecret.Value != nil {
				tmpl, err := template.New("value").Funcs(sprig.FuncMap()).Parse(*scrt.PassboltSecret.Value)
				if err != nil {
					logr.Error(err, "unable to parse template")
					return ctrl.Result{}, nil
				}

				target := bytes.NewBuffer([]byte{})
				err = tmpl.Execute(target, *secretData)
				if err != nil {
					panic(err)
				}
				k8sSecret.StringData[scrt.KubernetesSecretKey] = target.String()
				continue
			}
			logr.Info("no field or value defined for secret", "name", scrt.PassboltSecret.Name)
		}
	}

	// set owner reference if LeaveOnDelete was set to false
	if !secret.Spec.LeaveOnDelete {
		logr.Info("leave on delete is false, setting owner reference")
		// set owner reference
		err = ctrl.SetControllerReference(secret, k8sSecret, r.Scheme)
		if err != nil {
			secret.Status.SyncStatus = passboltv1alpha2.SyncStatusError
			secret.Status.SyncErrors = append(secret.Status.SyncErrors, passboltv1alpha2.SyncError{
				Message: fmt.Sprintf("failed to set controller reference to secret %s.%s: %s", req.Name, req.Namespace, err.Error()),
				Time:    metav1.Now(),
			})
			err2 := updateStatus(ctx, secret)
			if err2 != nil {
				// if the status update fails, we do not want to return the error
				logr.Error(err2, "unable to update PassboltSecret status")
				return ctrl.Result{}, nil
			}
			return ctrl.Result{}, nil
		}
	}

	logr.Info("creating or updating secret")
	opRslt, err := controllerutil.CreateOrUpdate(ctx, r.Client, k8sSecret, func() error {
		// we don't need to update the secret, as the secret is defined already above
		return nil
	})
	if err != nil {
		secret.Status.SyncStatus = passboltv1alpha2.SyncStatusError
		secret.Status.SyncErrors = append(secret.Status.SyncErrors, passboltv1alpha2.SyncError{
			Message: fmt.Sprintf("unable to create or patch secret %q.%q: %s", req.Name, req.Namespace, err.Error()),
			Time:    metav1.Now(),
		})
		err2 := updateStatus(ctx, secret)
		if err2 != nil {
			// if the status update fails, we do not want to return the error
			logr.Error(err2, "unable to update PassboltSecret status")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, nil
	}

	if opRslt == controllerutil.OperationResultNone {
		// secret was not changed
		logr.V(10).Info("secret unchanged")
		return ctrl.Result{}, nil
	}

	// update status
	secret.Status.SyncStatus = passboltv1alpha2.SyncStatusSuccess
	secret.Status.LastSync = metav1.Now()
	err = r.Client.Status().Update(ctx, secret)
	if err != nil {
		logr.Error(err, "unable to update PassboltSecret status")
		// we don't return an error here, as the secret was successfully synced
		return ctrl.Result{}, nil
	}

	logr.Info("reconcile complete", "name", req.NamespacedName)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PassboltSecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&passboltv1alpha2.PassboltSecret{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}
