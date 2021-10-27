package controller

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"

	"github.com/dovics/wx-demo/pkg/user/model"
	"github.com/dovics/wx-demo/util/config"
	"github.com/gin-gonic/gin"
)

var (
	errActive          = errors.New("the user is not activated")
	errUserIDNotExists = errors.New("get Admin ID is not exists")
	errUserIDNotValid  = func(value interface{}) error {
		return fmt.Errorf("get Admin ID is not valid. Is %s", value)
	}
)

// Controller external service interface
type Controller struct {
	db     *sql.DB
	JWT    *jwt.GinJWTMiddleware
	client *http.Client
}

// New create an external service interface
func New(db *sql.DB) *Controller {
	c := &Controller{
		db:     db,
		client: http.DefaultClient,
	}
	var err error
	c.JWT, err = c.newJWTMiddleware()
	if err != nil {
		log.Fatal(err)
	}
	return c
}

// RegisterRouter register router. It fatal because there is no service if register failed.
func (c *Controller) RegisterRouter(r gin.IRouter) {
	if r == nil {
		log.Fatal("[InitRouter]: server is nil")
	}
	err := model.CreateDatabase(c.db)
	if err != nil {
		log.Fatal(err)
	}

	err = model.CreateTable(c.db)
	if err != nil {
		log.Fatal(err)
	}

	r.GET("/info", c.getUserInfo)
	r.POST("/modify/active", c.modifyUserActive)
	r.POST("/modify/info", c.modifyUserInfo)
}

//Login JWT validation
func (c *Controller) Login(ctx *gin.Context) (uint32, error) {
	var req struct {
		Code string `json:"code"      binding:"required"`
	}

	err := ctx.ShouldBind(&req)
	if err != nil {
		return 0, err
	}

	resp, err := c.client.Get(BuildWxLoginUrl(req.Code))
	if err != nil {
		return 0, err
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var wx WxResponse
	if err := json.Unmarshal(buf, &wx); err != nil {
		return 0, err
	}

	id, err := model.IsExist(c.db, wx.OpenID)
	if err != sql.ErrNoRows && err != nil {
		return 0, err
	}

	if id == 0 {
		id, err = model.CreateUser(c.db, wx.OpenID, wx.SessionKey)
		if err != nil {
			log.Println("create user fail: ", err)
			return 0, err
		}
	} else {
		if err := model.UpdateSessionKey(c.db, id, wx.SessionKey); err != nil {
			log.Println("update session key fail: ", err)
			return 0, err
		}
	}

	return id, nil
}

type WxResponse struct {
	OpenID     string
	SessionKey string
	UnionID    string
	ErrCode    int
	ErrMsg     string
}

func BuildWxLoginUrl(code string) string {
	appid := config.GetString("wx.appid")
	secret := config.GetString("wx.secret")
	return fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", appid, secret, code)
}

func (c *Controller) modifyUserActive(ctx *gin.Context) {
	var req struct {
		CheckID     uint32 `json:"check_id"    binding:"required"`
		CheckActive bool   `json:"check_active"`
	}

	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = model.ModifyUserActive(c.db, req.CheckID, req.CheckActive)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *Controller) modifyUserInfo(ctx *gin.Context) {
	var req struct {
		NickName string `json:"nick_name,omitempty"`
		Avatar   string `json:"avatar,omitempty"`
		Gender   int    `json:"gender,omitempty"`
		City     string `json:"city,omitempty"`
		Province string `json:"province,omitempty"`
		Country  string `json:"country,omitempty"`
	}

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	id, err := c.GetID(ctx)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	if err := model.ModifyUserInfo(c.db, id, req.NickName, req.Avatar, req.Gender); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *Controller) getUserInfo(ctx *gin.Context) {
	id, err := c.GetID(ctx)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	info, err := model.GetUserInfo(c.db, id)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "info": info})
}
