package server

import (
	"beacon/beacond/oci"
	"beacon/beacond/registry"
	"fmt"
	"net/http"
	"time"

	"beacon/beacond/models"

	"github.com/labstack/echo"
)

// @Title				beacond API
// @Version			0.1
// @Description	API for beacond server
func Run(ociClient oci.OCIRuntime, registryClient registry.Registry, port int, cleanOnExit bool) {
	NewBeacon(ociClient, registryClient, cleanOnExit)
	defer Beacon.Close()

	e := echo.New()

	e.GET("/health", health)

	e.GET("/beacon", getBeaconDetails)

	e.GET("/probes", listProbes)
	e.POST("/probe", createProbe)
	e.DELETE("/probe", deleteProbe)

	org := NewOrGroup()

	org.Go(Beacon.Start)
	org.Go(func() error { return e.Start(fmt.Sprintf(":%d", port)) })

	e.Logger.Fatal(org.Wait())
}

// health handles the GET /health method for beacond
//
//	@Summary		Health check
//	@Description	reports the health of the beacond server
//	@Produce		json
//	@Success		200	{object}	BaseResponse
//	@Router			/health [get]
func health(c echo.Context) error {
	var r models.ServerBaseResponse

	r.Message = "beacond is happily running :)"

	return c.JSON(http.StatusOK, r)
}

// deleteProbe handles the DELETE /probe method for beacond
//
//	@Summary		Delete a probe
//	@Description	deletes the probe for the namespace and repo provided in the URL query parameters
//	@Produce		json
//	@Param			namespace	query		string	true	"the repo namespace the probe should check for image updates"
//	@Param			repo		query		string	true	"the repo name which the probe should check for image updates"
//	@Success		201			{object}	BaseResponse
//	@Failure		404			{object}	BaseResponse
//	@Failure		400			{object}	BaseResponse
//	@Failure		500			{object}	BaseResponse
//	@Router			/probe [delete]
func deleteProbe(c echo.Context) error {
	var r models.ServerBaseResponse

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

	if err != nil {
		r.Error = err.Error()
		r.Message = fmt.Sprintf("Failed to create probe for repo %s at namespace %s", repo, namespace)

		return c.JSON(http.StatusInternalServerError, r)
	}

	r.Message = fmt.Sprintf("Probe successfully deleted for repo %s at namespace %s", repo, namespace)
	return c.JSON(http.StatusCreated, r)
}

// createProbe handles the POST /probe method for beacond
//
//	@Summary		Create a probe
//	@Description	creates a probe for the namespace and repo provided in the URL query parameters
//	@Produce		json
//	@Param			namespace	query		string	true	"the repo namespace the probe should check for image updates"
//	@Param			repo		query		string	true	"the repo name which the probe should check for image updates"
//	@Success		201			{object}	BaseResponse
//	@Failure		409			{object}	BaseResponse
//	@Failure		404			{object}	BaseResponse
//	@Failure		400			{object}	BaseResponse
//	@Failure		500			{object}	BaseResponse
//	@Router			/probe [post]
func createProbe(c echo.Context) error {
	var r models.ServerBaseResponse

	namespace := c.QueryParam("namespace")
	repo := c.QueryParam("repo")

	if namespace == "" || repo == "" {
		r.Message = "Missing query parameters"
		r.Error = "Expect namespace and repo query params to be provided"

		return c.JSON(http.StatusBadRequest, r)
	}

	err := Beacon.Registry().TestRepo(namespace, repo)

	if err != nil {
		r.Message = fmt.Sprintf("Could not fetch repo %s in namespace %s", repo, namespace)

		if _, ok := err.(registry.GeneralServerError); ok {
			r.Error = err.Error()

			return c.JSON(http.StatusInternalServerError, r)
		}

		if _, ok := err.(registry.GeneralClientError); ok {
			r.Error = err.Error()

			return c.JSON(http.StatusBadRequest, r)
		}

		if _, ok := err.(registry.NotFoundError); ok {
			r.Error = err.Error()

			return c.JSON(http.StatusNotFound, r)
		}
	}

	err = Beacon.StartProbe(namespace, repo, time.Second*20)

	if _, ok := err.(BeaconErrorProbeAlreadyExists); ok {
		r.Error = err.Error()
		r.Message = fmt.Sprintf("Probe already exists for repo %s at namespace %s", repo, namespace)

		return c.JSON(http.StatusConflict, r)
	}

	r.Message = fmt.Sprintf("Probe successfully created for repo %s at namespace %s", repo, namespace)
	return c.JSON(http.StatusCreated, r)
}

// listProbes handles the GET /probes method for beacond
//
//	@Summary		Lists all probes
//	@Description	lists probes that are running for beacond
//	@Produce		json
//	@Success		200	{object}	ListProbesResponse
//	@Router			/probes [get]
func listProbes(c echo.Context) error {
	var r models.ServerListProbesResponse

	r.Probes = Beacon.ListProbes()

	return c.JSON(http.StatusOK, r)
}

// getBeaconDetails handles the GET /beacon method for beacond
//
//	@Summary		Get beacon details
//	@Description	describes the current status of beacond
//	@Produce		json
//	@Success		200	{object}	BeaconDescribeResponse
//	@Router			/beacon [get]
func getBeaconDetails(c echo.Context) error {
	var r models.ServerBeaconDescribeResponse

	r.Registry = Beacon.Registry().URL()
	r.Probes = Beacon.ListProbes()
	r.Runtime = string(Beacon.Runtime().Type())

	return c.JSON(http.StatusOK, r)
}
