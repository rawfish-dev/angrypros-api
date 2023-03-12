package storage

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/rawfish-dev/angrypros-api/config"
	"github.com/rawfish-dev/angrypros-api/models"
)

const (
	defaultPageSize = 10
)

var _ StorageService = new(Service)

type StorageService interface {
	UserStorage
	CountryStorage
	EntryStorage
}

type UserStorage interface {
	CreateUser(firebaseUserId, title, emailAddress, countryIsoAlpha2Code string) (*models.User, error)
	EditUser(user models.User, title, countryIsoAlpha2Code string) (*models.User, error)
	GetUserById(userId int64) (*models.User, error)
	GetUserByFirebaseUserId(firebaseUserId string) (*models.User, error)
	GetUserByEmailAddress(emailAddress string) (*models.User, error)
}

type CountryStorage interface {
	GetAllCountries() ([]models.Country, error)
}

type EntryStorage interface {
	GetAllAngerTiers() ([]models.AngerTier, error)
	CreateEntry(userId, angerTierId int64,
		isoAlpha2Code, textContent string) (*models.Entry, error)
	GetEntryById(entryId int64) (*models.Entry, error)
	EditEntry(entryId, userId int64, textContent string) (*models.Entry, error)
	GetEntries(beforeTimestampMicro int64, size int, userIdFilter *int64) ([]models.Entry, error)
}

type Service struct {
	db *gorm.DB
}

func NewService(p config.PostgresConfig) (*Service, error) {
	connectionStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		p.Host, p.Port, p.Username, p.Password, p.Database, p.SSLMode)

	db, err := gorm.Open(postgres.Open(connectionStr), &gorm.Config{
		// Logger: logger.New(
		// 	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		// 	logger.Config{
		// 		SlowThreshold: time.Second, // Slow SQL threshold
		// 		LogLevel:      logger.Info, // Log level
		// 	},
		// ),
	})
	if err != nil {
		return nil, ConnectionError{err.Error()}
	}

	err = db.AutoMigrate(&models.User{}, &models.Country{},
		&models.AngerTier{}, &models.Entry{},
		&models.Rarity{}, &models.Award{})
	if err != nil {
		return nil, GeneralDBError{fmt.Sprintf("could not auto migrate due to %s", err)}
	}

	return &Service{
		db: db,
	}, nil
}

func paginate(db *gorm.DB, offset, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if offset < 0 {
			offset = 0
		}
		if size <= 0 {
			size = defaultPageSize
		}

		return db.Offset(offset).Limit(size)
	}
}
