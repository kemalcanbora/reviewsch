package handler

import (
	"context"
	"fmt"
	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"reviewsch/internal/service/entity"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Service interface defines the required business operations
type Service interface {
	ApplyCoupon(entity.Basket, string) (*entity.Basket, error)
	CreateCoupon(int, string, float64) error
	GetCoupons([]string) ([]entity.Coupon, error)
}

// RateLimitConfig holds the rate limiting configuration
type RateLimitConfig struct {
	Enabled      bool
	RedisAddr    string
	RedisPass    string
	RatePerSec   uint
	BurstSize    int
	KeyFunc      func(*gin.Context) string
	ErrorHandler func(*gin.Context, ratelimit.Info)
}

// Config holds the gateway configuration
type Config struct {
	Host           string
	Port           int
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
	AllowedOrigins []string
	AllowedMethods []string
	RateLimit      RateLimitConfig
}

// Gateway represents the API Gateway
type Gateway struct {
	server      *http.Server
	Engine      *gin.Engine
	config      Config
	services    map[string]interface{}
	middleware  []gin.HandlerFunc
	routes      []RouteDefinition
	mu          sync.RWMutex
	redisClient *redis.Client
}

// RouteDefinition defines structure for route registration
type RouteDefinition struct {
	Path       string
	Method     string
	Handler    gin.HandlerFunc
	Middleware []gin.HandlerFunc
}

// New creates a new API Gateway instance
func New(cfg Config) *Gateway {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())

	g := &Gateway{
		Engine:   engine,
		config:   cfg,
		services: make(map[string]interface{}),
	}

	if cfg.AllowedOrigins != nil {
		engine.Use(CORSMiddleware(cfg))
	}
	if cfg.RateLimit.Enabled {
		log.Println("Rate limit enabled")
		g.setupRateLimit()
	}
	return g
}

func (g *Gateway) createRateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := fmt.Sprintf("rate:%s", ip)
		ctx := context.Background()

		_, err := g.redisClient.SetNX(ctx, key, 1, 30*time.Second).Result()
		if err != nil {
			c.Next()
			return
		}

		count, err := g.redisClient.Incr(ctx, key).Result()
		if err != nil {
			c.Next()
			return
		}

		log.Printf("IP: %s, Request count: %d", ip, count)

		// Check if over limit
		if count > int64(g.config.RateLimit.RatePerSec) {
			g.config.RateLimit.ErrorHandler(c, ratelimit.Info{
				ResetTime: time.Now().Add(time.Second),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (g *Gateway) setupRateLimit() {
	// Log current KeyFunc status
	log.Printf("Initial KeyFunc status: %v", g.config.RateLimit.KeyFunc != nil)

	// Set default key function if not provided
	if g.config.RateLimit.KeyFunc == nil {
		log.Println("Setting default KeyFunc")
		g.config.RateLimit.KeyFunc = func(c *gin.Context) string {
			clientIP := c.ClientIP()
			log.Printf("Rate limit key generated for IP: %s", clientIP)
			return clientIP
		}
	}
	log.Printf("Initial ErrorHandler status: %v", g.config.RateLimit.ErrorHandler != nil)

	// Set default error handler if not provided
	if g.config.RateLimit.ErrorHandler == nil {
		log.Println("Setting default ErrorHandler")
		g.config.RateLimit.ErrorHandler = func(c *gin.Context, info ratelimit.Info) {
			log.Printf("Rate limit exceeded for IP: %s. Reset in: %v",
				c.ClientIP(),
				time.Until(info.ResetTime),
			)
			c.JSON(429, gin.H{
				"error": fmt.Sprintf("Rate limit exceeded. Try again in %v",
					time.Until(info.ResetTime)),
			})
		}
	} else {
		// Wrap the existing error handler to add logging
		originalHandler := g.config.RateLimit.ErrorHandler
		g.config.RateLimit.ErrorHandler = func(c *gin.Context, info ratelimit.Info) {
			log.Printf("Rate limit exceeded for IP: %s. Reset in: %v",
				c.ClientIP(),
				time.Until(info.ResetTime),
			)
			originalHandler(c, info)
		}
	}

	// Initialize Redis client
	g.redisClient = redis.NewClient(&redis.Options{
		Addr:     g.config.RateLimit.RedisAddr,
		Password: g.config.RateLimit.RedisPass,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := g.redisClient.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Redis connection failed: %v", err)
	}

	// Add rate limit middleware to the engine
	g.UseMiddleware(g.createRateLimitMiddleware())
	log.Println("Rate limit middleware configured")
}

// RegisterService adds a service to the gateway
func (g *Gateway) RegisterService(name string, service interface{}) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.services[name] = service
}

// UseMiddleware adds middleware to the gateway
func (g *Gateway) UseMiddleware(middleware ...gin.HandlerFunc) {
	g.middleware = append(g.middleware, middleware...)
	for _, m := range middleware {
		g.Engine.Use(m)
	}
}

// RegisterRoutes registers multiple routes at once
func (g *Gateway) RegisterRoutes(routes []RouteDefinition) {
	g.routes = append(g.routes, routes...)
	for _, route := range routes {
		handlers := append(route.Middleware, route.Handler)
		g.Engine.Handle(route.Method, route.Path, handlers...)
	}
}

// Group creates a new route group
func (g *Gateway) Group(path string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return g.Engine.Group(path, handlers...)
}

// Start begins serving the API Gateway
func (g *Gateway) Start() error {
	g.server = &http.Server{
		Addr:           fmt.Sprintf("%s:%d", g.config.Host, g.config.Port),
		Handler:        g.Engine,
		ReadTimeout:    g.config.ReadTimeout,
		WriteTimeout:   g.config.WriteTimeout,
		MaxHeaderBytes: g.config.MaxHeaderBytes,
	}

	log.Printf("Starting server on %s", g.server.Addr)
	return g.server.ListenAndServe()
}

// Stop gracefully shuts down the API Gateway
func (g *Gateway) Stop(ctx context.Context) error {
	log.Println("Shutting down server...")
	if g.redisClient != nil {
		if err := g.redisClient.Close(); err != nil {
			log.Printf("Error closing Redis client: %v", err)
		}
	}

	return g.server.Shutdown(ctx)
}

// GetService retrieves a registered service
func (g *Gateway) GetService(name string) interface{} {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.services[name]
}

// CORSMiddleware adds CORS headers to the response
func CORSMiddleware(config Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods",
			fmt.Sprintf("%v", config.AllowedMethods))
		c.Writer.Header().Set("Access-Control-Allow-Headers",
			"Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
