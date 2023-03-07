package storage

import (
	"time"

	"github.com/rawfish-dev/angrypros-api/models"
)

func (s Service) GetAllAngerTiers() ([]models.AngerTier, error) {
	var angerTiers []models.AngerTier

	result := s.db.Order("rage_level ASC").Find(&angerTiers)
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
		AngerTierId:          angerTierId,
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

	result := s.db.
		Preload("User").
		Preload("Country").
		Preload("AngerTier").
		Find(&entry)
	if result.Error != nil {
		return nil, GeneralDBError{result.Error.Error()}
	}
	if result.RowsAffected == 0 {
		return nil, RecordNotFoundError{}
	}

	return &entry, nil
}

func (s Service) EditEntry(entryId, userId int64, textContent string) (*models.Entry, error) {
	editedEntry := models.Entry{
		Id:          entryId,
		TextContent: textContent,
	}

	result := s.db.Table("entries").
		Where("id = ? AND user_id = ?", entryId, userId).
		Updates(&editedEntry)
	if result.Error != nil {
		return nil, GeneralDBError{result.Error.Error()}
	}

	return s.GetEntryById(entryId)
}

func (s Service) GetEntries(beforeTimestampMicro int64, size int, userIdFilter *int64) ([]models.Entry, error) {
	var entries []models.Entry

	if size <= 0 {
		return entries, nil
	}

	queryTime := time.UnixMicro(beforeTimestampMicro)

	builtScopes := s.db.
		Preload("User.Country").
		Preload("AngerTier").
		Order("entries.created_at DESC").
		Limit(size).
		Where("entries.created_at < ?", queryTime)

	if userIdFilter != nil {
		builtScopes.Where("user_id = ?", *userIdFilter)
	}

	result := builtScopes.Find(&entries)
	if result.Error != nil {
		return nil, GeneralDBError{result.Error.Error()}
	}

	return entries, nil
}
