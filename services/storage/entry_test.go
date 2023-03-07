package storage_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/rawfish-dev/angrypros-api/services/storage"
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

	Context("GetEntryById", func() {
		Context("success", func() {
			It("should return an entry given a valid id", func() {
				entry, err := testStorageService.GetEntryById(seedEntries[0].Id)
				Expect(err).ToNot(HaveOccurred())
				Expect(entry.Id).To(Equal(seedEntries[0].Id))
				Expect(entry.TextContent).To(Equal(seedEntries[0].TextContent))
				Expect(entry.RageLevel).To(Equal(seedEntries[0].RageLevel))
				Expect(entry.UserId).To(Equal(seedEntries[0].UserId))
				Expect(entry.User.Id).To(Equal(seedEntries[0].UserId))
				Expect(entry.CountryIsoAlpha2Code).To(Equal(seedEntries[0].CountryIsoAlpha2Code))
				Expect(entry.Country.IsoAlpha2Code).To(Equal(seedEntries[0].CountryIsoAlpha2Code))
				Expect(entry.AngerTierId).To(Equal(seedAngerTiers[0].Id))
				Expect(entry.AngerTier.Id).To(Equal(seedAngerTiers[0].Id))
				Expect(entry.CreatedAt).To(BeTemporally(">=", seedTimeNow))
				Expect(entry.UpdatedAt).To(BeTemporally(">=", seedTimeNow))
			})
		})

		Context("failure", func() {
			It("should return an error if the entry id is invalid", func() {
				entry, err := testStorageService.GetEntryById(88888)
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(storage.RecordNotFoundError{}))
				Expect(entry).To(BeNil())
			})
		})
	})

	XContext("EditEntry", func() {
		Context("success", func() {})
		Context("failure", func() {})
	})

	XContext("GetEntries", func() {
		Context("success", func() {})
	})
})
