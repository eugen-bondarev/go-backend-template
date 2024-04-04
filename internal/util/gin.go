package util

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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
			ctx.Header("Content-Type", "application/problem+json")
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

		if result == nil {
			return
		}

		ctx.JSON(200, result)
	}
}

var validation = validator.New()

func GinGetBody[TOut any](ctx *gin.Context) (TOut, error) {
	var payload TOut
	ctx.ShouldBindBodyWith(&payload, binding.JSON)
	err := validation.Struct(&payload)

	if err == nil {
		return payload, nil
	}

	validationErrors := err.(validator.ValidationErrors)

	return payload, fmt.Errorf("field '%s' is invalid", validationErrors[0].Field())
}

func GinHealthz(r gin.IRouter) {
	r.GET("/healthz", func(ctx *gin.Context) {
		ctx.JSON(200, "ok")
	})
}

func GinConfigureCors(commaSeparatedOrigins string) gin.HandlerFunc {
	config := cors.DefaultConfig()
	if len(commaSeparatedOrigins) > 0 {
		config.AllowOrigins = strings.Split(commaSeparatedOrigins, ",")
	} else {
		config.AllowAllOrigins = true
	}
	config.AllowMethods = []string{"GET", "PUT", "POST", "PATCH", "DELETE", "OPTIONS", "HEAD"}
	fmt.Println(config)
	return cors.New(config)
}
