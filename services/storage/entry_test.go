package storage_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Entry", func() {
	Context("GetAllAngerTiers", func() {
		Context("success", func() {
			It("should return all anger tiers sorted by rage level ascending", func() {
				angerTiers, err := testStorageService.GetAllAngerTiers()
				Expect(err).ToNot(HaveOccurred())
				Expect(angerTiers).To(HaveLen(3))
				Expect(angerTiers[0]).To(Equal(seedAngerTiers[1]))
				Expect(angerTiers[1]).To(Equal(seedAngerTiers[2]))
				Expect(angerTiers[2]).To(Equal(seedAngerTiers[0]))
			})
		})
	})

	XContext("CreateEntry", func() {
		Context("success", func() {})
		Context("failure", func() {})
	})

	XContext("GetEntryById", func() {
		Context("success", func() {})
		Context("failure", func() {})
	})

	XContext("EditEntry", func() {
		Context("success", func() {})
		Context("failure", func() {})
	})

	XContext("GetEntries", func() {
		Context("success", func() {})
	})
})
