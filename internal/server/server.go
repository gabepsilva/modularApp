package server

import (
	"context"
	"net/http"
	"time"

	"modularApp/pkg/auth"
	"modularApp/pkg/frontend"
	"modularApp/pkg/items"
	"modularApp/pkg/web"
)

// Server represents the HTTP server
type Server struct {
	webModule      *web.Module
	authModule     *auth.Module
	itemsModule    *items.Module
	frontendModule *frontend.Module
	httpSrv        *http.Server
}

// New creates a new server instance
func New() *Server {
	// Initialize modules
	webModule := web.NewModule()
	authModule := auth.NewModule()
	itemsModule := items.NewModule(webModule, authModule)
	frontendModule := frontend.NewModule(webModule)

	// Register routes for base modules
	webModule.RegisterRoutes()

	// Register routes for other modules
	itemsModule.RegisterRoutes()
	frontendModule.RegisterRoutes()

	return &Server{
		webModule:      webModule,
		authModule:     authModule,
		itemsModule:    itemsModule,
		frontendModule: frontendModule,
	}
}

// Start starts the HTTP server
func (s *Server) Start(addr string) error {
	s.httpSrv = &http.Server{
		Addr:    addr,
		Handler: s.webModule.Router(),
	}

	return s.httpSrv.ListenAndServe()
}

// Stop gracefully stops the HTTP server
func (s *Server) Stop() error {
	if s.httpSrv == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.httpSrv.Shutdown(ctx)
}
