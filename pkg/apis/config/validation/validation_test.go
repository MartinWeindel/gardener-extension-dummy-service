// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package validation_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	gomegatypes "github.com/onsi/gomega/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"

	"github.com/MartinWeindel/gardener-extension-dummy-service/pkg/apis/config"
	"github.com/MartinWeindel/gardener-extension-dummy-service/pkg/apis/config/validation"
)

var _ = Describe("Validation", func() {

	DescribeTable("#ValidateConfiguration",
		func(config config.Configuration, match gomegatypes.GomegaMatcher) {
			err := validation.ValidateConfiguration(&config)
			Expect(err).To(match)
		},
		Entry("Empty configuration", config.Configuration{}, BeEmpty()),
		Entry("Good configuration", config.Configuration{
			Bar: ptr.To("good"),
		}, BeEmpty()),
		Entry("Bad configuration", config.Configuration{
			Bar: ptr.To("bad"),
		}, ConsistOf(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Type":   Equal(field.ErrorTypeInvalid),
				"Field":  Equal("bar"),
				"Detail": Equal("value 'bad' is not allowed"),
			})))),
	)
})
