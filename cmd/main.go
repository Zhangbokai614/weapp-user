package main

import (
	"database/sql"
	"fmt"
	"log"

	c "github.com/dovics/wx-demo/config"
	cart "github.com/dovics/wx-demo/pkg/cart/controller"
	goods "github.com/dovics/wx-demo/pkg/goods/controller"
	user "github.com/dovics/wx-demo/pkg/user/controller"

	"github.com/dovics/wx-demo/util/config"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	c.Initialize()
}

var (
	userRouterGroup        = "/api/v1/user"
	spuRouterGroup         = "/api/v1/spu"
	categoryRouterGroup    = "/api/v1/category"
	cartRouterGroup        = "/api/v1/cart"
	userRouterGroupLogin   = userRouterGroup + "/login"
	userRouterRefreshToken = userRouterGroup + "/refresh_token"
)

func main() {
	var (
		host     = config.GetString("database.mysql.host")
		port     = config.GetString("database.mysql.port")
		database = config.GetString("database.mysql.database")
		username = config.GetString("database.mysql.username")
		password = config.GetString("database.mysql.password")
		charset  = config.GetString("database.mysql.charset")
	)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%t&loc=%s",
		username, password, host, port, database, charset, true, "Local")

	dbConn, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	router := gin.Default()

	userController := user.New(dbConn)
	spuController := goods.NewSpuController(dbConn)
	categoryController := goods.NewCatagoryController(dbConn)
	cartController := cart.New(dbConn)
	router.POST(userRouterGroupLogin, userController.JWT.LoginHandler)
	router.POST(userRouterRefreshToken, userController.JWT.RefreshHandler)

	router.Use(userController.JWT.MiddlewareFunc())
	router.Use(userController.CheckActive())
	userController.RegisterRouter(router.Group(userRouterGroup))
	spuController.RegisterRouter(router.Group(spuRouterGroup))
	categoryController.RegisterRouter(router.Group(categoryRouterGroup))
	cartController.RegisterRouter(router.Group(cartRouterGroup))

	fmt.Println("port" + config.GetString("app.port"))
	log.Fatal(router.Run("0.0.0.0:" + config.GetString("app.port")))
}
