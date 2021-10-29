package user

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

var (
	errUserIDNotExists = errors.New("get Admin ID is not exists")
	errUserIDNotValid  = func(value interface{}) error {
		return fmt.Errorf("get Admin ID is not valid. Is %s", value)
	}
)

func GetID(ctx *gin.Context) (uint32, error) {
	id, ok := ctx.Get("userID")
	if !ok {
		return 0, errUserIDNotExists
	}

	v, ok := id.(float64)
	if !ok {
		return 0, errUserIDNotValid(id)
	}

	return uint32(v), nil
}
