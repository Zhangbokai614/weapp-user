package controller

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/dovics/wx-demo/pkg/goods/model"
	"github.com/gin-gonic/gin"
)

// Controller external service interface
type CatagoryController struct {
	db *sql.DB
}

// New create an external service interface
func NewCatagoryController(db *sql.DB) *CatagoryController {
	return &CatagoryController{
		db: db,
	}
}

func (c *CatagoryController) RegisterRouter(r gin.IRouter) {
	if r == nil {
		log.Fatal("[InitRouter]: server is nil")
	}

	if err := model.CreateDatabase(c.db); err != nil {
		log.Fatal(err)
	}

	if err := model.CreateCatagoryTable(c.db); err != nil {
		log.Fatal(err)
	}

	if err := model.CreateSkuTable(c.db); err != nil {
		log.Fatal(err)
	}

	if err := model.CreateSpecTable(c.db); err != nil {
		log.Fatal(err)
	}

	r.GET("/all", c.getAll)
	r.POST("/insert", c.insert)
}

func (c *CatagoryController) insert(ctx *gin.Context) {
	var req struct {
		CatagoryName string `json:"catagory_name,omitempty"`
	}

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	if err := model.InsertCatagory(c.db, req.CatagoryName); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *CatagoryController) getAll(ctx *gin.Context) {
	catagorys, err := model.InfoAllCatagory(c.db)
	if err != nil {

		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "catagorys": catagorys})
}
