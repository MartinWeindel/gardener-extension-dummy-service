// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"os"

	controllercmd "github.com/gardener/gardener/extensions/pkg/controller/cmd"
	heartbeatcmd "github.com/gardener/gardener/extensions/pkg/controller/heartbeat/cmd"

	dummyservicecmd "github.com/MartinWeindel/gardener-extension-dummy-service/pkg/cmd"
	"github.com/MartinWeindel/gardener-extension-dummy-service/pkg/constants"
)

// ExtensionName is the name of the extension.
const ExtensionName = "extension-dummy-service"

// Options holds configuration passed to the Dummy Service controller.
type Options struct {
	generalOptions                *controllercmd.GeneralOptions
	dummyOptions                  *dummyservicecmd.DummyServiceOptions
	restOptions                   *controllercmd.RESTOptions
	managerOptions                *controllercmd.ManagerOptions
	lifecycleControllerOptions    *controllercmd.ControllerOptions
	controlPlaneControllerOptions *controllercmd.ControllerOptions
	healthOptions                 *controllercmd.ControllerOptions
	heartbeatOptions              *heartbeatcmd.Options
	controllerSwitches            *controllercmd.SwitchOptions
	reconcileOptions              *controllercmd.ReconcilerOptions
	optionAggregator              controllercmd.OptionAggregator
}

// NewOptions creates a new Options instance.
func NewOptions() *Options {
	options := &Options{
		generalOptions: &controllercmd.GeneralOptions{},
		dummyOptions:   &dummyservicecmd.DummyServiceOptions{},
		restOptions:    &controllercmd.RESTOptions{},
		managerOptions: &controllercmd.ManagerOptions{
			// These are default values.
			LeaderElection:          true,
			LeaderElectionID:        controllercmd.LeaderElectionNameID(ExtensionName),
			LeaderElectionNamespace: os.Getenv(constants.EnvLeaderElectionNamespace),
		},
		lifecycleControllerOptions: &controllercmd.ControllerOptions{
			// This is a default value.
			MaxConcurrentReconciles: 5,
		},
		controlPlaneControllerOptions: &controllercmd.ControllerOptions{
			// This is a default value.
			MaxConcurrentReconciles: 1,
		},
		healthOptions: &controllercmd.ControllerOptions{
			// This is a default value.
			MaxConcurrentReconciles: 5,
		},
		heartbeatOptions: &heartbeatcmd.Options{
			// This is a default value.
			ExtensionName:        ExtensionName,
			RenewIntervalSeconds: 30,
			Namespace:            os.Getenv(constants.EnvLeaderElectionNamespace),
		},
		controllerSwitches: dummyservicecmd.ControllerSwitches(),
		reconcileOptions:   &controllercmd.ReconcilerOptions{},
	}

	options.optionAggregator = controllercmd.NewOptionAggregator(
		options.generalOptions,
		options.restOptions,
		options.managerOptions,
		options.lifecycleControllerOptions,
		controllercmd.PrefixOption("controlplane-", options.controlPlaneControllerOptions),
		options.dummyOptions,
		controllercmd.PrefixOption("healthcheck-", options.healthOptions),
		controllercmd.PrefixOption("heartbeat-", options.heartbeatOptions),
		options.controllerSwitches,
		options.reconcileOptions,
	)

	return options
}
