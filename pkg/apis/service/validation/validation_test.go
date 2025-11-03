// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package validation_test

import (
	"github.com/gardener/gardener/extensions/pkg/controller"
	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	gomegatypes "github.com/onsi/gomega/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"

	"github.com/MartinWeindel/gardener-extension-dummy-service/pkg/apis/service"
	"github.com/MartinWeindel/gardener-extension-dummy-service/pkg/apis/service/validation"
)

var _ = Describe("Validation", func() {
	var (
		cluster = &controller.Cluster{
			Shoot: &gardencorev1beta1.Shoot{},
		}
	)
	DescribeTable("#ValidateDummyConfigShoot",
		func(config service.DummyConfig, match gomegatypes.GomegaMatcher) {
			err := validation.ValidateDummyConfig(&config, cluster)
			Expect(err).To(match)
		},
		Entry("Empty", service.DummyConfig{}, BeEmpty()),
		Entry("Good configuration", service.DummyConfig{
			Foo: ptr.To("good"),
		}, BeEmpty()),
		Entry("Bad configuration", service.DummyConfig{
			Foo: ptr.To("bad"),
		}, ConsistOf(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Type":  Equal(field.ErrorTypeInvalid),
				"Field": Equal("providerConfig.foo"),
			})),
		)),
	)
})
