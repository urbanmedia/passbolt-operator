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

package controller

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	passboltv1alpha1 "github.com/urbanmedia/passbolt-operator/api/v1alpha1"
	passboltv1alpha2 "github.com/urbanmedia/passbolt-operator/api/v1alpha2"
	"github.com/urbanmedia/passbolt-operator/pkg/cache"
	"github.com/urbanmedia/passbolt-operator/pkg/passbolt"
	//+kubebuilder:scaffold:imports
)

var (
	passboltURL string = "http://localhost:8088"
)

// define global constants
const (
	passboltPassword = "TestTest123!"
	passboltGPGKey   = `-----BEGIN PGP PRIVATE KEY BLOCK-----

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
=jwOw
-----END PGP PRIVATE KEY BLOCK-----`
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	cfg            *rest.Config
	k8sClient      client.Client
	testEnv        *envtest.Environment
	passboltClient *passbolt.Client
	ctx            context.Context
	cancel         context.CancelFunc
)

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
	}

	var err error

	// initialize passbolt client
	ctx, cancel = context.WithCancel(context.TODO())
	clnt, err := passbolt.NewClient(ctx, cache.NewInMemoryCache(), passboltURL, passboltGPGKey, passboltPassword)
	Expect(err).NotTo(HaveOccurred())
	passboltClient = clnt
	// define timeout context for loading cache
	ctx2, cf := context.WithTimeout(ctx, 10*time.Second)
	defer cf()
	// fill cache
	err = passboltClient.LoadCache(ctx2)
	Expect(err).NotTo(HaveOccurred())

	// cfg is defined in this file globally.
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = passboltv1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = passboltv1alpha2.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	// setup manager and webhook server
	webhookInstallOptions := &testEnv.WebhookInstallOptions
	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:             scheme.Scheme,
		Host:               webhookInstallOptions.LocalServingHost,
		Port:               webhookInstallOptions.LocalServingPort,
		CertDir:            webhookInstallOptions.LocalServingCertDir,
		LeaderElection:     false,
		MetricsBindAddress: "0",
	})
	Expect(err).ToNot(HaveOccurred())

	err = (&PassboltSecretReconciler{
		Client:         k8sManager.GetClient(),
		Scheme:         k8sManager.GetScheme(),
		PassboltClient: passboltClient,
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	if err = (&passboltv1alpha1.PassboltSecret{}).SetupWebhookWithManager(k8sManager); err != nil {
		Expect(err).ToNot(HaveOccurred(), "unable to create webhook", "webhook", "PassboltSecret")

	}
	if err = (&passboltv1alpha2.PassboltSecret{}).SetupWebhookWithManager(k8sManager); err != nil {
		Expect(err).ToNot(HaveOccurred(), "unable to create webhook", "webhook", "PassboltSecret")
	}

	go func() {
		defer GinkgoRecover()
		err = k8sManager.Start(ctx)
		Expect(err).ToNot(HaveOccurred(), "failed to run manager")
	}()
})

var _ = AfterSuite(func() {
	cancel()
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})
