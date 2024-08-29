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

func (bh BillHandler) ListBills(c *gin.Context) {
	billList := bh.billService.ListBills()
	c.JSON(http.StatusOK, billList)
}

func (bh BillHandler) ListBill(c *gin.Context) {
	id := c.Param("id")
	// TODO Averiguar porque está usando el modelo de recomendation
	bill, err := bh.billService.ListBill(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bill)
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

func (bh BillHandler) MarkSpentItem(c *gin.Context) {
	idBill := c.Param("id")
	idItem := c.Param("idItem")

	err := bh.billService.MarkSpentItem(idBill, idItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, "")
}

func (bh BillHandler) MergeBills(c *gin.Context) {
	idDestination := c.Param("idDestination")
	var idsOrigen []string

	if err := c.BindJSON(&idsOrigen); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(idsOrigen) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No source IDs provided"})
		return
	}
	err := bh.billService.MergeBills(idDestination, idsOrigen)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, "")
}

func (bh BillHandler) UpdateBill(c *gin.Context) {
	id := c.Param("id")
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

	err := bh.billService.UdpateBill(id, reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}
