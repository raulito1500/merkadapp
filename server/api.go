package server

import (
	"github.com/gin-gonic/gin"
	"github.com/raulito1500/merkadapp/config"
	"github.com/raulito1500/merkadapp/database"
	"github.com/raulito1500/merkadapp/src/handlers"
	"github.com/raulito1500/merkadapp/src/repository"
	"github.com/raulito1500/merkadapp/src/services"
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

	server := gin.Default()
	api.initHandlers(db, server)
	server.Run(config.Port)
}

func (api *Api) initHandlers(db *mongo.Database, r *gin.Engine) {

	productRepository := repository.NewProductMongoRepository(db)
	productService := services.NewProductService(productRepository)
	productHandler := handlers.NewProductHandler(productService)
	productRoutes := r.Group("/products")
	{
		productRoutes.GET("/", productHandler.ListProducts)
		productRoutes.POST("/", productHandler.InsertProduct)
		productRoutes.PUT("/:id", productHandler.UpdateProduct)
	}

	billRepository := repository.NewBillMongoRepository(db)
	billService := services.NewBillService(billRepository)
	billHandler := handlers.NewBillHandler(billService)
	billRoutes := r.Group("/bills")
	{
		billRoutes.POST("/", billHandler.InsertBill)
	}

	marketListRepository := repository.NewMarketListMongoRepository(db)
	marketListService := services.NewMarketListService(marketListRepository)
	marketListHandler := handlers.NewMarketListHandler(marketListService)
	marketListRoutes := r.Group("/market-list")
	{
		marketListRoutes.GET("/", marketListHandler.ListMarketLists)
		marketListRoutes.GET("/:id", marketListHandler.ListMarketList)
		marketListRoutes.POST("/", marketListHandler.InsertMarketList)
		marketListRoutes.GET("/suggested", marketListHandler.SuggestMarketList)
		marketListRoutes.PUT("/:id/check/:idProduct", marketListHandler.MarkItemCheck)
	}
}
