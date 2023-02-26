package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseCode string

const (
	AuthIncomplete       ResponseCode = "auth-incomplete"
	GeneralServerError   ResponseCode = "general-server-error"
	InvalidAuth          ResponseCode = "invalid-auth"
	MalformedRequest     ResponseCode = "malformed-request"
	NoAuth               ResponseCode = "no-auth"
	ResourceNotFound     ResponseCode = "resource-not-found"
	UnprocessableRequest ResponseCode = "unprocessable-request"
)

type ResponseError struct {
	Code   string `json:"code,omitempty"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func WrapJSONAPI(c *gin.Context, httpStatus int, payload interface{}, errors []ResponseError, meta interface{}) {
	if payload == nil && errors == nil && meta == nil {
		c.JSON(httpStatus, nil)
		return
	}

	responseObj := make(map[string]interface{})

	if payload != nil {
		responseObj["data"] = payload
	}

	if meta != nil {
		responseObj["meta"] = meta
	}

	if errors != nil {
		responseObj["errors"] = errors
		c.Abort()
	}

	c.JSON(httpStatus, responseObj)
}

func MalformedRequestError(c *gin.Context, err error) {
	log.Printf("returning malformed request error due to %s", err)

	WrapJSONAPI(c, http.StatusBadRequest, nil, []ResponseError{
		{
			Code:   string(MalformedRequest),
			Title:  "Request was malformed",
			Detail: "Request was either malformed or not in the right form",
		},
	}, nil)
}

func ResourceNotFoundError(c *gin.Context) {
	WrapJSONAPI(c, http.StatusNotFound, nil, []ResponseError{
		{
			Code:   string(ResourceNotFound),
			Title:  "Resource was not found",
			Detail: "Resource either does not exist or params were incorrect",
		},
	}, nil)
}

func UnprocessableRequestError(c *gin.Context, errors []error) {
	log.Printf("returning unprocessable request error due to %+v", errors)

	responseErrors := make([]ResponseError, len(errors))
	for i := 0; i < len(errors); i++ {
		responseErrors[i] = ResponseError{
			Code:   string(UnprocessableRequest),
			Title:  "Request could not be processed",
			Detail: errors[i].Error(),
		}
	}

	WrapJSONAPI(c, http.StatusUnprocessableEntity, nil, responseErrors, nil)
}

func InternalServerError(c *gin.Context, err error) {
	log.Printf("returning internal server error due to %v", err)

	WrapJSONAPI(c, http.StatusInternalServerError, nil, []ResponseError{
		{
			Code:   string(GeneralServerError),
			Title:  "General server error",
			Detail: "A server error was encountered, please try again later",
		},
	}, nil)
}
