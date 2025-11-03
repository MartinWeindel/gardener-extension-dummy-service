// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"errors"
	"os"

	extensionsapisconfigv1alpha1 "github.com/gardener/gardener/extensions/pkg/apis/config/v1alpha1"
	"github.com/gardener/gardener/extensions/pkg/controller/cmd"
	extensionshealthcheckcontroller "github.com/gardener/gardener/extensions/pkg/controller/healthcheck"
	extensionsheartbeatcontroller "github.com/gardener/gardener/extensions/pkg/controller/heartbeat"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"

	"github.com/MartinWeindel/gardener-extension-dummy-service/pkg/apis/config"
	"github.com/MartinWeindel/gardener-extension-dummy-service/pkg/apis/config/v1alpha1"
	"github.com/MartinWeindel/gardener-extension-dummy-service/pkg/apis/config/validation"
	healthcheckcontroller "github.com/MartinWeindel/gardener-extension-dummy-service/pkg/controller/healthcheck"
	"github.com/MartinWeindel/gardener-extension-dummy-service/pkg/controller/lifecycle"
)

var (
	scheme  *runtime.Scheme
	decoder runtime.Decoder
)

func init() {
	scheme = runtime.NewScheme()
	utilruntime.Must(config.AddToScheme(scheme))
	utilruntime.Must(v1alpha1.AddToScheme(scheme))

	decoder = serializer.NewCodecFactory(scheme).UniversalDecoder()
}

// DummyServiceOptions holds options related to the certificate service.
type DummyServiceOptions struct {
	ConfigLocation string
	config         *DummyServiceConfig
}

// AddFlags implements Flagger.AddFlags.
func (o *DummyServiceOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.ConfigLocation, "config", "", "Path to dummy service configuration")
}

// Complete implements Completer.Complete.
func (o *DummyServiceOptions) Complete() error {
	if o.ConfigLocation == "" {
		return errors.New("config location is not set")
	}
	data, err := os.ReadFile(o.ConfigLocation)
	if err != nil {
		return err
	}

	config := config.Configuration{}
	_, _, err = decoder.Decode(data, nil, &config)
	if err != nil {
		return err
	}

	if errs := validation.ValidateConfiguration(&config); len(errs) > 0 {
		return errs.ToAggregate()
	}

	o.config = &DummyServiceConfig{
		config: config,
	}

	return nil
}

// Completed returns the decoded CertificatesServiceConfiguration instance. Only call this if `Complete` was successful.
func (o *DummyServiceOptions) Completed() *DummyServiceConfig {
	return o.config
}

// DummyServiceConfig contains configuration information about the certificate service.
type DummyServiceConfig struct {
	config config.Configuration
}

// Apply applies the DummyServiceOptions to the passed ControllerOptions instance.
func (c *DummyServiceConfig) Apply(config *config.Configuration) {
	*config = c.config
}

// ControllerSwitches are the cmd.SwitchOptions for the provider controllers.
func ControllerSwitches() *cmd.SwitchOptions {
	return cmd.NewSwitchOptions(
		cmd.Switch(lifecycle.ControllerName, lifecycle.AddToManager),
		cmd.Switch(extensionshealthcheckcontroller.ControllerName, healthcheckcontroller.AddToManager),
		cmd.Switch(extensionsheartbeatcontroller.ControllerName, extensionsheartbeatcontroller.AddToManager),
	)
}

func (c *DummyServiceConfig) ApplyHealthCheckConfig(config *extensionsapisconfigv1alpha1.HealthCheckConfig) {
	if c.config.HealthCheckConfig != nil {
		*config = *c.config.HealthCheckConfig
	}
}
