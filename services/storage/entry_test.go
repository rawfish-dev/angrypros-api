package storage_test

import (
	"time"

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

	Context("CreateEntry", func() {
		Context("success", func() {
			It("should create an entry given valid values", func() {
				entry, err := testStorageService.CreateEntry(seedUsers[0].Id, seedAngerTiers[0].Id, seedCountries[0].IsoAlpha2Code, "some valid content")
				Expect(err).ToNot(HaveOccurred())
				Expect(entry.Id).ToNot(BeZero())
				Expect(entry.TextContent).To(Equal("some valid content"))
				Expect(entry.UserId).To(Equal(seedEntries[0].UserId))
				Expect(entry.User.Id).To(Equal(seedEntries[0].UserId))
				Expect(entry.User.Country.IsoAlpha2Code).To(Equal(seedEntries[0].CountryIsoAlpha2Code))
				Expect(entry.CountryIsoAlpha2Code).To(Equal(seedCountries[0].IsoAlpha2Code))
				Expect(entry.Country.IsoAlpha2Code).To(Equal(seedCountries[0].IsoAlpha2Code))
				Expect(entry.AngerTierId).To(Equal(seedAngerTiers[0].Id))
				Expect(entry.AngerTier.Id).To(Equal(seedAngerTiers[0].Id))
				Expect(entry.CreatedAt).To(BeTemporally(">=", seedTimeNow))
				Expect(entry.UpdatedAt).To(BeTemporally(">=", seedTimeNow))
			})
		})

		Context("failure", func() {
			It("should return an error if the user id is invalid", func() {
				entry, err := testStorageService.CreateEntry(88888, seedAngerTiers[0].Id, seedCountries[0].IsoAlpha2Code, "some valid content")
				Expect(err).ToNot(BeNil())
				Expect(err).To(BeAssignableToTypeOf(storage.UserIdInvalidError{}))
				Expect(entry).To(BeNil())
			})

			It("should return an error if the anger tier id is invalid", func() {
				entry, err := testStorageService.CreateEntry(seedUsers[0].Id, 88888, seedCountries[0].IsoAlpha2Code, "some valid content")
				Expect(err).ToNot(BeNil())
				Expect(err).To(BeAssignableToTypeOf(storage.AngerTierIdInvalidError{}))
				Expect(entry).To(BeNil())
			})

			It("should return an error if the country iso alpha2 code is invalid", func() {
				entry, err := testStorageService.CreateEntry(seedUsers[0].Id, seedAngerTiers[0].Id, "XYZ", "some valid content")
				Expect(err).ToNot(BeNil())
				Expect(err).To(BeAssignableToTypeOf(storage.CountryCodeInvalidError{}))
				Expect(entry).To(BeNil())
			})
		})
	})

	Context("GetEntryById", func() {
		Context("success", func() {
			It("should return an entry given a valid id", func() {
				entry, err := testStorageService.GetEntryById(seedEntries[0].Id)
				Expect(err).ToNot(HaveOccurred())
				Expect(entry.Id).To(Equal(seedEntries[0].Id))
				Expect(entry.TextContent).To(Equal(seedEntries[0].TextContent))
				Expect(entry.UserId).To(Equal(seedEntries[0].UserId))
				Expect(entry.User.Id).To(Equal(seedEntries[0].UserId))
				Expect(entry.User.Country.IsoAlpha2Code).To(Equal(seedEntries[0].CountryIsoAlpha2Code))
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

	Context("EditEntry", func() {
		Context("success", func() {
			It("should update an entry given valid values", func() {
				entry, err := testStorageService.EditEntry(seedEntries[0].Id, seedUsers[0].Id, "some updated valid content")
				Expect(err).ToNot(HaveOccurred())
				Expect(entry.Id).To(Equal(seedEntries[0].Id))
				Expect(entry.TextContent).To(Equal("some updated valid content"))
				Expect(entry.UserId).To(Equal(seedEntries[0].UserId))
				Expect(entry.User.Id).To(Equal(seedEntries[0].UserId))
				Expect(entry.User.Country.IsoAlpha2Code).To(Equal(seedEntries[0].CountryIsoAlpha2Code))
				Expect(entry.CountryIsoAlpha2Code).To(Equal(seedEntries[0].CountryIsoAlpha2Code))
				Expect(entry.Country.IsoAlpha2Code).To(Equal(seedEntries[0].CountryIsoAlpha2Code))
				Expect(entry.AngerTierId).To(Equal(seedAngerTiers[0].Id))
				Expect(entry.AngerTier.Id).To(Equal(seedAngerTiers[0].Id))
				Expect(entry.CreatedAt).To(BeTemporally(">=", seedTimeNow))
				Expect(entry.UpdatedAt).To(BeTemporally(">", entry.CreatedAt))
			})
		})

		Context("failure", func() {
			It("should return an error if the entry id is invalid", func() {
				entry, err := testStorageService.EditEntry(88888, seedUsers[0].Id, "some updated valid content")
				Expect(err).To(BeAssignableToTypeOf(storage.RecordNotFoundError{}))
				Expect(entry).To(BeNil())
			})

			It("should return an error if the user id is invalid", func() {
				entry, err := testStorageService.EditEntry(seedEntries[0].Id, 88888, "some updated valid content")
				Expect(err).To(BeAssignableToTypeOf(storage.RecordNotFoundError{}))
				Expect(entry).To(BeNil())
			})
		})
	})

	Context("GetEntries", func() {
		Context("success", func() {
			It("should return a list of entries before the timestamp within the size limit (all)", func() {
				seedTotalCount := len(seedEntries)
				afterAllTimestamp := seedEntries[seedTotalCount-1].CreatedAt.Add(time.Millisecond).UnixMicro()

				entries, err := testStorageService.GetEntries(afterAllTimestamp, seedTotalCount+1, nil)
				Expect(err).ToNot(HaveOccurred())

				Expect(entries).To(HaveLen(seedTotalCount))
				Expect(entries[0].Id).To(Equal(seedEntries[2].Id))
				Expect(entries[0].TextContent).To(Equal(seedEntries[2].TextContent))
				Expect(entries[0].UserId).To(Equal(seedEntries[2].UserId))
				Expect(entries[0].User.Id).To(Equal(seedEntries[2].UserId))
				Expect(entries[0].User.Country.IsoAlpha2Code).To(Equal(seedEntries[2].CountryIsoAlpha2Code))
				Expect(entries[0].CountryIsoAlpha2Code).To(Equal(seedEntries[2].CountryIsoAlpha2Code))
				Expect(entries[0].Country.IsoAlpha2Code).To(Equal(seedEntries[2].CountryIsoAlpha2Code))
				Expect(entries[0].AngerTierId).To(Equal(seedEntries[2].AngerTierId))
				Expect(entries[0].AngerTier.Id).To(Equal(seedEntries[2].AngerTierId))
				Expect(entries[0].CreatedAt).To(BeTemporally(">=", seedTimeNow))
				Expect(entries[0].UpdatedAt).To(BeTemporally(">=", seedTimeNow))

				Expect(entries[1].Id).To(Equal(seedEntries[1].Id))
				Expect(entries[1].TextContent).To(Equal(seedEntries[1].TextContent))
				Expect(entries[1].UserId).To(Equal(seedEntries[1].UserId))
				Expect(entries[1].User.Id).To(Equal(seedEntries[1].UserId))
				Expect(entries[1].User.Country.IsoAlpha2Code).To(Equal(seedEntries[1].CountryIsoAlpha2Code))
				Expect(entries[1].CountryIsoAlpha2Code).To(Equal(seedEntries[1].CountryIsoAlpha2Code))
				Expect(entries[1].Country.IsoAlpha2Code).To(Equal(seedEntries[1].CountryIsoAlpha2Code))
				Expect(entries[1].AngerTierId).To(Equal(seedEntries[1].AngerTierId))
				Expect(entries[1].AngerTier.Id).To(Equal(seedEntries[1].AngerTierId))
				Expect(entries[1].CreatedAt).To(BeTemporally(">=", seedTimeNow))
				Expect(entries[1].UpdatedAt).To(BeTemporally(">=", seedTimeNow))

				Expect(entries[2].Id).To(Equal(seedEntries[0].Id))
				Expect(entries[2].TextContent).To(Equal(seedEntries[0].TextContent))
				Expect(entries[2].UserId).To(Equal(seedEntries[0].UserId))
				Expect(entries[2].User.Id).To(Equal(seedEntries[0].UserId))
				Expect(entries[2].User.Country.IsoAlpha2Code).To(Equal(seedEntries[0].CountryIsoAlpha2Code))
				Expect(entries[2].CountryIsoAlpha2Code).To(Equal(seedEntries[0].CountryIsoAlpha2Code))
				Expect(entries[2].Country.IsoAlpha2Code).To(Equal(seedEntries[0].CountryIsoAlpha2Code))
				Expect(entries[2].AngerTierId).To(Equal(seedEntries[0].AngerTierId))
				Expect(entries[2].AngerTier.Id).To(Equal(seedEntries[0].AngerTierId))
				Expect(entries[2].CreatedAt).To(BeTemporally(">=", seedTimeNow))
				Expect(entries[2].UpdatedAt).To(BeTemporally(">=", seedTimeNow))
			})
		})

		It("should return a list of entries before the timestamp within the size limit (partial)", func() {
			seedTotalCount := len(seedEntries)
			afterAllTimestamp := seedEntries[seedTotalCount-1].CreatedAt.Add(time.Millisecond).UnixMicro()

			entries, err := testStorageService.GetEntries(afterAllTimestamp, 1, nil)
			Expect(err).ToNot(HaveOccurred())

			Expect(entries).To(HaveLen(1))
			Expect(entries[0].Id).To(Equal(seedEntries[2].Id))
			Expect(entries[0].TextContent).To(Equal(seedEntries[2].TextContent))
			Expect(entries[0].UserId).To(Equal(seedEntries[2].UserId))
			Expect(entries[0].User.Id).To(Equal(seedEntries[2].UserId))
			Expect(entries[0].User.Country.IsoAlpha2Code).To(Equal(seedEntries[2].CountryIsoAlpha2Code))
			Expect(entries[0].CountryIsoAlpha2Code).To(Equal(seedEntries[2].CountryIsoAlpha2Code))
			Expect(entries[0].Country.IsoAlpha2Code).To(Equal(seedEntries[2].CountryIsoAlpha2Code))
			Expect(entries[0].AngerTierId).To(Equal(seedEntries[2].AngerTierId))
			Expect(entries[0].AngerTier.Id).To(Equal(seedEntries[2].AngerTierId))
			Expect(entries[0].CreatedAt).To(BeTemporally(">=", seedTimeNow))
			Expect(entries[0].UpdatedAt).To(BeTemporally(">=", seedTimeNow))
		})

		It("should return an empty list of entries if the size limit is zero", func() {
			seedTotalCount := len(seedEntries)
			afterAllTimestamp := seedEntries[seedTotalCount-1].CreatedAt.Add(time.Millisecond).UnixMicro()

			entries, err := testStorageService.GetEntries(afterAllTimestamp, 0, nil)
			Expect(err).ToNot(HaveOccurred())
			Expect(entries).To(BeEmpty())
		})

		It("should return an empty list of entries if the timestamp is before the earliest entry", func() {
			seedTotalCount := len(seedEntries)
			afterAllTimestamp := seedEntries[0].CreatedAt.Add(-time.Millisecond).UnixMicro()

			entries, err := testStorageService.GetEntries(afterAllTimestamp, seedTotalCount+1, nil)
			Expect(err).ToNot(HaveOccurred())
			Expect(entries).To(BeEmpty())
		})

		It("should return a list of entries before the timestamp within the size limit (partial, mid)", func() {
			seedTotalCount := len(seedEntries)
			afterAllTimestamp := seedEntries[seedTotalCount-2].CreatedAt.Add(time.Millisecond).UnixMicro()

			entries, err := testStorageService.GetEntries(afterAllTimestamp, seedTotalCount+1, nil)
			Expect(err).ToNot(HaveOccurred())

			Expect(entries).To(HaveLen(2))
			Expect(entries[0].Id).To(Equal(seedEntries[1].Id))
			Expect(entries[0].TextContent).To(Equal(seedEntries[1].TextContent))
			Expect(entries[0].UserId).To(Equal(seedEntries[1].UserId))
			Expect(entries[0].User.Id).To(Equal(seedEntries[1].UserId))
			Expect(entries[0].User.Country.IsoAlpha2Code).To(Equal(seedEntries[1].CountryIsoAlpha2Code))
			Expect(entries[0].CountryIsoAlpha2Code).To(Equal(seedEntries[1].CountryIsoAlpha2Code))
			Expect(entries[0].Country.IsoAlpha2Code).To(Equal(seedEntries[1].CountryIsoAlpha2Code))
			Expect(entries[0].AngerTierId).To(Equal(seedEntries[1].AngerTierId))
			Expect(entries[0].AngerTier.Id).To(Equal(seedEntries[1].AngerTierId))
			Expect(entries[0].CreatedAt).To(BeTemporally(">=", seedTimeNow))
			Expect(entries[0].UpdatedAt).To(BeTemporally(">=", seedTimeNow))

			Expect(entries[1].Id).To(Equal(seedEntries[0].Id))
			Expect(entries[1].TextContent).To(Equal(seedEntries[0].TextContent))
			Expect(entries[1].UserId).To(Equal(seedEntries[0].UserId))
			Expect(entries[1].User.Id).To(Equal(seedEntries[0].UserId))
			Expect(entries[1].User.Country.IsoAlpha2Code).To(Equal(seedEntries[0].CountryIsoAlpha2Code))
			Expect(entries[1].CountryIsoAlpha2Code).To(Equal(seedEntries[0].CountryIsoAlpha2Code))
			Expect(entries[1].Country.IsoAlpha2Code).To(Equal(seedEntries[0].CountryIsoAlpha2Code))
			Expect(entries[1].AngerTierId).To(Equal(seedEntries[0].AngerTierId))
			Expect(entries[1].AngerTier.Id).To(Equal(seedEntries[0].AngerTierId))
			Expect(entries[1].CreatedAt).To(BeTemporally(">=", seedTimeNow))
			Expect(entries[1].UpdatedAt).To(BeTemporally(">=", seedTimeNow))
		})

		It("should return a list of entries before the timestamp filtered by user id", func() {
			seedTotalCount := len(seedEntries)
			afterAllTimestamp := seedEntries[seedTotalCount-1].CreatedAt.Add(time.Millisecond).UnixMicro()

			entries, err := testStorageService.GetEntries(afterAllTimestamp, seedTotalCount+1, &seedUsers[1].Id)
			Expect(err).ToNot(HaveOccurred())

			Expect(entries).To(HaveLen(1))
			Expect(entries[0].Id).To(Equal(seedEntries[2].Id))
			Expect(entries[0].TextContent).To(Equal(seedEntries[2].TextContent))
			Expect(entries[0].UserId).To(Equal(seedEntries[2].UserId))
			Expect(entries[0].User.Id).To(Equal(seedEntries[2].UserId))
			Expect(entries[0].User.Country.IsoAlpha2Code).To(Equal(seedEntries[2].CountryIsoAlpha2Code))
			Expect(entries[0].CountryIsoAlpha2Code).To(Equal(seedEntries[2].CountryIsoAlpha2Code))
			Expect(entries[0].Country.IsoAlpha2Code).To(Equal(seedEntries[2].CountryIsoAlpha2Code))
			Expect(entries[0].AngerTierId).To(Equal(seedEntries[2].AngerTierId))
			Expect(entries[0].AngerTier.Id).To(Equal(seedEntries[2].AngerTierId))
			Expect(entries[0].CreatedAt).To(BeTemporally(">=", seedTimeNow))
			Expect(entries[0].UpdatedAt).To(BeTemporally(">=", seedTimeNow))
		})
	})
})
