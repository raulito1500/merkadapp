package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raulito1500/merkadapp/products/models"
	"github.com/raulito1500/merkadapp/products/services"
)

type BillHandler struct {
	billService services.BillService
}

func NewBillHandler(bs services.BillService) BillHandler {
	return BillHandler{
		billService: bs,
	}
}

func (bh BillHandler) InsertBill(c *gin.Context) {
	reqBody := new(models.Bill)

	if err := c.Bind(reqBody); err != nil {
		c.JSON(http.StatusBadRequest, "Bad request")
		return
	}
	// TODO: Responde el mensaje como un error y detener
	insertedID, err := bh.billService.Insert(reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, insertedID)
}
