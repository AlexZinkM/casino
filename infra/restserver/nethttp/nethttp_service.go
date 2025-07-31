package nethttp

import (
	"casino/boundary/logging"
	"casino/infra/middleware"
	"casino/infra/restserver"
	"net/http"
	"strings"
	"time"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type route struct {
	method  string
	path    string
	handler http.HandlerFunc
}

type NetHttpServer struct {
	routes []route
}

func NewNetHttpServer() restserver.Server {
	return &NetHttpServer{
		routes: make([]route, 0),
	}
}

func (s *NetHttpServer) RegisterPublicRoute(method, path string,
	handler http.HandlerFunc, logger logging.Logger) {

	wrappedHandler := middleware.LoggingMiddleware(handler, logger)
	s.routes = append(s.routes, route{
		method:  method,
		path:    path,
		handler: wrappedHandler,
	})
}

func (s *NetHttpServer) RegisterSwaggerRoutes() {
	swaggerHandler := httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)

	http.Handle("/swagger/", swaggerHandler)
}

func (s *NetHttpServer) Start(address string) error {
	var handler http.Handler = http.HandlerFunc(s.handleAll)
	
	timeoutHandler := http.TimeoutHandler(handler, 1*time.Second, "Service is not available")

	server := http.Server{
		Addr:    address,
		Handler: timeoutHandler,
	}

	return server.ListenAndServe()
}

func (s *NetHttpServer) handleAll(w http.ResponseWriter, r *http.Request) {

	if strings.HasPrefix(r.URL.Path, "/swagger/") {
		http.DefaultServeMux.ServeHTTP(w, r)
		return
	}

	for _, route := range s.routes {
		if route.method == r.Method && route.path == r.URL.Path {
			route.handler(w, r)
			return
		}
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"error": "no such path"}`))

}
