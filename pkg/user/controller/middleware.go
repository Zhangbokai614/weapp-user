package controller

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/dovics/wx-demo/pkg/user/model"
	"github.com/dovics/wx-demo/util/user"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

//CheckActive middleware that checks the active
func (c *Controller) CheckActive() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		a, err := user.GetID(ctx)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		active, err := model.IsActive(c.db, a)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusConflict, err)
			return
		}

		if !active {
			_ = ctx.AbortWithError(http.StatusLocked, errActive)
			ctx.JSON(http.StatusLocked, gin.H{"status": http.StatusLocked})
			return
		}
	}
}

func (c *Controller) newJWTMiddleware() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test-pet",
		Key:         []byte("moli-tech-cats-member"),
		Timeout:     140 * time.Hour,
		MaxRefresh:  140 * time.Hour,
		IdentityKey: "userID",
		// use data as userID here.
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			return jwt.MapClaims{
				"userID": data,
			}
		},
		// just get the ID
		IdentityHandler: func(ctx *gin.Context) interface{} {
			claims := jwt.ExtractClaims(ctx)
			return claims["userID"]
		},
		Authenticator: func(ctx *gin.Context) (interface{}, error) {
			return c.Login(ctx)
		},
		// no need to check user valid every time.
		Authorizator: func(data interface{}, ctx *gin.Context) bool {
			return true
		},
		Unauthorized: func(ctx *gin.Context, code int, message string) {
			ctx.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: JWT",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
}
