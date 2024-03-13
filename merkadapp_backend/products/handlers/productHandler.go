package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raulito1500/merkadapp/products/services"
)

type ProductHandler struct {
	productService services.ProductService
}

func NewProductHandler(ps services.ProductService) ProductHandler {
	return ProductHandler{
		productService: ps,
	}
}

// TODO: Procesar la información de ingreso antes de enviarla a las demás capas
func (ph ProductHandler) ListProducts(c *gin.Context) {
	products := ph.productService.ListProducts()
	c.JSON(http.StatusOK, products)
}

// TODO: Actualizar con data del body
func (ph ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	err := ph.productService.UpdateProduct(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
	}
	c.JSON(http.StatusOK, "")
}
