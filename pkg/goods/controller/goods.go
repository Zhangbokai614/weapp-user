package controller

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/dovics/wx-demo/pkg/goods/model"
	"github.com/gin-gonic/gin"
)

// Controller external service interface
type Controller struct {
	db *sql.DB
}

// New create an external service interface
func New(db *sql.DB) *Controller {
	return &Controller{
		db: db,
	}
}

// RegisterRouter register router. It fatal because there is no service if register failed.
func (c *Controller) RegisterRouter(r gin.IRouter) {
	if r == nil {
		log.Fatal("[InitRouter]: server is nil")
	}

	if err := model.CreateDatabase(c.db); err != nil {
		log.Fatal(err)
	}

	if err := model.CreateGoodsTable(c.db); err != nil {
		log.Fatal(err)
	}

	if err := model.CreateKindTable(c.db); err != nil {
		log.Fatal(err)
	}

	r.GET("/info", c.getGoodsInfo)
	r.GET("/info/recommend", c.getRecommendGoodsInfo)
	r.POST("/insert", c.insertGoods)
	r.GET("/kinds", c.getAllKind)
	r.POST("/kind/insert", c.insertKind)
}

func (c *Controller) insertGoods(ctx *gin.Context) {
	var req model.Goods

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	if err := model.InsertGoods(c.db, req); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *Controller) getGoodsInfo(ctx *gin.Context) {
	kindIDStr := ctx.DefaultQuery("kind", "0")
	kindID, err := strconv.Atoi(kindIDStr)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	goodses, err := model.GetGoodsByKind(c.db, uint32(kindID))
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": goodses})
}

func (c *Controller) getRecommendGoodsInfo(ctx *gin.Context) {
	goodses, err := model.GetRecommendGoods(c.db)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": goodses})
}
