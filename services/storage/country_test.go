package storage_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Country", func() {
	Context("GetAllCountries", func() {
		Context("success", func() {
			It("should return all countries sorted by name alphabetically", func() {
				countries, err := testStorageService.GetAllCountries()
				Expect(err).ToNot(HaveOccurred())
				Expect(countries).To(HaveLen(3))
				Expect(countries[0]).To(Equal(seedCountries[2]))
				Expect(countries[1]).To(Equal(seedCountries[0]))
				Expect(countries[2]).To(Equal(seedCountries[1]))
			})
		})
	})
})
