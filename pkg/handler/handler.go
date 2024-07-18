package handler

import (
	"gems_go_back/docs"
	"gems_go_back/pkg/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

// @title Gems API
// @version 1.0
// @description API for managing gems
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1

func (h *Handler) InitRoutes() *gin.Engine {

	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Укажите домены, которым разрешен доступ
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	go h.services.Crash.CrashGame()
	go h.services.Roulette.RouletteGame()
	h.services.TelegramBot()

	router.GET("/crash", h.handleConnectionsCrash)
	router.GET("/roulette", h.handleConnectionsRoulette)
	router.GET("/crash/init-bets-for-new-client", h.initCrashBetsForNewClient)
	router.GET("/roulette/init-bets-for-new-client", h.initRouletteBetsForNewClient)

	router.GET("/all-crash-records", h.getAllCrashRecords)
	router.GET("/all-roulette-records", h.getAllRouletteRecords)
	router.POST("/replenishment", h.NewReplenishment)

	admin := router.Group("/admin")
	{
		admin.POST("/change-status", h.adminChangeStatus)
	}

	fk := router.Group("/fk")
	{
		fk.POST("/msg", h.MSGFromFrekassa)
		fk.GET("/accepted", h.RedirectAccepted)
		fk.POST("/denied", h.RedirectDenied)
	}

	auth := router.Group("/user")
	{
		auth.POST("/sign-in", h.signIn)
		auth.POST("/sign-up", h.signUp)
		auth.PATCH("/update", h.updateUser)
		auth.GET("/user", h.getUserById)
		auth.GET("/sell-item", h.sellItem)
		auth.GET("/sell-all-items", h.sellAllItems)
	}

	item := router.Group("/item")
	{
		item.POST("/create", h.createItem)
		item.GET("/get-item", h.getItem)
		item.GET("/get-all-items", h.getAllItems)
		item.PATCH("/update", h.updateItem)
	}

	cases := router.Group("/case")
	{
		cases.POST("/create", h.createCase)
		cases.GET("/get-case", h.getCase)
		cases.GET("/get-all-cases", h.getAllCases)
		cases.PUT("/update", h.updateCase)
		cases.DELETE("/delete", h.deleteCase)
		cases.GET("/case-story", h.getAllCaseRecords)
	}

	games := router.Group("/games", h.userIdentity)
	{
		games.GET("/open", h.openCase)
	}

	withdraw := router.Group("/withdraw")
	{
		withdraw.POST("/create", h.createWithdraw)
	}
	return router
}
