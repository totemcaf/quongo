package delivery

import (
	"runtime"

	"github.com/ant0ine/go-json-rest/rest"
)

// BuildInfo ...
type BuildInfo struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	GoVersion string `json:"goVersion"`
}

// SystemView ...
type SystemView struct {
	routes  []*rest.Route
	version string
}

// NewSystemView ...
func NewSystemView(version string) *SystemView {

	ctrl := SystemView{version: version}

	routes := []*rest.Route{
		rest.Get("/info", ctrl.sysStatus),
	}

	ctrl.routes = routes

	return &ctrl
}

// Routes ...
func (v *SystemView) Routes() []*rest.Route {
	return v.routes
}

func (v *SystemView) sysStatus(w rest.ResponseWriter, r *rest.Request) {
	info := BuildInfo{"quongo", v.version, runtime.Version()}
	w.WriteJson(info)
}
