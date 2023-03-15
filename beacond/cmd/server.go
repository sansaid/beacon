package cmd

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

type base struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func run(beacon *Beacon, port int) {
	e := echo.New()

	e.GET("/health", health)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

func health(c echo.Context) error {
	var b base

	b.Message = "beacond is happily running :)"

	return c.JSON(http.StatusOK, b)
}
