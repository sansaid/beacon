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
	e.GET("/probes", listProbes)
	e.GET("/beacon", getBeaconDetails)
	e.POST("/probe", createProbe)
	e.DELETE("/probe", deleteProbe)

	org := NewOrGroup()

	org.Go(Beacon.Start)
	org.Go(func() error { return e.Start(fmt.Sprintf(":%d", port)) })

	e.Logger.Fatal(org.Wait())
}

func health(c echo.Context) error {
	var r response

	r.Message = "beacond is happily running :)"

	return c.JSON(http.StatusOK, r)
}

func deleteProbe(c echo.Context) error {
	var r response

	namespace := c.QueryParam("namespace")
	repo := c.QueryParam("repo")

	if namespace == "" || repo == "" {
		r.Error = "Expect namespace and repo query params to be provided"

		return c.JSON(http.StatusBadRequest, r)
	}

	err := Beacon.StopProbe(namespace, repo, time.Second*20)

	if _, ok := err.(BeaconErrorProbeDoesNotExist); ok {
		r.Error = err.Error()
		r.Message = fmt.Sprintf("Probe not found for repo %s at namespace %s", repo, namespace)

		return c.JSON(http.StatusNotFound, r)
	}

	r.Message = fmt.Sprintf("Probe successfully deleted for repo %s at namespace %s", repo, namespace)
	return c.JSON(http.StatusCreated, r)
}

func createProbe(c echo.Context) error {
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

	err := Beacon.StartProbe(namespace, repo, time.Second*20)

	if _, ok := err.(BeaconErrorProbeAlreadyExists); ok {
		r.Error = err.Error()
		r.Message = fmt.Sprintf("Probe already exists for repo %s at namespace %s", repo, namespace)

		return c.JSON(http.StatusConflict, r)
	}

	r.Message = fmt.Sprintf("Probe successfully created for repo %s at namespace %s", repo, namespace)
	return c.JSON(http.StatusCreated, r)
}

func listProbes(c echo.Context) error {
	var r struct {
		Probes []string `json:"probes"`
	}

	r.Probes = Beacon.ListProbes()

	return c.JSON(http.StatusConflict, r)
}

func getBeaconDetails(c echo.Context) error {
	var r struct {
		Registry string   `json:"registry"`
		Probes   []string `json:"probes"`
		Runtime  string   `json:"runtime"`
	}

	r.Registry = Beacon.Registry().URL()
	r.Probes = Beacon.ListProbes()
	r.Runtime = string(Beacon.Runtime().Type())

	return c.JSON(http.StatusConflict, r)
}
