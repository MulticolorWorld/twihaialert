package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type ErrorHandler struct {
}

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

func (eh ErrorHandler) WrongToken(c echo.Context) error {
	return c.Render(http.StatusBadRequest, "wrongToken", nil)
}

func (eh ErrorHandler) AccountAlreadyExist(c echo.Context) error {
	return c.Render(http.StatusBadRequest, "accountAlreadyExist", nil)
}
