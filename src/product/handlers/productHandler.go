package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raulito1500/merkadapp/helpers"
	"github.com/raulito1500/merkadapp/src/product/models"
	"github.com/raulito1500/merkadapp/src/product/services"
)

type ProductHandler struct {
	productService services.ProductService
}

func NewProductHandler(ps services.ProductService) ProductHandler {
	return ProductHandler{
		productService: ps,
	}
}

func (ph ProductHandler) ListProducts(c *gin.Context) {
	products := ph.productService.ListProducts()
	c.JSON(http.StatusOK, products)
}

func (ph ProductHandler) InsertProduct(c *gin.Context) {
	reqBody := new(models.Product)

	if err := c.Bind(reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := helpers.ValidateMandatory(reqBody.Category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf(err.Error(), "Category")})
		return
	}
	if err := helpers.ValidateInEnum(reqBody.Category, models.CATEGORIES); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf(err.Error(), "Category")})
		return
	}
	if err := helpers.ValidateMandatory(reqBody.Name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf(err.Error(), "Name")})
		return
	}
	if err := helpers.ValidateMandatory(reqBody.Repeat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf(err.Error(), "Repeat")})
		return
	}
	if err := helpers.ValidateIntNonZeroPositive(reqBody.Quantity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf(err.Error(), "Quantity")})
		return
	}

	insertedID, err := ph.productService.InsertProduct(reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, insertedID)
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
