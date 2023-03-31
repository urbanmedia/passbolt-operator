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
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	corev1 "k8s.io/api/core/v1"

	passboltv1alpha2 "github.com/urbanmedia/passbolt-operator/api/v1alpha2"
	"github.com/urbanmedia/passbolt-operator/pkg/passbolt"
	"github.com/urbanmedia/passbolt-operator/pkg/util"
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

	// define Kubernetes secret to be created or updated
	k8sSecret := &corev1.Secret{
		ObjectMeta: ctrl.ObjectMeta{
			Name:        secret.Name,
			Namespace:   secret.Namespace,
			Labels:      secret.Labels,
			Annotations: secret.Annotations,
		},
		Type: secret.Spec.SecretType,
	}

	opRslt, err := controllerutil.CreateOrUpdate(ctx, r.Client, k8sSecret, util.UpdateSecret(ctx, r.PassboltClient, r.Scheme, secret, k8sSecret))
	if err != nil {
		if snErr, ok := err.(passboltv1alpha2.SyncError); ok {
			secret.Status.SyncStatus = passboltv1alpha2.SyncStatusError
			secret.Status.SyncErrors = append(secret.Status.SyncErrors, snErr)
			if err := updateStatus(ctx, secret); err != nil {
				logr.Error(err, "unable to update status")
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}
		logr.Error(err, "unable to create or update secret")
		return ctrl.Result{}, err
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
