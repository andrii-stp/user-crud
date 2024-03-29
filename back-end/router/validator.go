package router

import (
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
	"strings"

	"github.com/andrii-stp/users-crud/model"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// UserValidation example
type UserValidator struct {
	logger    *slog.Logger
	validator *validator.Validate
}

// Validation example
func (uv *UserValidator) Validate(i interface{}) error {
	err := uv.validator.RegisterValidation("status", userStatusValidation)
	if err != nil {
		uv.logger.Error("Error registering custom validation", slog.String("err", err.Error()))
	}

	if err = uv.validator.Struct(i); err != nil {
		var messages []string

		validationErrors := err.(validator.ValidationErrors)

		elem := reflect.TypeOf(&model.User{}).Elem()
		for _, validationError := range validationErrors {
			field, _ := elem.FieldByName(validationError.Field())
			state := "empty"

			if validationError.Tag() != "required" {
				state = "invalid"
			}

			messages = append(messages, fmt.Sprintf(`'%s' is %s`, field.Tag.Get("json"), state))
		}

		return echo.NewHTTPError(http.StatusBadRequest, strings.Join(messages, "\n"))
	}

	return nil
}

func userStatusValidation(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if len(value) == 1 && (value == "I" || value == "A" || value == "T") {
		return true
	}

	return false
}
