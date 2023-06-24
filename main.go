package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"

	"github.com/olahol/melody"
)

func main() {
	e := echo.New()
	m := melody.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// http://localhost:1323
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// ws://localhost:1323
	e.GET("/ws", func(c echo.Context) error {
		m.HandleRequest(c.Response().Writer, c.Request())
		return nil
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.Broadcast(msg)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
