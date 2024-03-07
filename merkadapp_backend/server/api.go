package server

import (
	"github.com/gin-gonic/gin"
	"github.com/raulito1500/merkadapp/handlers"
)

type Api struct {
}

func NewApi() *Api {
	return &Api{}
}

func (api *Api) Run() {
	router := gin.Default()
	api.configRoutes(router)
	router.Run()
}

func (api *Api) configRoutes(r *gin.Engine) {
	productRoutes := r.Group("/products")
	{
		productRoutes.GET("/", func(c *gin.Context) { handlers.ListProducts(c) })
		productRoutes.PUT("/:id", func(c *gin.Context) { handlers.InsertProduct(c) })
	}
}
