package storage

import (
	"fmt"
	"strings"
)

// TODO:: Breakdown filter constraints by table name to reduce looping?
const (
	userAlreadyRegisteredErr   = "duplicate key value violates unique constraint \"idx_users_firebase_user_id\""
	userCountryCodeInvalidErr  = "insert or update on table \"users\" violates foreign key constraint \"fk_users_country\""
	userIdInvalidErr           = "insert or update on table \"entries\" violates foreign key constraint \"fk_entries_user\""
	entryAngerTierIdInvalidErr = "insert or update on table \"entries\" violates foreign key constraint \"fk_entries_anger_tier\""
	entryCountryCodeInvalidErr = "insert or update on table \"entries\" violates foreign key constraint \"fk_entries_country\""
)

var (
	knownPartialErrorMessages = []string{
		userAlreadyRegisteredErr,
		userCountryCodeInvalidErr,
		userIdInvalidErr,
		entryAngerTierIdInvalidErr,
		entryCountryCodeInvalidErr,
	}
)

type ConnectionError struct {
	message string
}

func (c ConnectionError) Error() string {
	return fmt.Sprintf("connection error caused by %s", c.message)
}

type RecordNotFoundError struct{}

func (r RecordNotFoundError) Error() string {
	return "unable to find requested record"
}

type GeneralDBError struct {
	message string
}

func (g GeneralDBError) Error() string {
	return fmt.Sprintf("general db error caused by %s", g.message)
}

type UserAlreadyRegisteredError struct{}

func (u UserAlreadyRegisteredError) Error() string {
	return "user is already registered"
}

type CountryCodeInvalidError struct{}

func (c CountryCodeInvalidError) Error() string {
	return "country code is invalid"
}

type UserIdInvalidError struct{}

func (u UserIdInvalidError) Error() string {
	return "user id is invalid"
}

type AngerTierIdInvalidError struct{}

func (u AngerTierIdInvalidError) Error() string {
	return "anger tier id is invalid"
}

func filterConstraintErrors(err error) error {
	var matchedErrorString string

	for _, partialErrorMessage := range knownPartialErrorMessages {
		if strings.Contains(err.Error(), partialErrorMessage) {
			matchedErrorString = partialErrorMessage
			break
		}
	}

	switch matchedErrorString {
	case userAlreadyRegisteredErr:
		return UserAlreadyRegisteredError{}
	case userCountryCodeInvalidErr, entryCountryCodeInvalidErr:
		return CountryCodeInvalidError{}
	case userIdInvalidErr:
		return UserIdInvalidError{}
	case entryAngerTierIdInvalidErr:
		return AngerTierIdInvalidError{}
	}

	return nil
}
