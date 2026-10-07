package server

import (
	"log"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/raulito1500/merkadapp/config"
	"github.com/raulito1500/merkadapp/database"
	"github.com/raulito1500/merkadapp/middleware"
	BillHandlers "github.com/raulito1500/merkadapp/src/bill/handlers"
	BillRepository "github.com/raulito1500/merkadapp/src/bill/repository"
	BillServices "github.com/raulito1500/merkadapp/src/bill/services"
	MarketListHandlers "github.com/raulito1500/merkadapp/src/market_list/handlers"
	MarketListRepository "github.com/raulito1500/merkadapp/src/market_list/repository"
	MarketListServices "github.com/raulito1500/merkadapp/src/market_list/services"
	NotificationHandler "github.com/raulito1500/merkadapp/src/notification/handlers"
	ProductHandlers "github.com/raulito1500/merkadapp/src/product/handlers"
	ProductRepository "github.com/raulito1500/merkadapp/src/product/repository"
	ProductServices "github.com/raulito1500/merkadapp/src/product/services"
	"go.mongodb.org/mongo-driver/mongo"
)

type Api struct {
	db *mongo.Database
}

func NewApi() *Api {
	return &Api{}
}

func (api *Api) Run() {
	config := config.NewConfig()
	db := database.NewMongoDatabase(config).GetDb()

	authClient, err := middleware.NewFirebaseAuthClient(config)
	if err != nil {
		log.Fatalf("failed to initialize Firebase auth client: %v", err)
	}

	server := gin.Default()
	configCors := cors.DefaultConfig()
	configCors.AllowAllOrigins = true
	configCors.AllowHeaders = append(configCors.AllowHeaders, "Authorization")
	server.Use(cors.New(configCors))
	api.initHandlers(db, server, authClient)
	server.Run(":" + config.Port)
}

func (api *Api) initHandlers(db *mongo.Database, r *gin.Engine, authClient *auth.Client) {
	requireAuth := middleware.RequireFirebaseAuth(authClient)

	notificationHandler := NotificationHandler.NewNotificationHandler()
	r.GET("/ws", notificationHandler.WebSocketHandler)

	productRepository := ProductRepository.NewProductMongoRepository(db)
	productService := ProductServices.NewProductService(productRepository)
	productHandler := ProductHandlers.NewProductHandler(productService)
	productRoutes := r.Group("/products", requireAuth)
	{
		productRoutes.GET(":id", productHandler.ListProducts)
		productRoutes.GET("available-ingredients", productHandler.AvailableIngredients)
		productRoutes.POST("", productHandler.InsertProduct)
		productRoutes.PUT(":id", productHandler.UpdateProduct)
	}

	billRepository := BillRepository.NewBillMongoRepository(db)
	billService := BillServices.NewBillService(billRepository)
	billHandler := BillHandlers.NewBillHandler(billService)
	billRoutes := r.Group("/bills", requireAuth)
	{
		billRoutes.GET("", billHandler.ListBills)
		billRoutes.GET(":id", billHandler.ListBill)
		billRoutes.PUT(":id", billHandler.UpdateBill)
		billRoutes.POST("/", billHandler.InsertBill)
		billRoutes.PUT("/merge/:idDestination", billHandler.MergeBills)
		billRoutes.GET("/byMonth", billHandler.TotalByMonth)
		billRoutes.POST("upload/xml", billHandler.UploadXML)
	}

	r.GET("/products/:id/bill-items", requireAuth, billHandler.BillItemsByProduct)
	r.GET("/products/recommendations", requireAuth, billHandler.RecommendedProducts)

	marketListRepository := MarketListRepository.NewMarketListMongoRepository(db)
	marketListService := MarketListServices.NewMarketListService(marketListRepository)
	marketListHandler := MarketListHandlers.NewMarketListHandler(marketListService, billService)
	marketListRoutes := r.Group("/market-list", requireAuth)
	{
		marketListRoutes.GET("", marketListHandler.ListMarketLists)
		marketListRoutes.GET(":id", marketListHandler.ListMarketList)
		marketListRoutes.POST("", marketListHandler.InsertMarketList)
		marketListRoutes.GET("suggested", marketListHandler.SuggestMarketList)
		marketListRoutes.PUT(":id/check/:idItem", marketListHandler.MarkItemCheck)
	}
}
