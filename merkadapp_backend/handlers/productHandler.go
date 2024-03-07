package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raulito1500/merkadapp/service"
)

//TODO: Procesar la información de ingreso antes de enviarla a las demás capas
func ListProducts(c *gin.Context) {
	products := service.GetAll()
	c.JSON(http.StatusOK, products)
}

func InsertProduct(c *gin.Context) {
	err := service.UpdateProduct("65e8db1583431fd8de4cb127")
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
	}
	c.JSON(http.StatusOK, "")
}
