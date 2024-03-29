package router

import (
	"context"
	"log/slog"

	"github.com/andrii-stp/users-crud/handler"
	"github.com/andrii-stp/users-crud/storage"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	swagger "github.com/swaggo/echo-swagger"
)

func Router(logger *slog.Logger, repo storage.UserRepository) *echo.Echo {
	e := echo.New()
	version := e.Group("/api/v1")
	users := version.Group("/users")

	e.Use(middleware.CORS())

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:     true,
		LogURI:        true,
		LogError:      true,
		HandleError:   true,
		LogValuesFunc: logValues(logger),
	}))

	e.Validator = &UserValidator{logger: logger, validator: validator.New(validator.WithRequiredStructEnabled())}

	e.GET("/swagger/*", swagger.WrapHandler)

	userHandler := handler.NewUserHandler(repo)

	users.GET("", userHandler.List)
	users.POST("", userHandler.Create)
	users.PUT("/:id", userHandler.Update)
	users.DELETE("/:id", userHandler.Delete)

	return e
}

func logValues(logger *slog.Logger) func(c echo.Context, v middleware.RequestLoggerValues) error {
	return func(c echo.Context, v middleware.RequestLoggerValues) error {
		if v.Error == nil {
			logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
				slog.String("uri", v.URI),
				slog.Int("status", v.Status),
			)

			return nil
		}

		logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
			slog.String("uri", v.URI),
			slog.Int("status", v.Status),
			slog.String("err", v.Error.Error()),
		)

		return nil
	}
}
