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
	apiv1.GET("/monitors", GetMonitors)
	apiv1.GET("/results/:id", GetResults)

	return &HTTPServer{e}
}

func GetMonitors(c echo.Context) error {
	m, err := GetAllMonitors()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, m) 
}

func GetResults(c echo.Context) error {
	id := c.Param("id")

	i, err := strconv.ParseInt(id, 10, 0)
	if err != nil {
		return err
	}

	r, err := GetResultsByMonitorID(int(i))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, r)	
}
