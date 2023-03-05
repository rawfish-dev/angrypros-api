package storage

import "github.com/rawfish-dev/angrypros-api/models"

func (s Service) GetAllAngerTiers() ([]models.AngerTier, error) {
	var angerTiers []models.AngerTier

	result := s.db.Find(&angerTiers).Order("rage_level asc")
	if result.Error != nil {
		return nil, GeneralDBError{result.Error.Error()}
	}

	return angerTiers, nil
}

func (s Service) CreateEntry(userId, angerTierId int64,
	countryIsoAlpha2Code, textContent string, rageLevel int) (*models.Entry, error) {
	entry := models.Entry{
		UserId:               userId,
		CountryIsoAlpha2Code: countryIsoAlpha2Code,
		AngryTierId:          angerTierId,
		TextContent:          textContent,
		RageLevel:            rageLevel,
	}

	result := s.db.Create(&entry)
	if result.Error != nil {
		return nil, GeneralDBError{result.Error.Error()}
	}

	return &entry, nil
}

func (s Service) GetEntryById(entryId int64) (*models.Entry, error) {
	entry := models.Entry{
		Id: entryId,
	}

	result := s.db.Preload("AngerTier").Find(&entry)
	if result.Error != nil {
		return nil, GeneralDBError{result.Error.Error()}
	}

	return &entry, nil
}
