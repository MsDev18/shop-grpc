package server

import (
	"fmt"
	"log"
	"net/http"
	"shop/internal/api/handler/address"
	"shop/internal/api/handler/auth"
	"shop/internal/api/handler/category"
	"shop/internal/api/handler/health"
	"shop/internal/api/handler/product"
	"shop/internal/api/handler/province"
	"shop/internal/api/handler/user"
	authmiddleware "shop/internal/api/middleware/auth"
	"shop/internal/api/router"
	"shop/internal/pkg/richerror"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config     Config
	httpServer *http.Server
}

type Config struct {
	Host         string        `koanf:"host"`
	Port         uint          `koanf:"port"`
	WriteTimeout time.Duration `koanf:"write_timeout"`
	ReadTimeout  time.Duration `koanf:"read_timeout"`
	IdleTimeout  time.Duration `koanf:"idle_timeout"`
	Env          string        `koanf:"env"`
}

func New(config Config, healthHandler health.Handler, authHandler auth.Handler, userHandler user.Handler, cateogryHandler category.Handler, provinceHandler province.Handler, addressHandler address.Handler,productHandler product.Handler, authMiddleware authmiddleware.Middleware) Server {
	// validation env
	env := Env(config.Env)
	if !env.IsValid() {
		log.Fatalf("Invalid environment: %s", config.Env)
	}
	// validate timeout values
	if config.ReadTimeout == 0 || config.WriteTimeout == 0 || config.IdleTimeout == 0 {
		log.Fatalf("Invalid timeout values: read_timeout=%d, write_timeout=%d, idle_timeout=%d", config.ReadTimeout, config.WriteTimeout, config.IdleTimeout)
	}

	// set from gin mode based on env
	if env == EnvDevelopment {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	// create gin engine
	engine := gin.Default()
	// register routes
	appRouter := router.New(engine, healthHandler, authHandler, userHandler, cateogryHandler, provinceHandler, addressHandler, productHandler, authMiddleware)
	appRouter.Register()

	// manually create http server to set timeouts
	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler:      engine,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	}

	return Server{
		config:     config,
		httpServer: httpServer,
	}
}

func (s Server) Run() error {
	const op = "server.Run"
	log.Printf("start running server on %s:%d", s.config.Host, s.config.Port)
	err := s.httpServer.ListenAndServe()
	if err != nil {
		return richerror.New().
			SetOp(op).
			SetMsg("failed to run server").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}
	return nil
}
