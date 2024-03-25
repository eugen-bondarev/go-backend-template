package util

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetParamString(params *gin.Params, key string) (string, bool) {
	return params.Get(key)
}

func GetParamInt(params *gin.Params, key string) (int, bool) {
	value, ok := params.Get(key)

	if !ok {
		return 0, ok
	}

	converted, err := strconv.Atoi(value)

	if err != nil {
		return 0, false
	}

	return converted, true
}

func DecorateHandler(handler func(*gin.Context) (any, error)) func(*gin.Context) {
	return func(ctx *gin.Context) {
		result, err := handler(ctx)

		if err != nil {
			parsedErr, ok := err.(*RequestError)
			if ok {
				ctx.JSON(parsedErr.StatusCode, gin.H{
					"error": parsedErr.Err.Error(),
				})
				return
			}
			ctx.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(200, result)
	}
}
