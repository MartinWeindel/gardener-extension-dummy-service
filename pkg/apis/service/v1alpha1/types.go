// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DummyConfig configuration resource
type DummyConfig struct {
	metav1.TypeMeta `json:",inline"`

	// Foo is an example field of DummyConfig.
	// optional
	Foo *string `json:"foo,omitempty"`
}
