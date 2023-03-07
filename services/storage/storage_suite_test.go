package storage_test

import (
	"fmt"
	"log"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/rawfish-dev/angrypros-api/config"
	"github.com/rawfish-dev/angrypros-api/models"
	"github.com/rawfish-dev/angrypros-api/services/storage"
)

const (
	testPostgresHost     = "localhost"
	testPostgresPort     = "15432"
	testPostgresUser     = "testuser"
	testPostgresPassword = "testpassword"
	testPostgresDBName   = "angrypros-tests"
)

var (
	testStorageService storage.StorageService
	testDB             *gorm.DB
	seedTimeNow        time.Time
	seedCountries      []models.Country
	seedAngerTiers     []models.AngerTier

	seedUsers []models.User
)

var _ = BeforeSuite(func() {
	var err error

	testStorageService, err = storage.NewService(config.PostgresConfig{
		Username: testPostgresUser,
		Password: testPostgresPassword,
		Host:     testPostgresHost,
		Port:     testPostgresPort,
		Database: testPostgresDBName,
		SSLMode:  "disable",
	})
	if err != nil {
		Fail(err.Error())
	}

	connectionStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		testPostgresHost, testPostgresPort, testPostgresUser,
		testPostgresPassword, testPostgresDBName)

	testDB, err = gorm.Open(postgres.Open(connectionStr), &gorm.Config{
		// Logger: logger.New(
		// 	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		// 	logger.Config{
		// 		SlowThreshold: time.Second,  // Slow SQL threshold
		// 		LogLevel:      logger.Error, // Log level
		// 	},
		// ),
	})
	if err != nil {
		Fail(err.Error())
	}

	seedTimeNow = time.Now()

	truncateStatic()
	seedStatic()
})

var _ = BeforeEach(func() {
	truncateDynamic()
	seedDynamic()
})

func truncateStatic() {
	commands := []string{
		"TRUNCATE countries RESTART IDENTITY CASCADE;",
		"TRUNCATE anger_tiers RESTART IDENTITY CASCADE;",
	}

	for _, command := range commands {
		result := testDB.Exec(command)
		if result.Error != nil {
			log.Fatalf("Unable to truncate with command %s due to %v",
				command, result.Error.Error())
		}
	}
}

func seedStatic() {
	prepareSeedCountries()
	prepareAngerTiers()
}

func prepareSeedCountries() {
	seedCountries = []models.Country{
		{
			IsoAlpha2Code: "MY",
			Name:          "Malaysia",
		},
		{
			IsoAlpha2Code: "SG",
			Name:          "Singapore",
		},
		{
			IsoAlpha2Code: "AU",
			Name:          "Australia",
		},
	}
	for idx := range seedCountries {
		country := seedCountries[idx]
		result := testDB.Create(&country)
		if result.Error != nil {
			Fail(result.Error.Error())
		}
		seedCountries[idx] = country
	}
}

func prepareAngerTiers() {
	seedAngerTiers = []models.AngerTier{
		{
			Label:     "Cross",
			RageLevel: 3,
		},
		{
			Label:     "Displeased",
			RageLevel: 1,
		},
		{
			Label:     "Annoyed",
			RageLevel: 2,
		},
	}
	for idx := range seedAngerTiers {
		angerTier := seedAngerTiers[idx]
		result := testDB.Create(&angerTier)
		if result.Error != nil {
			Fail(result.Error.Error())
		}
		seedAngerTiers[idx] = angerTier
	}
}

func truncateDynamic() {
	commands := []string{
		"TRUNCATE users RESTART IDENTITY CASCADE;",
	}

	for _, command := range commands {
		result := testDB.Exec(command)
		if result.Error != nil {
			log.Fatalf("Unable to truncate with command %s due to %v",
				command, result.Error.Error())
		}
	}
}

func seedDynamic() {
	prepareSeedUsers()
}

func prepareSeedUsers() {
	seedUsers = []models.User{
		{
			FirebaseUserId:         "abc1",
			Title:                  "VP of Engineering",
			NormalisedEmailAddress: "rawfishy-1@gmail.com",
			CountryIsoAlpha2Code:   "SG",
		},
		{
			FirebaseUserId:         "abc2",
			Title:                  "CEO",
			NormalisedEmailAddress: "rawfishy-2@gmail.com",
			CountryIsoAlpha2Code:   "MY",
		},
	}
	for idx := range seedUsers {
		user := seedUsers[idx]
		result := testDB.Create(&user)
		if result.Error != nil {
			Fail(result.Error.Error())
		}
		seedUsers[idx] = user
	}
}

func TestStorage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Storage Suite")
}
