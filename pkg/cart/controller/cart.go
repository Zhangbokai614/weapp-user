package controller

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/dovics/wx-demo/pkg/cart/model"
	"github.com/dovics/wx-demo/util/user"
	"github.com/gin-gonic/gin"
)

type CartController struct {
	db *sql.DB
}

func New(db *sql.DB) *CartController {
	return &CartController{
		db: db,
	}
}

func (c *CartController) RegisterRouter(r gin.IRouter) {
	if r == nil {
		log.Fatal("[InitRouter]: server is nil")
	}

	if err := model.CreateDatabase(c.db); err != nil {
		log.Fatal(err)
	}

	if err := model.CreateCartTable(c.db); err != nil {
		log.Fatal(err)
	}

	r.POST("/insert", c.insert)
	r.GET("/info", c.info)
}

func (c *CartController) insert(ctx *gin.Context) {
	var req struct {
		SkuID uint32 `json:"sku_id,omitempty"`
		SpuID uint32 `json:"spu_id,omitempty"`
		Count uint32 `json:"count,omitempty"`
	}

	userID, err := user.GetID(ctx)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	if err := model.InsertCart(c.db, userID, req.SkuID, req.SpuID, req.Count); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *CartController) info(ctx *gin.Context) {
	userID, err := user.GetID(ctx)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	goods, err := model.InfoByUserID(c.db, userID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": goods})
}
