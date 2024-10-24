package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reviewsch/internal/api/handler"
	"reviewsch/internal/api/middleware/auth"
	"reviewsch/internal/api/router"
	"reviewsch/internal/config"
	"reviewsch/internal/repository/memdb"
	"reviewsch/internal/service"
	"reviewsch/swagger"
	"syscall"
	"time"
)

var (
	repo = memdb.New()
)

func Run() error {
	swagger.SetupSwagger()

	conf := config.NewDefault()
	gateway := handler.New(*conf)

	// Swagger documentation
	gateway.Engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Register services
	couponService := service.New(repo)
	gateway.RegisterService("coupon", couponService)

	// Register middleware
	gateway.UseMiddleware(handler.CORSMiddleware(*conf))

	setupRoutes(gateway, couponService)

	return startServer(gateway)
}

func setupRoutes(gateway *handler.Gateway, couponService *service.Service) {
	apiGroup := gateway.Engine.Group("/api")
	v1 := apiGroup.Group("/v1")

	// NewCouponHandler
	couponHandler := router.NewCouponHandler(couponService)

	// Coupons group
	coupons := v1.Group("/coupons")

	// Applied JWT middleware to all coupon routes
	coupons.Use(auth.AdminAuth())
	{
		coupons.POST("/apply", couponHandler.Apply)
		coupons.POST("/create", couponHandler.Create)
		coupons.GET("/", couponHandler.Get)
	}

	// Health check
	v1.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": true})
	})
}

func startServer(gateway *handler.Gateway) error {
	// Start server in goroutine
	go func() {
		if err := gateway.Start(); !errors.Is(err, http.ErrServerClosed) && err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := gateway.Stop(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %v", err)
	}

	return nil
}
