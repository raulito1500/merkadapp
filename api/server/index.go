package server

import (
	"github.com/gin-gonic/gin"
	"github.com/raulito1500/merkadapp/database"
	"github.com/raulito1500/merkadapp/products/handlers"
	"github.com/raulito1500/merkadapp/products/repository"
	"github.com/raulito1500/merkadapp/products/services"
	"go.mongodb.org/mongo-driver/mongo"
)

type Api struct {
	db *mongo.Database
}

func NewApi() *Api {
	return &Api{}
}

func (api *Api) Run() {
	router := gin.Default()
	db := database.NewMongoDatabase().GetDb()
	api.initHandlers(db, router)
	router.Run()
}
func (api *Api) initHandlers(db *mongo.Database, r *gin.Engine) {

	productRepository := repository.NewProductMongoRepository(db)
	productService := services.NewProductService(productRepository)
	productHandler := handlers.NewProductHandler(productService)

	billRepository := repository.NewBillMongoRepository(db)
	billService := services.NewBillService(billRepository)
	billHandler := handlers.NewBillHandler(billService)

	marketListRepository := repository.NewMarketListMongoRepository(db)
	marketListService := services.NewMarketListService(marketListRepository)
	marketListHandler := handlers.NewMarketListHandler(marketListService)

	productRoutes := r.Group("/products")
	{
		productRoutes.GET("/", productHandler.ListProducts)
		productRoutes.POST("/", productHandler.InsertProduct)
		productRoutes.PUT("/:id", productHandler.UpdateProduct)
	}
	billRoutes := r.Group("/bills")
	{
		billRoutes.POST("/", billHandler.InsertBill)
	}
	marketListRoutes := r.Group("/market-list")
	{
		marketListRoutes.POST("/", marketListHandler.InsertMarketList)
		marketListRoutes.GET("/suggested", marketListHandler.SuggestMarketList)
		marketListRoutes.PUT("/:id/check/:idProduct", marketListHandler.MarkCheck)
	}
}
