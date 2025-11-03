// SPDX-FileCopyrightText: 2025 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package garden

import (
	"context"
	"time"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	operatorv1alpha1 "github.com/gardener/gardener/pkg/apis/operator/v1alpha1"
	"github.com/gardener/gardener/pkg/client/kubernetes"
	operatorclient "github.com/gardener/gardener/pkg/operator/client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

var _ = Describe("Dummy-Service Tests", func() {
	var (
		operatorExtension = &operatorv1alpha1.Extension{ObjectMeta: metav1.ObjectMeta{Name: "extension-dummy-service"}}

		rawExtension = &runtime.RawExtension{
			Raw: []byte(`{
  "apiVersion": "service.dummy.extensions.gardener.cloud/v1alpha1",
  "kind": "DummyConfig",
  "foo": "foo-value"
}`),
		}
	)

	It("Create, Delete", Label("simple"), func() {
		ctx, cancel := context.WithTimeout(parentCtx, 15*time.Minute)
		defer cancel()

		By("Deploy Extension")
		Expect(execMake(ctx, "extension-up")).To(Succeed())

		By("Get Virtual Garden Client")
		gardenClientSet, err := kubernetes.NewClientFromSecret(ctx, runtimeClient, v1beta1constants.GardenNamespace, "gardener",
			kubernetes.WithDisabledCachedClient(),
			kubernetes.WithClientOptions(client.Options{Scheme: operatorclient.VirtualScheme}),
		)
		Expect(err).NotTo(HaveOccurred())

		By("Create workerless shoot")
		shoot := &gardencorev1beta1.Shoot{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "local-wl",
				Namespace: "garden-local",
			},
		}
		_, err = controllerutil.CreateOrUpdate(ctx, gardenClientSet.Client(), shoot, func() error {
			shoot.Spec.CloudProfile = &gardencorev1beta1.CloudProfileReference{
				Name: "local",
			}
			shoot.Spec.Region = "local"
			shoot.Spec.Provider = gardencorev1beta1.Provider{
				Type: "local",
			}
			shoot.Spec.Extensions = []gardencorev1beta1.Extension{
				{
					Type:           "dummy-service",
					ProviderConfig: rawExtension,
				},
			}
			return nil
		})
		Expect(err).NotTo(HaveOccurred())

		By("Wait for Shoot to be 'Ready'")
		waitForShootToBeReconciled(ctx, gardenClientSet.Client(), shoot)

		By("Check Operator Extension status")
		waitForOperatorExtensionToBeReconciled(ctx, operatorExtension)

		//By("Delete Extension")
		//Expect(runtimeClient.Delete(ctx, operatorExtension)).To(Succeed())
		//waitForOperatorExtensionToBeDeleted(ctx, operatorExtension)
	})
})
