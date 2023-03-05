package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rawfish-dev/angrypros-api/config"
	"github.com/rawfish-dev/angrypros-api/models"
	"github.com/rawfish-dev/angrypros-api/services/storage"
)

type BaseEntryRequest struct {
	TextContent string `json:"textContent"`
}

func (b BaseEntryRequest) validate(e config.EntryConfig) []error {
	var validationErrors []error

	return validationErrors
}

type CreateEntryRequest struct {
	BaseEntryRequest
	RageLevel int `json:"rageLevel"` // TODO:: Encrypt in future to prevent user input
}

type UpdateEntryRequest struct {
	BaseEntryRequest
}

type EntryResponse struct {
	Id             int64        `json:"id"`
	User           UserResponse `json:"user"`
	AngerTierLabel string       `json:"angerTierLabel"`
	TextContent    string       `json:"textContent"`
	RageLevel      int          `json:"rageLevel"`
	CreatedAt      time.Time    `json:"createdAt"`
	UpdatedAt      time.Time    `json:"updatedAt"`
}

func (s Server) CreateEntryHandler(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(*models.User)

	jsonReqData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		MalformedRequestError(c, err)
		return
	}

	var req CreateEntryRequest
	err = json.Unmarshal(jsonReqData, &req)
	if err != nil {
		MalformedRequestError(c, err)
		return
	}

	validationErrors := req.validate(s.config.EntryConfig)
	if validationErrors != nil {
		UnprocessableRequestError(c, validationErrors)
		return
	}

	orderedAngerTiers, err := s.storageService.GetAllAngerTiers()
	if err != nil {
		InternalServerError(c, err)
		return
	}

	angerTier := s.mapRageLevelToAngerTier(orderedAngerTiers, req.RageLevel)
	if angerTier == nil {
		// Anger tier cannot be found which shouldn't happen
		InternalServerError(c, err)
		return
	}

	entry, err := s.storageService.CreateEntry(currentUser.Id, angerTier.Id,
		currentUser.Country.IsoAlpha2Code, req.TextContent, req.RageLevel)
	if err != nil {
		InternalServerError(c, err)
		return
	}

	resp := buildEntryResponse(*entry)

	WrapJSONAPI(c, http.StatusCreated, resp, nil, nil)
}

func (s Server) GetEntryDetailsHandler(c *gin.Context) {
	entryId := c.Param("entryId")

	parsedEntryId, err := strconv.ParseInt(entryId, 10, 64)
	if err != nil {
		ResourceNotFoundError(c)
		return
	}

	entry, err := s.storageService.GetEntryById(parsedEntryId)
	if err != nil {
		switch err.(type) {
		case storage.RecordNotFoundError:
			ResourceNotFoundError(c)
			return
		}

		InternalServerError(c, err)
		return
	}

	resp := buildEntryResponse(*entry)

	WrapJSONAPI(c, http.StatusOK, resp, nil, nil)
}

func (s Server) mapRageLevelToAngerTier(orderedAngerTiers []models.AngerTier, rageLevel int) *models.AngerTier {
	// Rage levels dictate order but also represent the actual bands
	for _, angerTier := range orderedAngerTiers {
		// Multiply by 10 as stored rage levels are 1 - 8
		if rageLevel < angerTier.RageLevel*10 {
			return &angerTier
		}
	}

	return nil
}

func buildEntryResponse(entry models.Entry) EntryResponse {
	return EntryResponse{
		Id:             entry.Id,
		User:           buildMinimalUserResponse(entry.User),
		AngerTierLabel: entry.AngerTier.Label,
		TextContent:    entry.TextContent,
		RageLevel:      entry.RageLevel,
		CreatedAt:      entry.CreatedAt,
		UpdatedAt:      entry.UpdatedAt,
	}
}