package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/raulito1500/merkadapp/helpers"
	ibm "github.com/raulito1500/merkadapp/src/bill/models"
	bs "github.com/raulito1500/merkadapp/src/bill/services"
	"github.com/raulito1500/merkadapp/src/market_list/entities"
	imlm "github.com/raulito1500/merkadapp/src/market_list/models"
	mls "github.com/raulito1500/merkadapp/src/market_list/services"
)

type MarketListHandler struct {
	marketListService mls.MarketListService
	billService       bs.BillService
}

func NewMarketListHandler(ms mls.MarketListService, bs bs.BillService) MarketListHandler {
	return MarketListHandler{
		marketListService: ms,
		billService:       bs,
	}
}

func (mh MarketListHandler) ListMarketLists(c *gin.Context) {
	marketlists, err := mh.marketListService.ListMarketLists()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, marketlists)
}

func (mh MarketListHandler) ListMarketList(c *gin.Context) {
	id := c.Param("id")
	// TODO Averiguar porque está usando el modelo de recomendation
	marketlist, err := mh.marketListService.ListMarketList(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, marketlist)
}

func (mh MarketListHandler) InsertMarketList(c *gin.Context) {
	reqBody := new(entities.MarketList)

	if err := c.Bind(reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, i := range reqBody.Items {
		if err := helpers.ValidateMandatory(i.ProductName); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf(err.Error(), "Product name")})
			return
		}
		if err := helpers.ValidateFloatNonZeroPositive(i.Quantity); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf(err.Error(), "Product quantity")})
			return
		}
	}

	insertedID, err := mh.marketListService.InsertMarketList(reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": insertedID})
}
func (mh MarketListHandler) SuggestMarketList(c *gin.Context) {
	marketList := mh.marketListService.SuggestMarketList()
	c.JSON(http.StatusOK, marketList)
}

func (mh MarketListHandler) MarkItemCheck(c *gin.Context) {
	idMarketList := c.Param("id")
	idItem := c.Param("idItem")

	err := mh.marketListService.MarkItemCheck(idMarketList, idItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	marketlist, err := mh.marketListService.ListMarketList(idMarketList)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	foundedItem := new(imlm.ListItemRecommendation)
	for _, item := range marketlist.Items {
		if item.ID == idItem {
			foundedItem = item
		}
	}

	bill := new(ibm.Bill)
	bill.Date = time.Now()
	item := new(ibm.BillItem)
	item.ProductId = foundedItem.ProductId
	item.Description = foundedItem.ProductName
	bill.Items = append(bill.Items, item)
	_, err = mh.billService.InsertBill(bill)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "")
}
