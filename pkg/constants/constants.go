// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package constants

const (
	// ExtensionType is the name of the extension type.
	ExtensionType = "dummy-service"
	// ServiceName is the name of the service.
	ServiceName = ExtensionType

	extensionServiceName = "extension-" + ServiceName
	// Origin is the origin field used in managed resources created by the extension.
	Origin = "gardener-extension-dummy-service"
	// ManagedResourceNamesShoot is the name used to describe the managed shoot resources.
	ManagedResourceNamesShoot = extensionServiceName + "-shoot"
	// ManagedResourceNamesSeed is the name used to describe the managed resources for the seed.
	ManagedResourceNamesSeed = extensionServiceName + "-seed"

	// EnvLeaderElectionNamespace is the environment variable name set in the deployment for providing the pod namespace.
	EnvLeaderElectionNamespace = "LEADER_ELECTION_NAMESPACE"
)
