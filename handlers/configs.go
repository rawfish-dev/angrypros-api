package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rawfish-dev/angrypros-api/config"
	"github.com/rawfish-dev/angrypros-api/models"
)

type CreateEntryConfigResponse struct {
	EntryTextContentMinimumLength int                 `json:"entryTextContentMinimumLength"`
	EntryTextContentMaximumLength int                 `json:"entryTextContentMaximumLength"`
	AngerTiers                    []AngerTierResponse `json:"angerTiers"`
}

func (s Server) GetCreateEntryConfigHandler(c *gin.Context) {
	entryConfig := s.config.EntryConfig

	angerTiers, err := s.storageService.GetAllAngerTiers()
	if err != nil {
		InternalServerError(c, err)
		return
	}

	resp := buildCreateEntryResponse(entryConfig, angerTiers)

	WrapJSONAPI(c, http.StatusOK, resp, nil, nil)
}

func buildCreateEntryResponse(entryConfig config.EntryConfig, angerTiers []models.AngerTier) CreateEntryConfigResponse {
	return CreateEntryConfigResponse{
		EntryTextContentMinimumLength: 1, // TODO :: Move to config
		EntryTextContentMaximumLength: entryConfig.EntryTextContentMaximumLength,
		AngerTiers:                    buildAngerTierResponses(angerTiers),
	}
}
