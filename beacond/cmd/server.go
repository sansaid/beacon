package cmd

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

type response struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func run(beacon *Beacon, port int) {
	defer beacon.Close()

	e := echo.New()

	e.GET("/health", health)
	e.PUT("/probe", probe)

	go beacon.Start()

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

func health(c echo.Context) error {
	var r response

	r.Message = "beacond is happily running :)"

	return c.JSON(http.StatusOK, r)
}

func probe(c echo.Context) error {
	var r response

	namespace := c.QueryParam("namespace")
	repo := c.QueryParam("repo")

	if namespace == "" || repo == "" {
		r.Error = "Expect namepsace and repo query params to be provided"

		return c.JSON(http.StatusBadRequest, r)
	}

	// TODO: add some awareness to the global beacon instance
	r.Message = "Probe created for "
	return c.JSON(http.StatusOK, r)
}
