// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	extensionsconfigv1alpha1 "github.com/gardener/gardener/extensions/pkg/apis/config/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Configuration contains information about the dummy service configuration.
type Configuration struct {
	metav1.TypeMeta

	// Bar is an example field of Configuration.
	Bar *string

	// HealthCheckConfig is the config for the health check controller.
	HealthCheckConfig *extensionsconfigv1alpha1.HealthCheckConfig
}
