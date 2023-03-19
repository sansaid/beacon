package cmd

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

type response struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func run(port int) {
	defer Beacon.Close()

	e := echo.New()

	e.GET("/health", health)
	e.PUT("/probe", probe)

	// TODO: use an error group so we can exit when Beacon errors
	go Beacon.Start()

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
		r.Error = "Expect namespace and repo query params to be provided"

		return c.JSON(http.StatusBadRequest, r)
	}

	if sc, err := Beacon.Registry().TestRepo(namespace, repo); err != nil {
		r.Message = fmt.Sprintf("Could not fetch repo %s in namespace %s", repo, namespace)
		r.Error = err.Error()

		return c.JSON(sc, r)
	}

	Beacon.StartProbe(namespace, repo, time.Second*20)

	r.Message = fmt.Sprintf("Probe successfully created for repo %s at namespace %s", repo, namespace)
	return c.JSON(http.StatusCreated, r)
}

// CONT: test starting a probe
