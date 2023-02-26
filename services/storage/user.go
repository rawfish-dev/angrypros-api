package storage

import (
	"strings"
	"time"

	"github.com/rawfish-dev/angrypros-api/models"
)

func (s Service) CreateUser(firebaseUserId, username, emailAddress, countryIsoAlpha2Code string) (*models.User, error) {
	now := time.Now()

	newUser := models.User{
		FirebaseUserId:         firebaseUserId,
		Username:               username,
		NormalisedUsername:     strings.ToLower(username),
		NormalisedEmailAddress: strings.ToLower(emailAddress),
		CountryIsoAlpha2Code:   countryIsoAlpha2Code,
		CreatedAt:              now,
		UpdatedAt:              now,
	}

	result := s.db.Create(&newUser)
	if result.Error != nil {
		constraintError := filterConstraintErrors(result.Error)
		if constraintError != nil {
			return nil, constraintError
		}

		return nil, GeneralDBError{result.Error.Error()}
	}

	return s.GetUserById(newUser.Id)
}

func (s Service) EditUser(user models.User, username, countryIsoAlpha2Code string) (*models.User, error) {
	now := time.Now()

	editedUser := models.User{
		Username:             username,
		NormalisedUsername:   strings.ToLower(username),
		CountryIsoAlpha2Code: countryIsoAlpha2Code,
		UpdatedAt:            now,
	}

	result := s.db.Model(&user).Updates(editedUser)
	if result.Error != nil {
		constraintError := filterConstraintErrors(result.Error)
		if constraintError != nil {
			return nil, constraintError
		}

		return nil, GeneralDBError{result.Error.Error()}
	}

	return s.GetUserById(user.Id)
}

func (s Service) GetUserById(userId int64) (*models.User, error) {
	var user models.User

	result := s.db.
		Joins("Country").
		Find(&user, models.User{Id: userId})
	if result.Error != nil {
		return nil, GeneralDBError{result.Error.Error()}
	}
	if result.RowsAffected == 0 {
		return nil, RecordNotFoundError{}
	}

	return &user, nil
}

func (s Service) GetUserByFirebaseUserId(firebaseUserId string) (*models.User, error) {
	var user models.User

	result := s.db.
		Joins("Country").
		Find(&user, models.User{FirebaseUserId: firebaseUserId})
	if result.Error != nil {
		return nil, GeneralDBError{result.Error.Error()}
	}
	if result.RowsAffected == 0 {
		return nil, RecordNotFoundError{}
	}

	return &user, nil
}

func (s Service) GetUserByEmailAddress(emailAddress string) (*models.User, error) {
	var user models.User

	result := s.db.Find(&user, models.User{NormalisedEmailAddress: strings.ToLower(emailAddress)})
	if result.Error != nil {
		return nil, GeneralDBError{result.Error.Error()}
	}
	if result.RowsAffected == 0 {
		return nil, RecordNotFoundError{}
	}

	return &user, nil
}
