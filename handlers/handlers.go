package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/rawfish-dev/angrypros-api/config"
	"github.com/rawfish-dev/angrypros-api/services/auth"
	"github.com/rawfish-dev/angrypros-api/services/feed"
	"github.com/rawfish-dev/angrypros-api/services/storage"
	timeS "github.com/rawfish-dev/angrypros-api/services/time"
)

type Server struct {
	config         config.AppConfig
	router         *gin.Engine
	authService    auth.AuthService
	feedService    feed.FeedService
	storageService storage.StorageService
	timeService    timeS.TimeService
}

func NewServer(config config.AppConfig, a auth.AuthService,
	f feed.FeedService, s storage.StorageService,
	t timeS.TimeService) (*Server, error) {
	return &Server{
		config:         config,
		router:         gin.Default(),
		authService:    a,
		feedService:    f,
		storageService: s,
		timeService:    t,
	}, nil
}

func (s Server) SetupRoutes() {
	s.router.Use(RecoverMiddleware())
	// s.router.Use(CORSMiddleware()) Might only be needed for browser

	apiPublic := s.router.Group("/api/public", optionalAuthMiddleware(s.authService, s.storageService))
	{
		apiPublic.GET("/healthcheck", s.HealthcheckHandler)
		apiPublic.GET("/countries", s.GetCountriesHandler)
		apiPublic.GET("/profiles/:userId", s.GetProfileHandler)
		apiPublic.GET("/feed", s.GetFeedHandler)
		// apiPublic.POST("/forgot-password", s.ForgotPasswordHandler)
	}

	apiAuthed := s.router.Group("/api", authMiddleware(s.authService, s.storageService))
	{
		apiAuthed.GET("/current-user", s.GetCurrentUserHandler)
		apiAuthed.POST("/users", s.CreateUserHandler)
		apiAuthed.PUT("/users", s.EditUserHandler)
		apiAuthed.POST("/entries", s.CreateEntryHandler)
		apiAuthed.GET("/entries/:entryId", s.GetEntryDetailsHandler)
		apiAuthed.PUT("/entries/:entryId", s.EditEntryHandler)
	}

	// s.router.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"*"},
	// 	AllowMethods:     []string{"OPTIONS", "DELETE", "POST", "GET", "PUT", "PATCH"},
	// 	AllowHeaders:     []string{"Origin"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	// AllowOriginFunc: func(origin string) bool {
	// 	//   return origin == "https://github.com"
	// 	// },
	// 	MaxAge: 12 * time.Hour,
	// }))

	// api.Use(cors.Default())
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s Server) HealthcheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}
