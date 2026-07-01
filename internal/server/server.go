package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gumla-hds/gumla-backend/config"
	"github.com/gumla-hds/gumla-backend/internal/auth"
	"github.com/gumla-hds/gumla-backend/internal/category"
	"github.com/gumla-hds/gumla-backend/internal/middleware"
	"github.com/gumla-hds/gumla-backend/internal/product"
	"github.com/gumla-hds/gumla-backend/pkg/firebase"
	jwtpkg "github.com/gumla-hds/gumla-backend/pkg/jwt"
	"github.com/gumla-hds/gumla-backend/pkg/razorpay"
)

type Server struct {
	Config        *config.Config
	DB            *pgxpool.Pool
	Firebase      *firebase.Client
	JWTManager    *jwtpkg.Manager
	Razorpay      *razorpay.Client
	Router        *gin.Engine
	httpServer    *http.Server
}

func New(cfg *config.Config, db *pgxpool.Pool, fb *firebase.Client, rp *razorpay.Client) *Server {
	jwtManager := jwtpkg.NewManager(
		cfg.JWTSecret,
		cfg.JWTConfig.AccessExpiry,
		cfg.JWTConfig.RefreshExpiry,
	)

	gin.SetMode(cfg.GinMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.RateLimit(cfg.RateLimit.RequestsPerMinute))

	s := &Server{
		Config:     cfg,
		DB:         db,
		Firebase:   fb,
		JWTManager: jwtManager,
		Razorpay:   rp,
		Router:     router,
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	authRepo := auth.NewRepository(s.DB)

	var firebaseVerifier auth.FirebaseIDTokenVerifier
	if s.Firebase != nil && s.Firebase.Auth != nil {
		firebaseVerifier = auth.NewFirebaseVerifier(s.Firebase.Auth)
	}

	authSvc := auth.NewService(authRepo, s.JWTManager, firebaseVerifier)

	authHandler := auth.NewHandler(authSvc, s.JWTManager)

	categoryRepo := category.NewRepository(s.DB)
	categorySvc := category.NewService(categoryRepo)
	categoryHandler := category.NewHandler(categorySvc)

	productRepo := product.NewRepository(s.DB)
	productSvc := product.NewService(productRepo)
	productHandler := product.NewHandler(productSvc)

	api := s.Router.Group("/api/v1")
	{
		api.GET("/health", s.handleHealth)

		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/dev-login", authHandler.DevLogin)
			authGroup.POST("/refresh", authHandler.Refresh)

			protected := authGroup.Group("")
			protected.Use(middleware.RequireAuth(s.JWTManager, authSvc.IsTokenBlacklisted))
			{
				protected.POST("/logout", authHandler.Logout)
				protected.POST("/register-device", authHandler.RegisterDevice)
			}
		}

		categoryGroup := api.Group("/categories")
		{
			categoryGroup.GET("", categoryHandler.List)
			categoryGroup.GET("/:id", categoryHandler.GetByID)

			adminCategories := categoryGroup.Group("")
			adminCategories.Use(middleware.RequireAuth(s.JWTManager, authSvc.IsTokenBlacklisted))
			adminCategories.Use(middleware.RequireRole("admin"))
			{
				adminCategories.POST("", categoryHandler.Create)
				adminCategories.PUT("/:id", categoryHandler.Update)
				adminCategories.DELETE("/:id", categoryHandler.Delete)
			}
		}

		productGroup := api.Group("/products")
		{
			productGroup.GET("", productHandler.List)
			productGroup.GET("/:id", productHandler.GetByID)

			adminProducts := productGroup.Group("")
			adminProducts.Use(middleware.RequireAuth(s.JWTManager, authSvc.IsTokenBlacklisted))
			adminProducts.Use(middleware.RequireRole("admin"))
			{
				adminProducts.POST("", productHandler.Create)
				adminProducts.PUT("/:id", productHandler.Update)
				adminProducts.DELETE("/:id", productHandler.Delete)
			}
		}
	}
}

func (s *Server) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"status":  "ok",
			"version": "1.0.0",
			"time":    time.Now().UTC().Format(time.RFC3339),
		},
	})
}

func (s *Server) Start(port string) error {
	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      s.Router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Server starting on port %s", port)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced shutdown: %w", err)
	}

	log.Println("Server exited gracefully")
	return nil
}
