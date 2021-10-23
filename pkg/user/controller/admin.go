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

	"github.com/dovics/wx-demo/pkg/user/model/mysql"
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
	err := mysql.CreateDatabase(c.db)
	if err != nil {
		log.Fatal(err)
	}

	err = mysql.CreateTable(c.db)
	if err != nil {
		log.Fatal(err)
	}

	r.POST("/modify/active", c.modifyAdminActive)
}

func (c *Controller) modifyAdminActive(ctx *gin.Context) {
	var (
		admin struct {
			CheckID     uint32 `json:"check_id"    binding:"required"`
			CheckActive bool   `json:"check_active"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.ModifyAdminActive(c.db, admin.CheckID, admin.CheckActive)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
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

	id, err := mysql.IsExist(c.db, wx.OpenID)
	if err != sql.ErrNoRows && err != nil {
		return 0, err
	}

	if id == 0 {
		id, err = mysql.CreateUser(c.db, wx.OpenID, wx.SessionKey)
		if err != nil {
			fmt.Println(1)
			return 0, err
		}
	} else {
		if err := mysql.UpdateSessionKey(c.db, id, wx.SessionKey); err != nil {
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
