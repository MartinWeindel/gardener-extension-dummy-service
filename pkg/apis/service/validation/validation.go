// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	"github.com/gardener/gardener/extensions/pkg/controller"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/MartinWeindel/gardener-extension-dummy-service/pkg/apis/service"
)

// ValidateDummyConfig validates the passed configuration instance.
func ValidateDummyConfig(config *service.DummyConfig, cluster *controller.Cluster) field.ErrorList {
	allErrs := field.ErrorList{}

	root := field.NewPath("providerConfig")

	if config.Foo != nil && *config.Foo == "bad" {
		allErrs = append(allErrs, field.Invalid(root.Child("foo"), *config.Foo, "value 'bad' is not allowed"))
	}

	return allErrs
}
