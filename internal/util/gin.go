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

func DecorateMiddleware(handler func(*gin.Context)) func(*gin.Context) {
	return func(ctx *gin.Context) {
		handler(ctx)
		ctx.Next()
	}
}

func DecorateRequiredMiddleware(handler func(*gin.Context) error) func(*gin.Context) {
	return func(ctx *gin.Context) {
		err := handler(ctx)

		if err == nil {
			ctx.Next()
			return
		}

		ctx.Header("Content-Type", "application/problem+json")
		parsedErr, ok := err.(*RequestError)
		if ok {
			ctx.JSON(parsedErr.StatusCode, gin.H{
				"error": parsedErr.Err.Error(),
			})
			ctx.Abort()
			return
		}

		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		ctx.Abort()
	}
}

func DecorateHandler(handler func(*gin.Context) (any, error)) func(*gin.Context) {
	return func(ctx *gin.Context) {
		result, err := handler(ctx)

		if err != nil {
			parsedErr, ok := err.(*RequestError)
			ctx.Header("Content-Type", "application/problem+json")
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
