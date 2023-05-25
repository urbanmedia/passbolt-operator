/*
Copyright 2022.

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

package main

import (
	"context"
	"flag"
	"os"
	"time"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	passboltv1alpha1 "github.com/urbanmedia/passbolt-operator/api/v1alpha1"
	passboltv1alpha2 "github.com/urbanmedia/passbolt-operator/api/v1alpha2"
	"github.com/urbanmedia/passbolt-operator/controllers"
	"github.com/urbanmedia/passbolt-operator/pkg/passbolt"
	//+kubebuilder:scaffold:imports
)

const (
	// loadCacheErrorRetryCounterReLogin defines how often we retry to re-login after a failed login attempt
	// must always be less than loadCacheErrorRetryCounterExit
	loadCacheErrorRetryCounterReLogin = 2
	// loadCacheErrorRetryCounterExit defines how often we retry to re-login after which we exit the program
	// must always be greater than loadCacheErrorRetryCounterReLogin
	loadCacheErrorRetryCounterExit = 3
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
	cacheLog = ctrl.Log.WithName("cache")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(passboltv1alpha1.AddToScheme(scheme))
	utilruntime.Must(passboltv1alpha2.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "19193eec.tagesspiegel.de",
		// LeaderElectionReleaseOnCancel defines if the leader should step down voluntarily
		// when the Manager ends. This requires the binary to immediately end when the
		// Manager is stopped, otherwise, this setting is unsafe. Setting this significantly
		// speeds up voluntary leader transitions as the new leader don't have to wait
		// LeaseDuration time first.
		//
		// In the default scaffold provided, the program ends immediately after
		// the manager stops, so would be fine to enable this option. However,
		// if you are doing or is intended to do any operation such as perform cleanups
		// after the manager stops then its usage might be unsafe.
		// LeaderElectionReleaseOnCancel: true,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// create context with timeout for passbolt client
	ctx, cf := context.WithTimeout(context.Background(), 5*time.Second)
	defer cf()

	// create passbolt client
	clnt, err := passbolt.NewClient(ctx, os.Getenv("PASSBOLT_URL"), os.Getenv("PASSBOLT_GPG"), os.Getenv("PASSBOLT_PASSWORD"))
	if err != nil {
		setupLog.Error(err, "unable to create passbolt client")
		os.Exit(1)
	}

	// initial cache load
	err = clnt.LoadCache(ctx)
	if err != nil {
		setupLog.Error(err, "unable to load passbolt cache")
		os.Exit(1)
	}

	// fill passbolt cache with existing secrets
	// ticker is set to 5 minutes
	ticker := time.NewTicker(5 * time.Minute)
	done := make(chan bool)
	defer func() {
		ticker.Stop()
		done <- true
	}()
	go func() {
		// define an error counter to prevent the operator from running in a broken state
		errorCounter := 0
		for {
			select {
			case <-done:
				// we exit here, because the ticker is stopped
				return
			case <-ticker.C:
				cacheLog.Info("loading passbolt cache")
				ctx, cf := context.WithTimeout(context.Background(), 5*time.Second)
				defer cf()
				err := clnt.LoadCache(ctx)
				if err != nil {
					cacheLog.Error(err, "unable to update local passbolt cache", "errorCounter", errorCounter)
					if errorCounter > loadCacheErrorRetryCounterReLogin {
						if errorCounter >= loadCacheErrorRetryCounterExit {
							// we exit here, because we failed to load the cache 3 times in a row
							// this is to prevent the operator from running in a broken state
							os.Exit(1)
						}
						// try to re-login to passbolt
						err := clnt.ReLogin(ctx)
						if err != nil {
							cacheLog.Error(err, "unable to re-login to passbolt")
							// we tried to re-login, but failed
							os.Exit(1)
						}
					}
					errorCounter++
					// we don't exit here, because we want to continue with the next tick
					continue
				}
				errorCounter = 0
			}
		}
	}()

	if err = (&controllers.PassboltSecretReconciler{
		Client:         mgr.GetClient(),
		Scheme:         mgr.GetScheme(),
		PassboltClient: clnt,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "PassboltSecret")
		os.Exit(1)
	}
	if os.Getenv("ENABLE_WEBHOOKS") != "false" {
		if err = (&passboltv1alpha1.PassboltSecret{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "PassboltSecret")
			os.Exit(1)
		}
		if err = (&passboltv1alpha2.PassboltSecret{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "PassboltSecret")
			os.Exit(1)
		}
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
