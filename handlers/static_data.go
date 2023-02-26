package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/rawfish-dev/angrypros-api/models"
)

type CountryResponse struct {
	IsoAlpha2Code string `json:"isoAlpha2Code"`
	Name          string `json:"name"`
}

type CountriesResponse struct {
	Countries []CountryResponse `json:"countries"`
}

func (s Server) GetCountriesHandler(c *gin.Context) {
	countries, err := s.storageService.GetAllCountries()
	if err != nil {
		InternalServerError(c, err)
		return
	}

	countryResponses := make([]CountryResponse, len(countries))
	for idx := range countries {
		countryResponses[idx] = buildCountryResponse(countries[idx])
	}

	resp := CountriesResponse{
		Countries: countryResponses,
	}

	WrapJSONAPI(c, http.StatusOK, resp, nil, nil)
}

func buildCountryResponse(country models.Country) CountryResponse {
	return CountryResponse{
		IsoAlpha2Code: country.IsoAlpha2Code,
		Name:          country.Name,
	}
}
