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
	go h.services.DirtyMoves()
	h.services.TelegramBot()

	router.GET("/online", h.getOnline)
	router.GET("/crash", h.handleConnectionsCrash)
	router.GET("/roulette", h.handleConnectionsRoulette)
	router.GET("/crash/init-bets-for-new-client", h.initCrashBetsForNewClient)
	router.GET("/roulette/init-bets-for-new-client", h.initRouletteBetsForNewClient)

	router.GET("/gems-prices", h.getPositionPrices)

	router.GET("/all-crash-records", h.getAllCrashRecords)
	router.GET("/all-roulette-records", h.getAllRouletteRecords)

	router.GET("/last-drops", h.getLastDrops)
	router.GET("/drops", h.handleConnectionsDrop)

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
		cases.GET("/last-drops", h.getLastDrops)
	}

	authenticated := router.Group("/authenticated", h.userIdentity)
	{
		authenticated.GET("/open-case", h.openCase)
		authenticated.POST("/replenishment", h.NewReplenishment)

		withdraw := authenticated.Group("/withdraw")
		{
			withdraw.POST("/create", h.createWithdraw)
			withdraw.GET("/info", h.getUsersWithdraws)
		}

		user := authenticated.Group("/user")
		{
			user.PATCH("/update", h.updateUser)
			user.GET("/user", h.getUserById)
			user.GET("/sell-item", h.sellItem)
			user.GET("/sell-all-items", h.sellAllItems)
			user.GET("/change-photo", h.changeAvatar)
		}

	}
	return router
}
