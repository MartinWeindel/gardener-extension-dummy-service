// SPDX-FileCopyrightText: 2025 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package lifecycle

import (
	"fmt"

	"github.com/gardener/gardener/extensions/pkg/controller"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/MartinWeindel/gardener-extension-dummy-service/pkg/apis/service"
	"github.com/MartinWeindel/gardener-extension-dummy-service/pkg/apis/service/v1alpha1"
	"github.com/MartinWeindel/gardener-extension-dummy-service/pkg/apis/service/validation"
)

// ValidateProviderConfig validates the given provider configuration.
func ValidateProviderConfig(cfg *service.DummyConfig, cluster *controller.Cluster) error {
	externalConfig := v1alpha1.DummyConfig{}
	if err := v1alpha1.Convert_service_DummyConfig_To_v1alpha1_DummyConfig(cfg, &externalConfig, nil); err != nil {
		return fmt.Errorf("failed to convert configuration: %w", err)
	}
	allErrs := validation.ValidateDummyConfig(cfg, cluster)
	if len(allErrs) > 0 {
		return fmt.Errorf("invalid network filter configuration: %s", field.ErrorList(allErrs).ToAggregate().Error())
	}

	return nil
}
