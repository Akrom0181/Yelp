// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"net/http"

	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Swagger docs.
	"github.com/Akorm0181/yelp/config"
	_ "github.com/Akorm0181/yelp/docs"
	"github.com/Akorm0181/yelp/internal/controller/http/v1/handler"
	"github.com/Akorm0181/yelp/internal/usecase"
	"github.com/Akorm0181/yelp/pkg/logger"
	rediscache "github.com/golanguzb70/redis-cache"
)

// NewRouter -.
// Swagger spec:
// @title       Yelp API
// @description This is a sample server Yelp server.
// @version     1.0
// @BasePath    /v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func NewRouter(engine *gin.Engine, l *logger.Logger, config *config.Config, useCase *usecase.UseCase, redis rediscache.RedisCache) {
	// Options
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	handlerV1 := handler.NewHandler(l, config, useCase, redis)

	// Initialize Casbin enforcer
	e := casbin.NewEnforcer("config/rbac.conf", "config/policy.csv")
	engine.Use(handlerV1.AuthMiddleware(e))

	// Swagger
	url := ginSwagger.URL("swagger/doc.json") // The url pointing to API definition
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// K8s probe
	engine.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routes
	v1 := engine.Group("/v1")

	user := v1.Group("/user")
	{
		user.POST("/", handlerV1.CreateUser)
		user.GET("/list", handlerV1.GetUsers)
		user.GET("/:id", handlerV1.GetUser)
		user.PUT("/", handlerV1.UpdateUser)
		user.DELETE("/:id", handlerV1.DeleteUser)
		user.POST("/upload", handlerV1.UploadProfilePic)
	}

	session := v1.Group("/session")
	{
		session.GET("/list", handlerV1.GetSessions)
		session.GET("/:id", handlerV1.GetSession)
		session.PUT("/", handlerV1.UpdateSession)
		session.DELETE("/:id", handlerV1.DeleteSession)
	}

	auth := v1.Group("/auth")
	{
		auth.POST("/logout", handlerV1.Logout)
		auth.POST("/register", handlerV1.Register)
		auth.POST("/verify-email", handlerV1.VerifyEmail)
		auth.POST("/login", handlerV1.Login)
	}

	business := v1.Group("/business")
	{
		business.POST("/", handlerV1.CreateBusiness)
		business.GET("/list", handlerV1.GetBusinesses)
		business.GET("/:id", handlerV1.GetBusiness)
		business.PUT("/", handlerV1.UpdateBusiness)
		business.DELETE("/:id", handlerV1.DeleteBusiness)
		business.POST("/upload/:id", handlerV1.UploadBusinessPic)
	}

	business_cat := v1.Group("/business-category")
	{
		business_cat.POST("/", handlerV1.CreateBusinessCategory)
		business_cat.GET("/list", handlerV1.GetBusinessCategories)
		business_cat.GET("/:id", handlerV1.GetBusinessCategory)
		business_cat.PUT("/", handlerV1.UpdateBusinessCategory)
		business_cat.DELETE("/:id", handlerV1.DeleteBusinessCategory)
	}

	review := v1.Group("/review")
	{
		review.POST("/", handlerV1.CreateReview)
		review.GET("/list", handlerV1.GetReviews)
		review.GET("/:id", handlerV1.GetReview)
		review.PUT("/", handlerV1.UpdateReview)
		review.DELETE("/:id", handlerV1.DeleteReview)
	}

	report := v1.Group("/report")
	{
		report.POST("/", handlerV1.CreateReport)
		report.GET("/list", handlerV1.GetReports)
		report.GET("/:id", handlerV1.GetReport)
		report.PUT("/", handlerV1.UpdateReport)
		report.DELETE("/:id", handlerV1.DeleteReport)
	}

	notification := v1.Group("/notification")
	{
		notification.POST("/", handlerV1.CreateNotification)
		notification.GET("/list", handlerV1.GetNotifications)
		notification.GET("/:id", handlerV1.GetNotification)
		notification.PUT("/update-status", handlerV1.UpdateStatusNotification)
		notification.DELETE("/:id", handlerV1.DeleteNotification)
		notification.PUT("/:id", handlerV1.UpdateNotification)
	}

	firebase := v1.Group("/firebase")
	{
		firebase.POST("/", handlerV1.UploadFiles)
		firebase.DELETE("/:id", handlerV1.DeleteFile)
	}

	event := v1.Group("/event")
	{
		event.POST("/", handlerV1.CreateEvent)
		event.PUT("/", handlerV1.UpdateEvent)
		event.GET("/list", handlerV1.GetEvents)
		event.GET("/:id", handlerV1.GetEvent)
		event.DELETE("/:id", handlerV1.DeleteEvent)
		event.POST("/add-participant", handlerV1.AddParticipant)
		event.DELETE("/remove-participant", handlerV1.RemoveParticipant)
		event.GET("/:id/participants", handlerV1.GetParticipants)
	}

	bookmark := v1.Group("/bookmark")
	{
		bookmark.POST("/", handlerV1.CreateBookmark)
		bookmark.GET("/list", handlerV1.GetBookmarks)
		bookmark.GET("/:id", handlerV1.GetBookmark)
		bookmark.PUT("/", handlerV1.UpdateBookmark)
		bookmark.DELETE("/:id", handlerV1.DeleteBookmark)
	}

	promotion := v1.Group("/promotion")
	{
		promotion.POST("/", handlerV1.CreatePromotion)
		promotion.GET("/list", handlerV1.GetPromotions)
		promotion.GET("/:id", handlerV1.GetPromotion)
		promotion.DELETE("/:id", handlerV1.DeletePromotion)
	}

	tag := v1.Group("/tag")
	{
		tag.POST("/", handlerV1.CreateTag)
		tag.GET("/list", handlerV1.GetTags)
		tag.GET("/:id", handlerV1.GetTag)
		tag.PUT("/", handlerV1.UpdateTag)
		tag.DELETE("/:id", handlerV1.DeleteTag)
	}

	follower := v1.Group("/follower")
	{
		follower.POST("/", handlerV1.FollowUnfollow)
		follower.GET("/list", handlerV1.GetFollowers)
	}
}
