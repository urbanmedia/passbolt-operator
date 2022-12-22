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
	"os"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	corev1 "k8s.io/api/core/v1"

	passboltv1alpha1 "github.com/urbanmedia/passbolt-operator/api/v1alpha1"
	"github.com/urbanmedia/passbolt-operator/pkg/passbolt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PassboltSecretReconciler reconciles a PassboltSecret object
type PassboltSecretReconciler struct {
	client.Client
	Scheme *runtime.Scheme
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

	secret := &passboltv1alpha1.PassboltSecret{}
	err := r.Client.Get(ctx, req.NamespacedName, secret)
	if err != nil {
		// If the resource no longer exists, in which case we stop processing.
		logr.Error(err, "unable to fetch PassboltSecret")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	// cleanup status
	secret.Status.SyncErrors = []passboltv1alpha1.SyncError{}

	// create context with timeout
	// TODO: make timeout configurable
	ctx2, cf := context.WithTimeout(context.Background(), 30*time.Second)
	defer cf()

	// create status update function
	updateStatus := func(ctx context.Context, passboltSecret *passboltv1alpha1.PassboltSecret) error {
		err := r.Client.Status().Update(ctx, passboltSecret)
		if err != nil {
			logr.Error(err, "unable to update PassboltSecret status")
			return err
		}
		return nil
	}

	// FIXME: pass real values to the passbolt client
	// create passbolt client
	clnt, err := passbolt.NewClient(ctx2, os.Getenv("PASSBOLT_URL"), os.Getenv("PASSBOLT_GPG"), os.Getenv("PASSBOLT_PASSWORD"))
	if err != nil {
		secret.Status.SyncStatus = passboltv1alpha1.SyncStatusError
		secret.Status.SyncErrors = append(secret.Status.SyncErrors, passboltv1alpha1.SyncError{
			Message: fmt.Sprintf("unable to create passbolt client: %s", err.Error()),
			Time:    metav1.Now(),
		})
		err2 := updateStatus(ctx2, secret)
		if err2 != nil {
			// if the status update fails, we do not want to return the error
			logr.Error(err2, "unable to update PassboltSecret status")
			return ctrl.Result{}, nil
		}
		logr.Error(err, "unable to create passbolt client")
		return ctrl.Result{}, err
	}
	defer clnt.Close(ctx2)

	err = clnt.LoadCache(ctx2)
	if err != nil {
		secret.Status.SyncStatus = passboltv1alpha1.SyncStatusError
		secret.Status.SyncErrors = append(secret.Status.SyncErrors, passboltv1alpha1.SyncError{
			Message: fmt.Sprintf("unable to load cache: %s", err.Error()),
			Time:    metav1.Now(),
		})
		err2 := updateStatus(ctx2, secret)
		if err2 != nil {
			// if the status update fails, we do not want to return the error
			logr.Error(err2, "unable to update PassboltSecret status")
			return ctrl.Result{}, nil
		}
		// if the cache load fails, we do not want to return the error
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
	for _, scrt := range secret.Spec.Secrets {
		secretData, err := clnt.GetSecret(ctx, scrt.Name)
		if err != nil {
			secret.Status.SyncStatus = passboltv1alpha1.SyncStatusError
			secret.Status.SyncErrors = append(secret.Status.SyncErrors, passboltv1alpha1.SyncError{
				Message:    fmt.Sprintf("unable to GET secret %s from passbolt: %s", scrt.Name, err.Error()),
				Time:       metav1.Now(),
				SecretName: scrt.Name,
				SecretKey:  scrt.KubernetesSecretKey,
			})
			err2 := updateStatus(ctx2, secret)
			if err2 != nil {
				// if the status update fails, we do not want to return the error
				logr.Error(err2, "unable to update PassboltSecret status")
				return ctrl.Result{}, nil
			}
			// if the secret cannot be retrieved, we do not want to return the error
			return ctrl.Result{}, nil
		}
		k8sSecret.StringData[scrt.KubernetesSecretKey] = secretData
	}

	// set owner reference if LeaveOnDelete was set to false
	if !secret.Spec.LeaveOnDelete {
		// set owner reference
		err = ctrl.SetControllerReference(secret, k8sSecret, r.Scheme)
		if err != nil {
			secret.Status.SyncStatus = passboltv1alpha1.SyncStatusError
			secret.Status.SyncErrors = append(secret.Status.SyncErrors, passboltv1alpha1.SyncError{
				Message: fmt.Sprintf("failed to set controller reference to secret %s.%s: %s", req.Name, req.Namespace, err.Error()),
				Time:    metav1.Now(),
			})
			err2 := updateStatus(ctx2, secret)
			if err2 != nil {
				// if the status update fails, we do not want to return the error
				logr.Error(err2, "unable to update PassboltSecret status")
				return ctrl.Result{}, nil
			}
			return ctrl.Result{}, nil
		}
	}

	// check if the secret already exists
	err = r.Client.Get(ctx, req.NamespacedName, &corev1.Secret{})
	if err != nil && !errors.IsNotFound(err) {
		secret.Status.SyncStatus = passboltv1alpha1.SyncStatusError
		secret.Status.SyncErrors = append(secret.Status.SyncErrors, passboltv1alpha1.SyncError{
			Message: fmt.Sprintf("unable to GET secret %q.%q from Kubernetes: %s", req.Name, req.Namespace, err.Error()),
			Time:    metav1.Now(),
		})
		err2 := updateStatus(ctx2, secret)
		if err2 != nil {
			// if the status update fails, we do not want to return the error
			logr.Error(err2, "unable to update PassboltSecret status")
			return ctrl.Result{}, nil
		}
		// if the secret cannot be retrieved, we do not want to return the error
		return ctrl.Result{}, nil
	}
	// check if the secret already exists
	if errors.IsNotFound(err) {
		err = r.Client.Create(ctx, k8sSecret)
		if err != nil {
			secret.Status.SyncStatus = passboltv1alpha1.SyncStatusError
			secret.Status.SyncErrors = append(secret.Status.SyncErrors, passboltv1alpha1.SyncError{
				Message: fmt.Sprintf("unable to CREATE secret %q.%q from Kubernetes: %s", k8sSecret.GetName(), k8sSecret.GetNamespace(), err.Error()),
				Time:    metav1.Now(),
			})
			err2 := updateStatus(ctx2, secret)
			if err2 != nil {
				// if the status update fails, we do not want to return the error
				logr.Error(err2, "unable to update PassboltSecret status")
				return ctrl.Result{}, nil
			}
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}
	// secret already exists, update it
	err = r.Client.Update(ctx, k8sSecret)
	if err != nil {
		secret.Status.SyncStatus = passboltv1alpha1.SyncStatusError
		secret.Status.SyncErrors = append(secret.Status.SyncErrors, passboltv1alpha1.SyncError{
			Message: fmt.Sprintf("unable to UPDATE secret %q.%q from Kubernetes: %s", k8sSecret.GetName(), k8sSecret.GetNamespace(), err.Error()),
			Time:    metav1.Now(),
		})
		err2 := updateStatus(ctx2, secret)
		if err2 != nil {
			// if the status update fails, we do not want to return the error
			logr.Error(err2, "unable to update PassboltSecret status")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	secret.Status.SyncStatus = passboltv1alpha1.SyncStatusSuccess
	secret.Status.LastSync = metav1.Now()
	err = r.Client.Status().Update(ctx, secret)
	if err != nil {
		logr.Error(err, "unable to update PassboltSecret status")
		// we don't return an error here, as the secret was successfully synced
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PassboltSecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&passboltv1alpha1.PassboltSecret{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}
