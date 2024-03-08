package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raulito1500/merkadapp/products/models"
)

func InsertBill(c *gin.Context) {
	reqBody := new(models.Bill)

	if err := c.Bind(reqBody); err != nil {
		c.JSON(http.StatusBadRequest, "Bad request")
	}
	c.JSON(http.StatusOK, reqBody)

}
