package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/aalexanderkevin/getstarvio-backend/internal/app"
	"github.com/aalexanderkevin/getstarvio-backend/internal/http/middleware"
	"github.com/aalexanderkevin/getstarvio-backend/internal/http/routes"
)

type Server struct {
	engine *gin.Engine
	cfg    app.Container
}

func NewServer(c *app.Container) *Server {
	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.HTTPLogger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "X-Internal-Token", "X-Request-ID"},
		ExposeHeaders:    []string{"X-Request-ID"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	routes.Register(r, c)

	return &Server{engine: r, cfg: *c}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.cfg.Cfg.Service.Host, s.cfg.Cfg.Service.Port)
	return s.engine.Run(addr)
}

func (s *Server) Handler() http.Handler {
	return s.engine
}
