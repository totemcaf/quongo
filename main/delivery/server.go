package delivery

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ant0ine/go-json-rest/rest"
)

// Handler ...
type Handler interface {
	Routes() []*rest.Route
}

// Server ...
type Server struct {
	port           int
	handlers       []Handler
	systemHandlers []Handler
}

// NewServer ...
func NewServer(port int) *Server {
	return &Server{port: port, handlers: []Handler{}, systemHandlers: []Handler{}}
}

// Add ...
func (s *Server) Add(handler Handler) {
	s.handlers = append(s.handlers, handler)
}

// Add ...
func (s *Server) AddSystem(handler Handler) {
	s.systemHandlers = append(s.systemHandlers, handler)
}

// Start ...
func (s *Server) Start() error {
	stack := []rest.Middleware{
		&rest.AccessLogApacheMiddleware{},
		//      Format: rest.CombinedLogFormat,
		&rest.TimerMiddleware{},
		&rest.RecorderMiddleware{},
		//		&rest.PoweredByMiddleware{},
		&rest.RecoverMiddleware{
			EnableResponseStackTrace: true,
		},
		&rest.JsonIndentMiddleware{},
		&rest.ContentTypeCheckerMiddleware{},
		&rest.GzipMiddleware{},
	}

	api := rest.NewApi()

	statusMw := &rest.StatusMiddleware{}
	api.Use(statusMw)

	api.Use(stack...)

	router := s.makeRouter(s.handlers)
	api.SetApp(router)

	system := rest.NewApi()

	systemRouter := s.makeRouter(s.systemHandlers)
	system.SetApp(systemRouter)

	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))
	http.Handle("/", system.MakeHandler())

	return http.ListenAndServe(":"+strconv.Itoa(s.port), nil)
}

func (s *Server) makeRouter(handlers []Handler) rest.App {
	routes := []*rest.Route{}

	for _, handler := range handlers {
		routes = append(routes, handler.Routes()...)
	}

	router, err := rest.MakeRouter(routes...)

	if err != nil {
		log.Fatal(err)
	}

	return router
}
