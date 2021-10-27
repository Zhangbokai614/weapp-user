package controller

import (
	"net/http"

	"github.com/dovics/wx-demo/pkg/goods/model"
	"github.com/gin-gonic/gin"
)

func (c *Controller) insertKind(ctx *gin.Context) {
	var req struct {
		KindName string `json:"kind_name,omitempty"`
	}

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	if err := model.InsertKind(c.db, req.KindName); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *Controller) getAllKind(ctx *gin.Context) {
	kinds, err := model.GetAllKind(c.db)
	if err != nil {

		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "kinds": kinds})
}
