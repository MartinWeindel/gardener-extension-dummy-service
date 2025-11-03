// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package garden

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	operatorv1alpha1 "github.com/gardener/gardener/pkg/apis/operator/v1alpha1"
	"github.com/gardener/gardener/pkg/client/kubernetes"
	"github.com/gardener/gardener/pkg/logger"
	. "github.com/gardener/gardener/pkg/utils/test"
	. "github.com/gardener/gardener/pkg/utils/test/matchers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"k8s.io/apimachinery/pkg/runtime"
	kubernetesscheme "k8s.io/client-go/kubernetes/scheme"
	componentbaseconfigv1alpha1 "k8s.io/component-base/config/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	parentCtx     context.Context
	runtimeClient client.Client
)

var _ = BeforeSuite(func() {
	Expect(os.Getenv("KUBECONFIG")).NotTo(BeEmpty(), "KUBECONFIG must be set")
	Expect(os.Getenv("REPO_ROOT")).NotTo(BeEmpty(), "REPO_ROOT must be set")

	logf.SetLogger(logger.MustNewZapLogger(logger.InfoLevel, logger.FormatJSON, zap.WriteTo(GinkgoWriter)))

	restConfig, err := kubernetes.RESTConfigFromClientConnectionConfiguration(&componentbaseconfigv1alpha1.ClientConnectionConfiguration{Kubeconfig: os.Getenv("KUBECONFIG")}, nil, kubernetes.AuthTokenFile, kubernetes.AuthClientCertificate)
	Expect(err).NotTo(HaveOccurred())

	scheme := runtime.NewScheme()
	Expect(kubernetesscheme.AddToScheme(scheme)).To(Succeed())
	Expect(operatorv1alpha1.AddToScheme(scheme)).To(Succeed())
	Expect(extensionsv1alpha1.AddToScheme(scheme)).To(Succeed())
	runtimeClient, err = client.New(restConfig, client.Options{Scheme: scheme})
	Expect(err).NotTo(HaveOccurred())
})

var _ = BeforeEach(func() {
	parentCtx = context.Background()
})

func waitForShootToBeReconciled(ctx context.Context, gardenClient client.Client, shoot *gardencorev1beta1.Shoot) {
	CEventually(ctx, func(g Gomega) gardencorev1beta1.LastOperationState {
		g.Expect(gardenClient.Get(ctx, client.ObjectKeyFromObject(shoot), shoot)).To(Succeed())
		if shoot.Status.LastOperation == nil || shoot.Status.ObservedGeneration != shoot.Generation {
			return ""
		}
		return shoot.Status.LastOperation.State
	}).WithPolling(2 * time.Second).Should(Equal(gardencorev1beta1.LastOperationStateSucceeded))
}

func waitForOperatorExtensionToBeReconciled(ctx context.Context, extension *operatorv1alpha1.Extension) {
	CEventually(ctx, func(g Gomega) []gardencorev1beta1.Condition {
		g.Expect(runtimeClient.Get(ctx, client.ObjectKeyFromObject(extension), extension)).To(Succeed())
		if extension.Status.ObservedGeneration != extension.Generation {
			return nil
		}
		return extension.Status.Conditions
	}).WithPolling(1 * time.Second).Should(ContainElements(MatchFields(IgnoreExtras, Fields{
		"Type":   Equal(operatorv1alpha1.ExtensionInstalled),
		"Status": Equal(gardencorev1beta1.ConditionTrue),
	}), MatchFields(IgnoreExtras, Fields{
		"Type":   Equal(gardencorev1beta1.ConditionType(operatorv1alpha1.ControllerInstallationsHealthy)),
		"Status": Equal(gardencorev1beta1.ConditionTrue),
	})))
}

func waitForOperatorExtensionToBeDeleted(ctx context.Context, extension *operatorv1alpha1.Extension) {
	CEventually(ctx, func() error {
		return runtimeClient.Get(ctx, client.ObjectKeyFromObject(extension), extension)
	}).WithPolling(2 * time.Second).Should(BeNotFoundError())
}

// ExecMake executes one or multiple make targets.
func execMake(ctx context.Context, targets ...string) error {
	cmd := exec.CommandContext(ctx, "make", targets...)
	cmd.Dir = os.Getenv("REPO_ROOT")
	for _, key := range []string{"PATH", "GOPATH", "HOME", "KUBECONFIG"} {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, os.Getenv(key)))
	}
	cmdString := fmt.Sprintf("running make %s", strings.Join(targets, " "))
	logf.Log.Info(cmdString)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s failed: %s\n%s", cmdString, err, string(output))
	}
	return nil
}
