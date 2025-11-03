// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/MartinWeindel/gardener-extension-dummy-service/pkg/apis/config"
)

// ValidateConfiguration validates the passed configuration instance.
func ValidateConfiguration(config *config.Configuration) field.ErrorList {
	allErrs := field.ErrorList{}

	if config.Bar != nil && *config.Bar == "bad" {
		allErrs = append(allErrs, field.Invalid(field.NewPath("bar"), *config.Bar, "value 'bad' is not allowed"))
	}

	return allErrs
}
