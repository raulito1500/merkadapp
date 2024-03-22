package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raulito1500/merkadapp/helpers"
	"github.com/raulito1500/merkadapp/src/bill/models"
	"github.com/raulito1500/merkadapp/src/bill/services"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := helpers.ValidateMandatory(reqBody.Where); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf(err.Error(), "Where")})
		return
	}

	for _, i := range reqBody.Items {
		if err := helpers.ValidateMandatory(i.Description); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf(err.Error(), "Item description")})
			return
		}
		if err := helpers.ValidateMandatory(i.Brand); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf(err.Error(), "Item brand")})
			return
		}
		if err := helpers.ValidateFloatNonZeroPositive(i.Quantity); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf(err.Error(), "Item quantity")})
			return
		}
		if err := helpers.ValidateInEnum(i.Unit, models.UNITS); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf(err.Error(), "Item unit")})
			return
		}
		if err := helpers.ValidateFloatNonZeroPositive(i.UnitValue); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf(err.Error(), "Item unit value")})
			return
		}
	}

	insertedID, err := bh.billService.InsertBill(reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": insertedID})
}
