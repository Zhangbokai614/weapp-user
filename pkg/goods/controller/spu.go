package controller

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/dovics/wx-demo/pkg/goods/model"
	"github.com/gin-gonic/gin"
)

// Controller external service interface
type SpuController struct {
	db *sql.DB
}

// New create an external service interface
func NewSpuController(db *sql.DB) *SpuController {
	return &SpuController{
		db: db,
	}
}

// RegisterRouter register router. It fatal because there is no service if register failed.
func (c *SpuController) RegisterRouter(r gin.IRouter) {
	if r == nil {
		log.Fatal("[InitRouter]: server is nil")
	}

	if err := model.CreateDatabase(c.db); err != nil {
		log.Fatal(err)
	}

	if err := model.CreateSpuTable(c.db); err != nil {
		log.Fatal(err)
	}

	r.GET("/info", c.getSpuInfoByKind)
	r.GET("/info/recommend", c.getRecommendSpuInfo)
	r.POST("/insert", c.insertSpu)
	r.GET("/info/detail", c.getSpuInfoDetail)
}

func (c *SpuController) insertSpu(ctx *gin.Context) {
	var req model.Spu

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	tx, err := c.db.Begin()
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}
	defer tx.Rollback()

	spuID, err := model.TxInsertSpu(tx, req)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	for _, spec := range req.Spec {
		if err := model.InsertSpec(c.db, spuID, spec.Kind, spec.Value); err != nil {
			ctx.Error(err)
			ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *SpuController) getSpuInfoByKind(ctx *gin.Context) {
	catagoryIDStr, ok := ctx.GetQuery("catagory")
	if !ok {
		ctx.Error(errors.New("request should contain catagory id"))
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}
	catagoryID, err := strconv.Atoi(catagoryIDStr)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	spus, err := model.GetSpuByCatagory(c.db, uint32(catagoryID))
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": spus})
}

func (c *SpuController) getRecommendSpuInfo(ctx *gin.Context) {
	spus, err := model.GetRecommendSpu(c.db)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": spus})
}

func (c *SpuController) getSpuInfoDetail(ctx *gin.Context) {
	spuIDStr, ok := ctx.GetQuery("spu_id")
	if !ok {
		ctx.Error(errors.New("request should contain spu id"))
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}
	spuID, err := strconv.Atoi(spuIDStr)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	tx, err := c.db.Begin()
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	defer tx.Rollback()

	spu, err := model.TxInfoSpuByID(tx, uint32(spuID))
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	spu.Spec, err = model.TxInfoSpecBySpuID(tx, uint32(spuID))
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	if err := tx.Commit(); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": spu})
}
