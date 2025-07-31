package restserver

import (
	"casino/boundary/logging"
	"net/http"
)

type Server interface {
	RegisterPublicRoute(method, path string, handler http.HandlerFunc, logger logging.Logger)
	RegisterSwaggerRoutes()
	Start(address string) error
}
