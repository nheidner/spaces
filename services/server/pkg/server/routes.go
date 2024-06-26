package server

import (
	"spaces-p/pkg/common"
	"spaces-p/pkg/controllers"
	"spaces-p/pkg/middlewares"
	localmemory "spaces-p/pkg/repositories/local_memory"
	"spaces-p/pkg/repositories/redis_repo"
	"spaces-p/pkg/services"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

// add all routes
func addRoutes(
	apiVersion string,
	router *gin.Engine,
	logger common.Logger,
	redisClient *redis.Client,
	postgresClient *sqlx.DB,
	authClient common.AuthClient,
	geoCodeRepo common.GeocodeRepository,
) {
	api := router.Group("/" + apiVersion)

	// set repos
	redisRepo := redis_repo.NewRedisRepository(redisClient)
	localMemoryRepo := localmemory.NewLocalMemoryRepo()

	// set up services
	userService := services.NewUserService(logger, redisRepo)
	spaceService := services.NewSpaceService(logger, redisRepo, localMemoryRepo)
	spaceNotificationService := services.NewSpaceNotificationsService(logger, redisRepo, localMemoryRepo)
	threadService := services.NewThreadService(logger, redisRepo, localMemoryRepo)
	messageService := services.NewMessageService(logger, redisRepo, localMemoryRepo)
	addressService := services.NewAddressService(logger, redisRepo, geoCodeRepo)
	healthService := services.NewHealthService(logger, postgresClient)

	// set up controllers
	userController := controllers.NewUserController(logger, userService, authClient)
	spaceController := controllers.NewSpaceController(logger, spaceService, spaceNotificationService, threadService, messageService)
	addressController := controllers.NewAddressController(logger, addressService)
	healthController := controllers.NewHealthController(logger, healthService)

	// middleware functions
	validateThreadInSpaceMiddleware := middlewares.ValidateThreadInSpace(logger, redisRepo)
	validateMessageInThreadMiddleware := middlewares.ValidateMessageInThread(logger, redisRepo)
	isSpaceSubscriberMiddleware := middlewares.IsSpaceSubscriber(logger, redisRepo)

	// USERS
	api.POST("/users", userController.CreateUserFromIdToken)                                                                       // to test
	api.GET("/users/:userid", middlewares.EnsureAuthenticated(logger, authClient, redisRepo, true, false), userController.GetUser) // to test

	// AUTHENTICATED USER
	api.GET("/user",
		middlewares.EnsureAuthenticated(logger, authClient, redisRepo, false, false),
		userController.GetAuthedUser,
	)
	api.PUT("/user", middlewares.EnsureAuthenticated(logger, authClient, redisRepo, true, false)) // TODO
	api.DELETE("/user")                                                                           // TODO

	// SPACES
	api.GET("/spaces", middlewares.EnsureAuthenticated(logger, authClient, redisRepo, true, false), spaceController.GetSpaces)    // tested
	api.POST("/spaces", middlewares.EnsureAuthenticated(logger, authClient, redisRepo, true, false), spaceController.CreateSpace) // tested
	api.GET("/spaces/:spaceid", spaceController.GetSpace)                                                                         // tested
	api.GET("/spaces/:spaceid/updates/ws",
		middlewares.EnsureAuthenticated(logger, authClient, redisRepo, true, true),
		isSpaceSubscriberMiddleware,
		spaceController.SpaceConnect,
	)
	api.GET("/spaces/:spaceid/subscribers", // tested
		middlewares.EnsureAuthenticated(logger, authClient, redisRepo, true, false),
		spaceController.GetSpaceSubscribers,
	)
	api.POST("/spaces/:spaceid/subscribers", // tested
		middlewares.EnsureAuthenticated(logger, authClient, redisRepo, true, false),
		spaceController.AddSpaceSubscriber,
	)
	api.GET("/spaces/:spaceid/toplevel-threads",
		middlewares.EnsureAuthenticated(logger, authClient, redisRepo, true, false),
		spaceController.GetTopLevelThreads,
	)
	api.POST("/spaces/:spaceid/toplevel-threads",
		middlewares.EnsureAuthenticated(logger, authClient, redisRepo, true, false),
		isSpaceSubscriberMiddleware,
		spaceController.CreateTopLevelThread,
	)
	api.GET("/spaces/:spaceid/threads/:threadid",
		validateThreadInSpaceMiddleware,
		spaceController.GetThreadWithMessages,
	)
	api.POST("/spaces/:spaceid/threads/:threadid/messages",
		middlewares.EnsureAuthenticated(logger, authClient, redisRepo, true, false),
		validateThreadInSpaceMiddleware,
		isSpaceSubscriberMiddleware,
		spaceController.CreateMessage,
	)
	api.GET("/spaces/:spaceid/threads/:threadid/messages/:messageid",
		middlewares.EnsureAuthenticated(logger, authClient, redisRepo, true, false),
		validateThreadInSpaceMiddleware,
		validateMessageInThreadMiddleware,
		spaceController.GetMessage,
	)
	api.POST("/spaces/:spaceid/threads/:threadid/messages/:messageid/threads",
		middlewares.EnsureAuthenticated(logger, authClient, redisRepo, true, false),
		validateThreadInSpaceMiddleware,
		validateMessageInThreadMiddleware,
		spaceController.CreateThread,
	)
	api.POST("/spaces/:spaceid/threads/:threadid/messages/:messageid/likes",
		middlewares.EnsureAuthenticated(logger, authClient, redisRepo, true, false),
		validateThreadInSpaceMiddleware,
		validateMessageInThreadMiddleware,
		spaceController.LikeMessage,
	)

	// ADDRESSES
	api.GET("/address",
		middlewares.EnsureAuthenticated(logger, authClient, redisRepo, true, false),
		addressController.GetAddress,
	) // tested

	// HEALTH
	api.GET("/health", healthController.HealthCheck)
	api.GET("/healthz", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "OK"})
	})
}
