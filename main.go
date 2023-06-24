package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"

	"github.com/olahol/melody"
)

func main() {
	e := echo.New()
	m := melody.New()
	m.Config.MaxMessageSize = 8192

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

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
		log.Println("message", string(msg))

		m.BroadcastFilter(msg, func(q *melody.Session) bool {
			// Retorna true para enviar a mensagem para todos os clientes, exceto o remetente original
			return q != s
		})
	})

	e.Logger.Fatal(e.Start(":1323"))
}
