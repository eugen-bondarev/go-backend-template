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
)

type App struct {
	signingSvc           svc.ISigningSvc
	userRepo             repo.IUserRepo
	userDataSigningSvc   svc.UserDataSigningSvc
	forgotPassSigningSvc svc.ForgotPassSigningSvc
	mailerSvc            svc.IMailerSvc
	authSvc              svc.IAuthSvc
	tokenInvalidator     svc.ITokenInvalidatorSvc
	policies             permissions.Policies
	// tmpStorageSvc        svc.ITmpStorageSvc
}

func NewApp() (App, error) {
	pg, err := postgres.NewPostgres(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_PORT"),
	)

	if err != nil {
		return App{}, err
	}

	err = pg.Migrate("./assets/migrations")

	if err != nil {
		return App{}, err
	}

	userRepo := repo.NewPGUserRepo(&pg)
	mailerSvc := svc.NewSMTPMailerSvc(
		os.Getenv("SMTP_USERNAME"),
		os.Getenv("SMTP_PASSWORD"),
		os.Getenv("SMTP_HOST"),
		os.Getenv("SMTP_PORT"),
	)
	authSvc := svc.NewDefaultAuthSvc(userRepo, os.Getenv("PEPPER"))
	signingSvc := svc.NewJWTSigningSvc(os.Getenv("JWT_SECRET"))
	userDataSigningSvc := svc.NewUserDataSigningSvc(signingSvc)
	forgotPassSigningSvc := svc.NewForgotPassSigningSvc(signingSvc)

	redis := redis.NewRedis(
		os.Getenv("REDIS_HOST"),
		os.Getenv("REDIS_PORT"),
		os.Getenv("REDIS_PASS"),
	)

	tmpStorageSvc := svc.NewRedisTempStorageSvc(&redis)
	tokenInvalidator := svc.NewDefaultTokenInvalidator(tmpStorageSvc)

	policies := permissions.NewPolicies()
	policies.Add("admin", "index", "users")
	policies.Add("admin", "manage", "users")

	return App{
		signingSvc:           signingSvc,
		userRepo:             userRepo,
		userDataSigningSvc:   userDataSigningSvc,
		forgotPassSigningSvc: forgotPassSigningSvc,
		tokenInvalidator:     tokenInvalidator,
		mailerSvc:            mailerSvc,
		authSvc:              authSvc,
		policies:             policies,
	}, nil
}

func main() {
	godotenv.Load()

	app, err := NewApp()

	if err != nil {
		panic(err.Error())
	}

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	mw := middleware.NewGinMiddlewareFactory(
		app.userDataSigningSvc,
		app.tokenInvalidator,
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
			return app.userRepo.GetUsers()
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

			return nil, app.userRepo.DeleteUserByID(id)
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

			token, refreshToken, err := app.refreshToken(payload.RefreshToken)

			if err != nil {
				return nil, err
			}

			return struct {
				Token        string `json:"token"`
				RefreshToken string `json:"refreshToken"`
			}{
				Token:        token,
				RefreshToken: refreshToken,
			}, nil
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

			parsedSessionToken, err := app.userDataSigningSvc.ParseSessionToken(payload.Token)

			if err != nil {
				return nil, err
			}

			parsedRefreshToken, err := app.userDataSigningSvc.ParseRefreshToken(payload.RefreshToken)

			if err != nil {
				return nil, err
			}

			app.tokenInvalidator.Invalidate(payload.Token, parsedSessionToken.ExpiresAt)
			app.tokenInvalidator.Invalidate(payload.RefreshToken, parsedRefreshToken.ExpiresAt)

			return nil, nil
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

			token, refreshToken, err := app.login(payload.Email, payload.Password)

			if err != nil {
				return nil, err
			}

			return struct {
				Token        string `json:"token"`
				RefreshToken string `json:"refreshToken"`
			}{
				Token:        token,
				RefreshToken: refreshToken,
			}, nil
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

			return nil, app.register(payload.Email, payload.Password)
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

			email, err := app.forgotPassSigningSvc.Parse(payload.Token)

			if err != nil {
				return nil, err
			}

			return nil, app.authSvc.SetPasswordByEmail(email, payload.Password)
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

			token, err := app.forgotPassSigningSvc.Sign(payload.Email)

			if err != nil {
				return nil, err
			}

			mail := svc.NewMailBuilder(
				payload.Email,
				"So you want to reset your password?\n"+
					"Your token is: "+token.Value,
			)

			return nil, app.mailerSvc.Send(mail)
		}),
	)

	r.Run("0.0.0.0:4200")
}
