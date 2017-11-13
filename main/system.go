package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
	"runtime"
)

const APP_NAME = "Kafka-Http"

type BuildInfo struct {
	Name      string       `json:"name"`
	Version   string       `json:"version"`
	GoVersion string       `json:"goVersion"`
	Stats     *rest.Status `json:"httpStats"`
	RepoStats *interface{} `json:"repoStats""`
}

type SysStatusController struct {
	repository Repository
	statusMw   *rest.StatusMiddleware
	Routes     []*rest.Route
}

func NewSystemController(repository Repository, statusMw *rest.StatusMiddleware) *SysStatusController {
	ctrl := SysStatusController{repository: repository, statusMw: statusMw}

	routes := []*rest.Route{
		rest.Get("/system/status", ctrl.status),
		rest.Get("/system/stats", ctrl.stats),
		rest.Get("/system/health", ctrl.health),
	}

	ctrl.Routes = routes

	return &ctrl
}

func (c *SysStatusController) status(w rest.ResponseWriter, r *rest.Request) {
	repoStats := c.repository.Stat()
	stats := c.statusMw.GetStatus()
	v := BuildInfo{APP_NAME, VERSION, runtime.Version(), stats, &repoStats}
	w.WriteJson(v)
}

func (c *SysStatusController) stats(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(c.statusMw.GetStatus())
}

func (c *SysStatusController) health(w rest.ResponseWriter, r *rest.Request) {
	if err := c.repository.Health(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.WriteJson(err.Error)
	} else {
		w.WriteJson("Ok")
	}
}
