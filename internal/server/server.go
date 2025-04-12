package server

import (
	"context"
	"net/http"
	"time"

	"modularApp/pkg/items"
	"modularApp/pkg/web"
)

// Server represents the HTTP server
type Server struct {
	webModule   *web.Module
	itemsModule *items.Module
	httpSrv     *http.Server
}

// New creates a new server instance
func New() *Server {
	// Initialize modules
	webModule := web.NewModule()
	itemsModule := items.NewModule(webModule)
	
	// Register routes for base modules
	webModule.RegisterRoutes()
	
	// Register routes for other modules
	itemsModule.RegisterRoutes()
	
	return &Server{
		webModule:   webModule,
		itemsModule: itemsModule,
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