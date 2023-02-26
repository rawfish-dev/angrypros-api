package handlers

import (
	"log"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/rawfish-dev/angrypros-api/services/auth"
	"github.com/rawfish-dev/angrypros-api/services/storage"
)

const (
	headerKeyAuthorization = "Authorization"
)

// func CORSMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		origin := c.Request.Header.Get("Origin")
// 		if origin == "http://localhost:8081" || origin == "https://modeefi.com" {
// 			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
// 		}
// 		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
// 		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
// 		c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE, POST, OPTIONS, GET, PUT")

// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(204)
// 			return
// 		}

// 		c.Next()
// 	}
// }

func RecoverMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("caught panic due to %s", rec)
				debug.PrintStack()

				InternalServerError(c, nil)
				return
			}
		}()

		c.Next()
	}
}

// Should not be used directly but via publicOrAuthedMiddleware or authedOnlyMiddleware
func optionalAuthMiddleware(a auth.AuthService, s storage.StorageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ignore ALL errors and only set currentUser in the event everything succeeds

		authHeaderValue := c.Request.Header.Get(headerKeyAuthorization)
		if len(authHeaderValue) != 0 {
			// Of the form "Bearer <id token>"
			tokens := strings.Split(authHeaderValue, "Bearer ")
			if len(tokens) == 2 {
				idToken := tokens[1]

				firebaseUserId, err := a.GetFirebaseUserId(idToken)
				if err == nil {
					c.Set("firebaseUserId", firebaseUserId)

					user, err := s.GetUserByFirebaseUserId(firebaseUserId)
					if err == nil && user != nil {
						c.Set("currentUser", user)
					}
				}
			}
		}

		c.Next()
	}
}

// Should not be used directly but via publicOrAuthedMiddleware or authedOnlyMiddleware
func authMiddleware(a auth.AuthService, s storage.StorageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeaderValue := c.Request.Header.Get(headerKeyAuthorization)
		if len(authHeaderValue) == 0 {
			errors := []ResponseError{
				{
					Code:   string(NoAuth),
					Title:  "Id token not found",
					Detail: "Expected id token to be passed",
				},
			}
			WrapJSONAPI(c, http.StatusUnauthorized, nil, errors, nil)
			return
		}

		// Of the form "Bearer <id token>"
		tokens := strings.Split(authHeaderValue, "Bearer ")
		if len(tokens) != 2 {
			errors := []ResponseError{
				{
					Code:   string(InvalidAuth),
					Title:  "Id token invalid",
					Detail: "Passed id token value was invalid",
				},
			}
			WrapJSONAPI(c, http.StatusUnauthorized, nil, errors, nil)
			return
		}
		idToken := tokens[1]

		var firebaseUserId string

		firebaseUserId, err := a.GetFirebaseUserId(idToken)
		if err != nil {
			errors := []ResponseError{
				{
					Code:   string(InvalidAuth),
					Title:  "Id token invalid",
					Detail: "Passed id token value was invalid",
				},
			}
			WrapJSONAPI(c, http.StatusUnauthorized, nil, errors, nil)
			return
		}
		c.Set("firebaseUserId", firebaseUserId)

		user, err := s.GetUserByFirebaseUserId(firebaseUserId)
		if err != nil {
			switch err.(type) {
			case storage.RecordNotFoundError:
				// Skip this error if path is create as this is expected
				if c.Request.Method == "POST" && isRegistrationRelated(c.Request.URL.Path) {
					return
				}

				// Check for authClient.GetUser to get user email during creation

				errors := []ResponseError{
					{
						Code:   string(AuthIncomplete),
						Title:  "Id token unknown",
						Detail: "Id token did not match a known user",
					},
				}
				WrapJSONAPI(c, http.StatusForbidden, nil, errors, nil)
				return

			case storage.GeneralDBError:
				InternalServerError(c, err)
				return
			}
		}

		// Setting with user value nil results in complications when
		// checking against nil later on as the interface value is not
		// nil, so skipping setting here makes it easy to check for
		// direct existence
		if user != nil {
			c.Set("currentUser", user)
		}

		c.Next()
	}
}

func isRegistrationRelated(urlPath string) bool {
	// "/api/users" -> Used for creating the actual user
	return urlPath == "/api/users"
}
