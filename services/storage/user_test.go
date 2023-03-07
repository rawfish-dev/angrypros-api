package storage_test

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/rawfish-dev/angrypros-api/services/storage"
)

var _ = Describe("User", func() {
	var (
		validFirebaseUserId = "abc123"
		validTitle          = "Chief Person Officer"
		validEmailAddress   = "SOME-USER@example.com"
	)

	Context("CreateUser", func() {
		Context("success", func() {
			It("should create a user and return no error", func() {
				newUser, err := testStorageService.CreateUser(validFirebaseUserId, validTitle, validEmailAddress,
					seedCountries[0].IsoAlpha2Code)
				Expect(err).ToNot(HaveOccurred())

				Expect(newUser.Id).ToNot(BeZero())
				Expect(newUser.FirebaseUserId).To(Equal(validFirebaseUserId))
				Expect(newUser.Title).To(Equal(validTitle))
				Expect(newUser.NormalisedEmailAddress).To(Equal(strings.ToLower(validEmailAddress)))
				Expect(newUser.CountryIsoAlpha2Code).To(Equal(seedCountries[0].IsoAlpha2Code))
				Expect(newUser.Country.IsoAlpha2Code).To(Equal(seedCountries[0].IsoAlpha2Code))
				Expect(newUser.CreatedAt.After(seedTimeNow)).To(BeTrue())
				Expect(newUser.UpdatedAt.After(seedTimeNow)).To(BeTrue())
			})
		})
		Context("failure", func() {
			It("should return an error if a user already has the same firebase user id", func() {
				newUser, err := testStorageService.CreateUser(seedUsers[0].FirebaseUserId, validTitle, validEmailAddress,
					seedCountries[0].IsoAlpha2Code)
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(storage.UserAlreadyRegisteredError{}))
				Expect(newUser).To(BeNil())
			})

			It("should return an error if a user is assigned an invalid country code", func() {
				newUser, err := testStorageService.CreateUser(validFirebaseUserId, validTitle, validEmailAddress,
					"SGX")
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(storage.CountryCodeInvalidError{}))
				Expect(newUser).To(BeNil())
			})
		})
	})

	XContext("EditUser", func() {
		Context("success", func() {})
		Context("failure", func() {})
	})

	Context("GetUserById", func() {
		Context("success", func() {
			It("should return a user given a valid id", func() {
				user, err := testStorageService.GetUserById(seedUsers[0].Id)
				Expect(err).ToNot(HaveOccurred())

				Expect(user.Id).To(Equal(seedUsers[0].Id))
				Expect(user.FirebaseUserId).To(Equal(seedUsers[0].FirebaseUserId))
				Expect(user.Title).To(Equal(seedUsers[0].Title))
				Expect(user.NormalisedEmailAddress).To(Equal(seedUsers[0].NormalisedEmailAddress))
				Expect(user.CountryIsoAlpha2Code).To(Equal(seedUsers[0].CountryIsoAlpha2Code))
				Expect(user.Country.IsoAlpha2Code).To(Equal(seedUsers[0].CountryIsoAlpha2Code))
				Expect(user.CreatedAt).To(BeTemporally("==", seedUsers[0].CreatedAt))
				Expect(user.UpdatedAt).To(BeTemporally(">=", seedUsers[0].UpdatedAt)) // Mod entry hook is updating user
			})
		})

		Context("failure", func() {
			It("should return an error if the user id is invalid", func() {
				user, err := testStorageService.GetUserById(88888)
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(storage.RecordNotFoundError{}))
				Expect(user).To(BeNil())
			})
		})
	})

	Context("GetUserByFirebaseUserId", func() {
		Context("success", func() {
			It("should return a user given a valid id", func() {
				user, err := testStorageService.GetUserByFirebaseUserId(seedUsers[0].FirebaseUserId)
				Expect(err).ToNot(HaveOccurred())

				Expect(user.Id).To(Equal(seedUsers[0].Id))
				Expect(user.FirebaseUserId).To(Equal(seedUsers[0].FirebaseUserId))
				Expect(user.Title).To(Equal(seedUsers[0].Title))
				Expect(user.NormalisedEmailAddress).To(Equal(seedUsers[0].NormalisedEmailAddress))
				Expect(user.CreatedAt).To(BeTemporally("==", seedUsers[0].CreatedAt))
				Expect(user.UpdatedAt).To(BeTemporally(">=", seedUsers[0].UpdatedAt)) // Mod entry hook is updating user
			})
		})

		Context("failure", func() {
			It("should return an error if the user id is invalid", func() {
				user, err := testStorageService.GetUserByFirebaseUserId("88888")
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(storage.RecordNotFoundError{}))
				Expect(user).To(BeNil())
			})
		})
	})

	Context("GetUserByEmailAddress", func() {
		Context("success", func() {
			DescribeTable("should return a user given a valid email", func(username string) {
				user, err := testStorageService.GetUserByEmailAddress(seedUsers[1].NormalisedEmailAddress)
				Expect(err).ToNot(HaveOccurred())

				Expect(user.Id).To(Equal(seedUsers[1].Id))
				Expect(user.FirebaseUserId).To(Equal(seedUsers[1].FirebaseUserId))
				Expect(user.Title).To(Equal(seedUsers[1].Title))
				Expect(user.NormalisedEmailAddress).To(Equal(seedUsers[1].NormalisedEmailAddress))
				Expect(user.CreatedAt).To(BeTemporally("==", seedUsers[1].CreatedAt))
				Expect(user.UpdatedAt).To(BeTemporally("==", seedUsers[1].UpdatedAt))
			},
				Entry("When same casing", "test.USER.2"),
				Entry("When mixed casing", "TEST.UsEr.2"),
			)
		})

		Context("failure", func() {
			It("should return an error if the username is invalid", func() {
				user, err := testStorageService.GetUserByEmailAddress("88888")
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(storage.RecordNotFoundError{}))
				Expect(user).To(BeNil())
			})
		})
	})
})
