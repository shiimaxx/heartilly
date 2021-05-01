package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type HTTPServer struct {
	*echo.Echo
}

func NewHTTPServer() *HTTPServer {
	e := echo.New()
	e.Use(middleware.Recover())

	apiv1 := e.Group("/api/v1")
	apiv1.GET("/monitors", ListMonitors)

	return &HTTPServer{e}
}

func ListMonitors(c echo.Context) error {
	m, err := GetMonitors()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, m) 
}
