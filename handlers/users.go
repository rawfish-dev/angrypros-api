package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/rawfish-dev/angrypros-api/config"
	"github.com/rawfish-dev/angrypros-api/models"
	"github.com/rawfish-dev/angrypros-api/services/storage"
)

const (
	invalidEmailAddressMessage        = "email is invalid"
	userEmailAlreadyRegisteredMessage = "email is already in use"
)

var (
	errEmailAlreadyExists = errors.New(userEmailAlreadyRegisteredMessage)
)

// Used for creating users and most of editing users
type BaseUserRequest struct {
	Title                string `json:"title"`
	CountryIsoAlpha2Code string `json:"countryIsoAlpha2Code"`
}

func (b BaseUserRequest) validate(userConfig config.UserConfig) []error {
	var validationErrors []error

	// TODO:: Add back regex, only validating length currently
	// re := regexp.MustCompile(userConfig.UsernameRegex)
	if len(b.Title) < userConfig.TitleMinimumLength ||
		len(b.Title) > userConfig.TitleMaximumLength {
		validationErrors = append(validationErrors,
			fmt.Errorf("title must be at least %d and at most %d in length",
				userConfig.TitleMinimumLength, userConfig.TitleMaximumLength))
	}

	return validationErrors
}

type EditUserRequest struct {
	BaseUserRequest
}

func (e EditUserRequest) validate(userConfig config.UserConfig) []error {
	validationErrors := e.BaseUserRequest.validate(userConfig)

	// if c.ProfileImageMediaId != nil {
	// 	parsedUrl, err := url.Parse(*c.ProfileImageUrl)
	// 	if err != nil {
	// 		validationErrors = append(validationErrors,
	// 			errors.New("profile image url domain is invalid"))

	// 		return validationErrors
	// 	}

	// 	urlProtocolHost := fmt.Sprintf("%s://%s", parsedUrl.Scheme, parsedUrl.Host)

	// 	invalidMediaDomain := false
	// 	for _, whitelistedDomain := range mediaDomainWhitelist {
	// 		if urlProtocolHost == whitelistedDomain {
	// 			invalidMediaDomain = true
	// 			break
	// 		}
	// 	}

	// 	if !invalidMediaDomain {
	// 		validationErrors = append(validationErrors,
	// 			errors.New("profile image url domain is invalid"))
	// 	}
	// }

	return validationErrors
}

type CurrentUserResponse struct {
	UserResponse
}

type UserResponse struct {
	Id      int64           `json:"id"`
	Title   string          `json:"title"`
	Country CountryResponse `json:"country"`
}

type ForgotPasswordRequest struct {
	EmailAddress   string `json:"emailAddress"`
	RecaptchaToken string `json:"recaptchaToken"`
}

func (f ForgotPasswordRequest) validate() []error {
	return nil
}

func (s Server) CreateUserHandler(c *gin.Context) {
	currentUser, _ := c.Get("currentUser")

	if currentUser != nil {
		UnprocessableRequestError(c, []error{
			errors.New("user id already exists"),
		})
		return
	}

	jsonReqData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		MalformedRequestError(c, err)
		return
	}

	var req BaseUserRequest
	err = json.Unmarshal(jsonReqData, &req)
	if err != nil {
		MalformedRequestError(c, err)
		return
	}

	validationErrors := req.validate(s.config.UserConfig)
	if validationErrors != nil {
		UnprocessableRequestError(c, validationErrors)
		return
	}

	firebaseUserId := c.MustGet("firebaseUserId").(string)
	email, err := s.authService.GetFirebaseUserEmail(firebaseUserId)
	if err != nil {
		InternalServerError(c, err)
		return
	}

	existingUser, err := s.storageService.GetUserByEmailAddress(email)
	if err != nil {
		switch err.(type) {
		case storage.GeneralDBError:
			InternalServerError(c, err)
			return
		}
	}
	if existingUser != nil {
		UnprocessableRequestError(c, []error{
			errEmailAlreadyExists,
		})
		return
	}

	// No transaction added to prevent complete conflicts for simplicity
	user, err := s.storageService.CreateUser(firebaseUserId, req.Title,
		email, req.CountryIsoAlpha2Code)
	if err != nil {
		InternalServerError(c, err)
		return
	}

	resp := buildMinimalUserResponse(*user)

	WrapJSONAPI(c, http.StatusCreated, resp, nil, nil)
}

func (s Server) EditUserHandler(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(*models.User)

	jsonReqData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		MalformedRequestError(c, err)
		return
	}

	var req EditUserRequest
	err = json.Unmarshal(jsonReqData, &req)
	if err != nil {
		MalformedRequestError(c, err)
		return
	}

	validationErrors := req.validate(s.config.UserConfig)
	if validationErrors != nil {
		UnprocessableRequestError(c, validationErrors)
		return
	}

	user, err := s.storageService.EditUser(*currentUser, req.Title, req.CountryIsoAlpha2Code)
	if err != nil {
		InternalServerError(c, err)
		return
	}

	resp := buildMinimalUserResponse(*user)

	WrapJSONAPI(c, http.StatusOK, resp, nil, nil)
}

func (s Server) GetCurrentUserHandler(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(*models.User)

	resp := buildCurrentUserResponse(*currentUser)

	WrapJSONAPI(c, http.StatusOK, resp, nil, nil)
}

func (s Server) ForgotPasswordHandler(c *gin.Context) {
	jsonReqData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		MalformedRequestError(c, err)
		return
	}

	var req ForgotPasswordRequest
	err = json.Unmarshal(jsonReqData, &req)
	if err != nil {
		MalformedRequestError(c, err)
		return
	}

	validationErrors := req.validate()
	if validationErrors != nil {
		UnprocessableRequestError(c, validationErrors)
		return
	}

	// err = s.authService.VerifyRecaptcha(req.RecaptchaToken)
	// if err != nil {
	// 	switch err.(type) {
	// 	case auth.RecaptchaVerificationFailedError:
	// 		log.Printf("recaptcha verification failed for email: %s", req.EmailAddress)
	// 	}

	// 	// Hide recaptcha verification for now
	// 	InternalServerError(c, err)
	// 	return
	// }

	// err = s.authService.SendForgotPasswordEmail(req.EmailAddress)
	// if err != nil {
	// 	log.Printf("forgot password email sending failed due to %s", err)
	// }

	c.Status(http.StatusOK)
}

func buildCurrentUserResponse(user models.User) CurrentUserResponse {
	return CurrentUserResponse{
		UserResponse: buildMinimalUserResponse(user),
	}
}

func buildMinimalUserResponse(user models.User) UserResponse {
	return UserResponse{
		Id:      user.Id,
		Title:   user.Title,
		Country: buildCountryResponse(user.Country),
	}
}

func requestCurrentUserId(c *gin.Context) *int64 {
	currentUser, exists := c.Get("currentUser")
	var currentUserId *int64
	if exists {
		currentUserId = &currentUser.(*models.User).Id
		return currentUserId
	}

	return nil
}
