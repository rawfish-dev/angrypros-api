package storage

import "github.com/rawfish-dev/angrypros-api/models"

func (s Service) GetAllCountries() ([]models.Country, error) {
	var countries []models.Country

	result := s.db.Order("name asc").Find(&countries)
	if result.Error != nil {
		return nil, GeneralDBError{result.Error.Error()}
	}

	return countries, nil
}

func (s Service) GetCountryByIsoAlpha2Code(IsoAlpha2Code string) (*models.Country, error) {
	return nil, nil
}
