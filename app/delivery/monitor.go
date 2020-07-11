package delivery

import "github.com/ant0ine/go-json-rest/rest"

// StatusInteractor ...
type StatusInteractor interface {
	Status() error
}

// MonitorView ...
type MonitorView struct {
	interactor StatusInteractor
	routes     []*rest.Route
}

// NewMonitorView ...
func NewMonitorView(interactor StatusInteractor) *MonitorView {

	ctrl := MonitorView{
		interactor: interactor,
	}

	ctrl.routes = []*rest.Route{
		rest.Get("/ping", ctrl.Ping),
	}

	return &ctrl
}

// Routes ...
func (v *MonitorView) Routes() []*rest.Route {
	return v.routes
}

// Ping ...
func (v *MonitorView) Ping(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson("ok")
}
