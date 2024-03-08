package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raulito1500/merkadapp/products/services"
)

// TODO: Procesar la información de ingreso antes de enviarla a las demás capas
func ListProducts(c *gin.Context) {
	products := services.GetAll()
	c.JSON(http.StatusOK, products)
}

// TODO: Actualizar con data del body
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	err := services.UpdateProduct(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
	}
	c.JSON(http.StatusOK, "")
}
