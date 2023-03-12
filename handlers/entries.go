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
	AngerTierId int `json:"angerTierId"`
}

type EditEntryRequest struct {
	BaseEntryRequest
}

type EntryResponse struct {
	Id          int64             `json:"id"`
	User        UserResponse      `json:"user"`
	AngerTier   AngerTierResponse `json:"angerTier"`
	TextContent string            `json:"textContent"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}

type AngerTierResponse struct {
	Id        int64  `json:"id"`
	Label     string `json:"label"`
	RageLevel int    `json:"rageLevel"`
}

type FeedResponse struct {
	Entries             []EntryResponse `json:"entries"`
	QueryTimestampMicro int64           `json:"queryTimestampMicro"`
	MoreResults         bool            `json:"moreResults"`
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

	entry, err := s.storageService.CreateEntry(currentUser.Id, int64(req.AngerTierId),
		currentUser.Country.IsoAlpha2Code, req.TextContent)
	if err != nil {
		InternalServerError(c, err)
		return
	}

	resp := buildEntryResponse(*entry)

	WrapJSONAPI(c, http.StatusCreated, resp, nil, nil)
}

func (s Server) GetEntryDetailsHandler(c *gin.Context) {
	entryIdStr := c.Param("entryId")

	entryId, err := strconv.ParseInt(entryIdStr, 10, 64)
	if err != nil {
		ResourceNotFoundError(c)
		return
	}

	entry, err := s.storageService.GetEntryById(entryId)
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

func (s Server) EditEntryHandler(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(*models.User)

	entryIdStr := c.Param("entryId")

	entryId, err := strconv.ParseInt(entryIdStr, 10, 64)
	if err != nil {
		ResourceNotFoundError(c)
		return
	}

	jsonReqData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		MalformedRequestError(c, err)
		return
	}

	var req EditEntryRequest
	err = json.Unmarshal(jsonReqData, &req)
	if err != nil {
		MalformedRequestError(c, err)
		return
	}

	entry, err := s.storageService.EditEntry(entryId, currentUser.Id, req.TextContent)
	if err != nil {
		InternalServerError(c, err)
		return
	}

	resp := buildEntryResponse(*entry)

	WrapJSONAPI(c, http.StatusOK, resp, nil, nil)
}

func (s Server) GetFeedHandler(c *gin.Context) {
	beforeTimestampMicro := s.getParsedBeforeTimestamp(c.Query("before"))

	userId := s.getParsedUserId(c.Query("userId"))
	var userIdFilter *int64
	if userId != 0 {
		userIdFilter = &userId
	}

	entryables, err := s.feedService.GetFeedItems(beforeTimestampMicro, s.config.FeedConfig.DefaultPageSize, userIdFilter)
	if err != nil {
		InternalServerError(c, err)
		return
	}

	resp := buildFeedResponse(s.config.FeedConfig, entryables, beforeTimestampMicro)

	WrapJSONAPI(c, http.StatusOK, resp, nil, nil)
}

func (s Server) getParsedBeforeTimestamp(beforeTimestampStr string) int64 {
	var beforeTimestampMicro int64
	if beforeTimestampStr != "" {
		beforeTimestampMicro, _ = strconv.ParseInt(beforeTimestampStr, 10, 64)
	}
	if beforeTimestampMicro == 0 {
		// In case of not being passed or errors in parsing
		beforeTimestampMicro = s.timeService.Now().UnixMicro()
	}

	return beforeTimestampMicro
}

func (s Server) getParsedUserId(userIdStr string) int64 {
	var userId int64
	if userIdStr != "" {
		userId, _ = strconv.ParseInt(userIdStr, 10, 64)
	}

	return userId
}

func buildEntryResponse(entry models.Entry) EntryResponse {
	return EntryResponse{
		Id:          entry.Id,
		User:        buildMinimalUserResponse(entry.User),
		AngerTier:   buildAngerTierResponse(entry.AngerTier),
		TextContent: entry.TextContent,
		CreatedAt:   entry.CreatedAt,
		UpdatedAt:   entry.UpdatedAt,
	}
}

func buildAngerTierResponse(angerTier models.AngerTier) AngerTierResponse {
	return AngerTierResponse{
		Id:        angerTier.Id,
		Label:     angerTier.Label,
		RageLevel: angerTier.RageLevel,
	}
}

func buildFeedResponse(feedConfig config.FeedConfig, entries []models.Entry,
	beforeTimestampMicro int64) FeedResponse {
	entryResponses := buildEntryResponses(entries)

	return FeedResponse{
		Entries:             entryResponses,
		QueryTimestampMicro: beforeTimestampMicro,
		MoreResults:         len(entryResponses) == feedConfig.DefaultPageSize,
	}
}

func buildEntryResponses(entries []models.Entry) []EntryResponse {
	entryResponses := make([]EntryResponse, len(entries))

	for idx := range entries {
		entryResponses[idx] = buildEntryResponse(entries[idx])
	}

	return entryResponses
}

func buildAngerTierResponses(angerTiers []models.AngerTier) []AngerTierResponse {
	angerTierResponses := make([]AngerTierResponse, len(angerTiers))

	for idx := range angerTiers {
		angerTierResponses[idx] = buildAngerTierResponse(angerTiers[idx])
	}

	return angerTierResponses
}
