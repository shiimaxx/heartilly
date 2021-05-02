package main

import (
	"net/http"
	"strconv"

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
	apiv1.GET("/results/:id", ListResults)

	return &HTTPServer{e}
}

func ListMonitors(c echo.Context) error {
	m, err := GetMonitors()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, m) 
}

func ListResults(c echo.Context) error {
	id := c.Param("id")

	i, err := strconv.ParseInt(id, 10, 0)
	if err != nil {
		return err
	}

	r, err := GetResults(int(i))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, r)	
}
