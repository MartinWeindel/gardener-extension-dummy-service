// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package lifecycle

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/gardener/gardener/extensions/pkg/controller"
	"github.com/gardener/gardener/extensions/pkg/controller/extension"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	kubernetesclient "github.com/gardener/gardener/pkg/client/kubernetes"
	"github.com/gardener/gardener/pkg/extensions"
	"github.com/gardener/gardener/pkg/utils/managedresources"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/MartinWeindel/gardener-extension-dummy-service/pkg/apis/config"
	"github.com/MartinWeindel/gardener-extension-dummy-service/pkg/apis/service"
	"github.com/MartinWeindel/gardener-extension-dummy-service/pkg/apis/service/v1alpha1"
	"github.com/MartinWeindel/gardener-extension-dummy-service/pkg/constants"
)

const (
	ControllerName = "dummy-service-lifecycle"
)

// NewActuator returns an actuator responsible for Extension resources.
func NewActuator(mgr manager.Manager, serviceConfig config.Configuration, extensionClasses []extensionsv1alpha1.ExtensionClass) (extension.Actuator, error) {
	a := &actuator{
		client:           mgr.GetClient(),
		config:           mgr.GetConfig(),
		scheme:           mgr.GetScheme(),
		decoder:          serializer.NewCodecFactory(mgr.GetScheme(), serializer.EnableStrict).UniversalDecoder(),
		logger:           log.Log.WithName(ControllerName),
		serviceConfig:    serviceConfig,
		extensionClasses: extensionClasses,
	}

	return a, nil
}

type actuator struct {
	client           client.Client
	config           *rest.Config
	decoder          runtime.Decoder
	extensionClasses []extensionsv1alpha1.ExtensionClass
	serviceConfig    config.Configuration
	logger           logr.Logger
	scheme           *runtime.Scheme
}

// Reconcile the Extension resource.
func (a *actuator) Reconcile(ctx context.Context, _ logr.Logger, ex *extensionsv1alpha1.Extension) error {
	var (
		namespace = ex.GetNamespace()
		cluster   *extensions.Cluster
		err       error
	)

	cluster, err = controller.GetCluster(ctx, a.client, namespace)
	if err != nil {
		return fmt.Errorf("failed to get cluster config: %w", err)
	}

	shootConfig := &v1alpha1.DummyConfig{}
	if ex.Spec.ProviderConfig != nil {
		if _, _, err := a.decoder.Decode(ex.Spec.ProviderConfig.Raw, nil, shootConfig); err != nil {
			return fmt.Errorf("failed to decode provider config: %w", err)
		}
	}

	internalShootConfig := &service.DummyConfig{}
	if err := a.scheme.Convert(shootConfig, internalShootConfig, nil); err != nil {
		return fmt.Errorf("failed to convert shoot config: %w", err)
	}

	if err := ValidateProviderConfig(internalShootConfig, cluster); err != nil {
		return fmt.Errorf("failed to validate provider config: %w", err)
	}

	dummyData := map[string]string{
		"foo": ptr.Deref(internalShootConfig.Foo, "<unset>"),
		"bar": ptr.Deref(a.serviceConfig.Bar, "<unset>"),
	}

	shootResources, err := getShootResources(namespace, dummyData)
	if err != nil {
		return err
	}

	if err := managedresources.CreateForShoot(ctx, a.client, namespace, constants.ManagedResourceNamesShoot, constants.Origin, false, shootResources); err != nil {
		return err
	}

	seedResources, err := getSeedResources(namespace, dummyData)
	if err != nil {
		return err
	}

	return managedresources.CreateForSeed(ctx, a.client, namespace, constants.ManagedResourceNamesSeed, false, seedResources)
}

// Delete the Extension resource.
func (a *actuator) Delete(ctx context.Context, _ logr.Logger, ex *extensionsv1alpha1.Extension) error {
	namespace := ex.GetNamespace()
	twoMinutes := 2 * time.Minute

	timeoutShootCtx, cancelShootCtx := context.WithTimeout(ctx, twoMinutes)
	defer cancelShootCtx()

	if err := managedresources.DeleteForSeed(ctx, a.client, namespace, constants.ManagedResourceNamesSeed); err != nil {
		return err
	}

	if err := managedresources.WaitUntilDeleted(timeoutShootCtx, a.client, namespace, constants.ManagedResourceNamesSeed); err != nil {
		return err
	}

	if err := managedresources.DeleteForShoot(ctx, a.client, namespace, constants.ManagedResourceNamesShoot); err != nil {
		return err
	}

	if err := managedresources.WaitUntilDeleted(timeoutShootCtx, a.client, namespace, constants.ManagedResourceNamesShoot); err != nil {
		return err
	}

	return nil
}

// ForceDelete implements Network.Actuator.
func (a *actuator) ForceDelete(ctx context.Context, log logr.Logger, ex *extensionsv1alpha1.Extension) error {
	return a.Delete(ctx, log, ex)
}

// Restore the Extension resource.
func (a *actuator) Restore(ctx context.Context, log logr.Logger, ex *extensionsv1alpha1.Extension) error {
	return a.Reconcile(ctx, log, ex)
}

// Migrate the Extension resource.
func (a *actuator) Migrate(ctx context.Context, log logr.Logger, ex *extensionsv1alpha1.Extension) error {
	return nil
}

func getShootResources(namespace string, data map[string]string) (map[string][]byte, error) {
	shootRegistry := managedresources.NewRegistry(kubernetesclient.ShootScheme, kubernetesclient.ShootCodec, kubernetesclient.ShootSerializer)

	var objects []client.Object
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "dummy-shoot-config",
			Namespace: namespace,
		},
		Data: data,
	}
	objects = append(objects, configMap)

	shootResources, err := shootRegistry.AddAllAndSerialize(objects...)
	if err != nil {
		return nil, err
	}
	return shootResources, nil
}

func getSeedResources(namespace string, data map[string]string) (map[string][]byte, error) {
	seedRegistry := managedresources.NewRegistry(kubernetesclient.SeedScheme, kubernetesclient.SeedCodec, kubernetesclient.SeedSerializer)

	var objects []client.Object
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "dummy-seed-config",
			Namespace: namespace,
		},
		Data: data,
	}
	objects = append(objects, configMap)

	seedResources, err := seedRegistry.AddAllAndSerialize(objects...)
	if err != nil {
		return nil, err
	}
	return seedResources, nil
}
