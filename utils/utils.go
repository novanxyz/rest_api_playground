package utils

import (
	"fmt"
	"net/url"
	"novanxyz/models"
	"os"

	"github.com/gin-gonic/gin"
)

func ErrorPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func ResponseError(err error, ctx *gin.Context) {
	if err != nil {
		ctx.Error(err)
		return
	}
}

func ErrorResponseRecovery(ctx *gin.Context, err interface{}) {
	ctx.JSON(500, models.Response{Code: 500, Status: "NOK", Data: err})
}

func ResponseErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		err := ctx.Errors.Last()
		fmt.Println(err)
		if err != nil {
			switch e := err.Err.(type) {
			default:
				ctx.JSON(500, models.Response{Code: 500, Status: "NOK", Data: e.Error()})
			}
			ctx.Abort()
		}
	}
}

func Getenv(key ...string) string {
	if value, ok := os.LookupEnv(key[0]); ok {
		return value
	}
	if len(key) == 1 {
		return ""
	} else {
		return key[1]
	}
}

func QueryParamMap(query url.Values) map[string]interface{} {
	paramMap := make(map[string]interface{}, 0)

	for k, v := range query {
		fmt.Println(k, v)
		if len(v) == 1 && len(v[0]) != 0 {
			paramMap[k] = v[0]
		} else {
			break
		}
	}
	return paramMap
}
