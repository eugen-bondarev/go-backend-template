package main

import (
	"go-backend-template/internal/dto"
	"go-backend-template/internal/middleware"
	"go-backend-template/internal/permissions"
	"go-backend-template/internal/postgres"
	"go-backend-template/internal/redis"
	"go-backend-template/internal/repo"
	"go-backend-template/internal/svc"
	"go-backend-template/internal/util"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/text/language"
)

type App struct {
	userRepo          repo.IUserRepo
	signing           svc.ISigning
	userDataSigning   svc.UserDataSigning
	forgotPassSigning svc.ForgotPassSigning
	mailer            svc.IMailer
	auth              svc.IAuth
	fieManager        svc.IFileManager
	policies          permissions.Policies
}

type Controller struct {
	app *App
}

func MustInitApp() App {
	pg, err := postgres.NewPostgres(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_PORT"),
	)
	util.PanicOnError(err)

	err = pg.Migrate("./assets/migrations")
	util.PanicOnError(err)

	userRepo := repo.NewPGUserRepo(&pg)
	mailer := svc.NewSMTPMailer(
		os.Getenv("SMTP_USERNAME"),
		os.Getenv("SMTP_PASSWORD"),
		os.Getenv("SMTP_HOST"),
		os.Getenv("SMTP_PORT"),
	)
	auth := svc.NewDefaultAuth(userRepo, os.Getenv("PEPPER"))
	signing := svc.NewJWTSigning(os.Getenv("JWT_SECRET"))
	forgotPassSigning := svc.NewForgotPassSigning(signing)

	redis, redisErr := redis.NewRedis(
		os.Getenv("REDIS_HOST"),
		os.Getenv("REDIS_PORT"),
		os.Getenv("REDIS_PASS"),
	)

	var tokenInvalidator svc.ITokenInvalidator
	if redisErr != nil {
		tokenInvalidator = svc.NewNoopTokenInvalidator()
	} else {
		tmpStorage := svc.NewRedisTempStorage(&redis)
		tokenInvalidator = svc.NewDefaultTokenInvalidator(tmpStorage)
	}

	userDataSigning := svc.NewUserDataSigning(signing, tokenInvalidator)

	fileRepo := repo.NewPGFileRepo(&pg)
	fileStorage := svc.NewDiskFileStorage("./storage")
	fileManager := svc.NewFileManager(fileRepo, fileStorage)

	policies := permissions.NewPolicies()
	policies.Add("admin", "index", "users")
	policies.Add("admin", "manage", "users")

	return App{
		signing:           signing,
		userRepo:          userRepo,
		userDataSigning:   userDataSigning,
		forgotPassSigning: forgotPassSigning,
		mailer:            mailer,
		auth:              auth,
		fieManager:        fileManager,
		policies:          policies,
	}
}

func main() {
	godotenv.Load()

	app := MustInitApp()

	util.Bundle = i18n.NewBundle(language.English)
	util.Bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	util.Bundle.LoadMessageFile("de.toml")
	util.Localizers = make(map[string]*i18n.Localizer)

	util.Localizers["de"] = i18n.NewLocalizer(util.Bundle, "de")
	util.Localizers["en"] = i18n.NewLocalizer(util.Bundle, "en")

	// msg := &i18n.LocalizeConfig{
	// 	DefaultMessage: &i18n.Message{
	// 		ID:    "greeting1",
	// 		Other: "Hello, {{.Name}}",
	// 	},
	// 	TemplateData: map[string]any{
	// 		"Name": "Diana",
	// 	},
	// }

	// fmt.Println(de.Localize(msg))
	// fmt.Println(en.Localize(msg))

	controller := Controller{app: &app}

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	mw := middleware.NewGinMiddlewareFactory(
		app.userDataSigning,
		&app.policies,
	)

	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.Use(util.GinConfigureCors(os.Getenv("CORS_ALLOW_ORIGINS")))

	util.GinHealthz(r)

	v1 := r.Group("/v1")

	v1.GET(
		"/users",
		mw.SetRole(),
		mw.EnforcePolicy("index", "users"),
		util.DecorateHandler(func(ctx *gin.Context) (any, error) {
			return controller.getUsers()
		}),
	)

	v1.DELETE(
		"/users/:id",
		mw.SetRole(),
		mw.EnforcePolicy("manage", "users"),
		util.DecorateHandler(func(ctx *gin.Context) (any, error) {
			id, err := strconv.Atoi(ctx.Params.ByName("id"))

			if err != nil {
				return nil, err
			}

			return controller.deleteUserByID(id)
		}),
	)

	auth := v1.Group("/auth")

	auth.POST(
		"/refresh",
		util.DecorateHandler(func(ctx *gin.Context) (any, error) {
			payload, err := util.GinGetBody[struct {
				dto.WithRefreshToken
			}](ctx)

			if err != nil {
				return nil, err
			}

			return controller.refreshToken(payload.RefreshToken)
		}),
	)

	auth.POST(
		"/logout",
		util.DecorateHandler(func(ctx *gin.Context) (any, error) {
			payload, err := util.GinGetBody[struct {
				dto.WithToken
				dto.WithRefreshToken
			}](ctx)

			if err != nil {
				return nil, err
			}

			return controller.logout(payload.Token, payload.RefreshToken)
		}),
	)

	auth.POST(
		"/login",
		util.DecorateHandler(func(ctx *gin.Context) (any, error) {
			payload, err := util.GinGetBody[struct {
				dto.WithEmail
				dto.WithPassword
			}](ctx)

			if err != nil {
				return nil, err
			}

			return controller.login(payload.Email, payload.Password)
		}),
	)

	auth.POST(
		"/register",
		util.DecorateHandler(func(ctx *gin.Context) (any, error) {
			payload, err := util.GinGetBody[struct {
				dto.WithEmail
				dto.WithPassword
			}](ctx)

			if err != nil {
				return nil, err
			}

			return controller.register(payload.Email, payload.Password)
		}),
	)

	auth.POST(
		"/reset-password",
		util.DecorateHandler(func(ctx *gin.Context) (any, error) {
			payload, err := util.GinGetBody[struct {
				dto.WithToken
				dto.WithPassword
			}](ctx)

			if err != nil {
				return nil, err
			}

			return controller.resetPassword(payload.Token, payload.Password)
		}),
	)

	auth.POST(
		"/forgot-password",
		util.DecorateHandler(func(ctx *gin.Context) (any, error) {
			payload, err := util.GinGetBody[struct {
				dto.WithEmail
			}](ctx)

			if err != nil {
				return nil, err
			}

			return controller.forgotPassword(payload.Email)
		}),
	)

	r.Run("0.0.0.0:4200")
}
