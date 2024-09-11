// Copyright (c) 2024 Red Hat, Inc.

package pkg

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	testdata "github.com/rhecosystemappeng/patch-utils/pkg/testdata/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"testing"
)

var clt client.Client
var testEnv *envtest.Environment

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controller Tests")
}

var _ = BeforeSuite(func() {
	By("bootstrapping the testing environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("testdata", "crd")},
	}

	// install the scheme and initialize the test api
	scheme := runtime.NewScheme()
	Expect(corev1.AddToScheme(scheme)).To(Succeed())
	Expect(testdata.InitTestApi(scheme)).To(Succeed())

	// start testing environment and get config for the client
	cfg, err := testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	// create and save the test client
	clt, err = client.New(cfg, client.Options{Scheme: scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(clt).NotTo(BeNil())
})

var _ = AfterSuite(func() {
	By("tearing down the testing environment")
	Expect(testEnv.Stop()).To(Succeed())
})
