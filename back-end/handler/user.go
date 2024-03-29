package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/andrii-stp/users-crud/model"
	"github.com/andrii-stp/users-crud/storage"
	"github.com/labstack/echo/v4"
)

// UserHandler example
type UserHandler struct {
	repository storage.UserRepository
}

// NewUserHandler example
func NewUserHandler(repository storage.UserRepository) *UserHandler {
	return &UserHandler{repository: repository}
}

// List godoc
//
//	@Summary		List users
//	@Description	get users
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		model.User
//	@Failure		500	{object}	echo.HTTPError
//	@Router			/users [get]
func (u UserHandler) List(c echo.Context) error {
	logger := c.Logger()

	users, err := u.repository.List(c.Request().Context())
	if err != nil {
		logger.Errorf("failed to get users from database: %v", err)

		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get users")
	}

	return c.JSON(http.StatusOK, users)
}

// Create godoc
//
//	@Summary		Create user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		model.User		true	"Create user"
//	@Success		201		{object}	model.User
//	@Failure		400		{object}	echo.HTTPError
//	@Failure		409		{object}	echo.HTTPError
//	@Failure		500		{object}	echo.HTTPError
//	@Router			/users [post]
func (u UserHandler) Create(c echo.Context) error {
	logger := c.Logger()

	var user model.User
	if err := c.Bind(&user); err != nil {
		logger.Errorf("failed to bind to user type: %v", err)

		return echo.NewHTTPError(http.StatusBadRequest, "Failed to bind request body")
	}

	if err := c.Validate(user); err != nil {
		return err
	}

	if err := u.repository.Create(c.Request().Context(), &user); err != nil {
		logger.Errorf("failed to create user: %v", err)

		if errors.Is(err, storage.ErrAlreadyExist) {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
	}

	return c.JSON(http.StatusCreated, user)
}

// Update godoc
//
//	@Summary		Update user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int				true	"Update user"
//	@Param			user	body		model.User		true	"Update user"
//	@Success		201		{object}	model.User
//	@Failure		400		{object}	echo.HTTPError
//	@Failure		409		{object}	echo.HTTPError
//	@Failure		500		{object}	echo.HTTPError
//	@Router			/users{id} [put]
func (u UserHandler) Update(c echo.Context) error {
	logger := c.Logger()
	idParam := c.Param("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		logger.Errorf("failed to convert id to int: %v", err)

		return echo.NewHTTPError(http.StatusBadRequest, `'id' is not a number`)
	}

	var user model.User
	if err := c.Bind(&user); err != nil {
		logger.Errorf("failed to bind to user type: %v", err)

		return echo.NewHTTPError(http.StatusBadRequest, "Failed to bind request body")
	}

	if err := c.Validate(user); err != nil {
		return err
	}

	err = u.repository.Update(c.Request().Context(), id, &user)
	if err != nil {
		logger.Errorf("failed to update user: %v", err)

		if errors.Is(err, storage.ErrAlreadyExist) {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}

		if errors.Is(err, storage.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update user")
	}

	return c.JSON(http.StatusOK, user)
}

// Delete godoc
//
//	@Summary		Delete user
//	@Description	Delete by username
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"	Format(int64)
//	@Success		204	{object}	model.User
//	@Failure		400	{object}	echo.HTTPError
//	@Failure		404	{object}	echo.HTTPError
//	@Failure		500	{object}	echo.HTTPError
//	@Router			/users/{id} [delete]
func (u UserHandler) Delete(c echo.Context) error {
	logger := c.Logger()
	idParam := c.Param("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		logger.Errorf("failed to convert id to int: %v", err)

		return echo.NewHTTPError(http.StatusBadRequest, `'id' is not a number`)
	}

	if err := u.repository.Delete(c.Request().Context(), id); err != nil {
		logger.Errorf("failed to delete user: %v", err)

		if errors.Is(err, storage.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete user")
	}

	return c.NoContent(http.StatusNoContent)
}
