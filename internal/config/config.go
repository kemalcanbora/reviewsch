package config

import (
	"fmt"
	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"reviewsch/internal/api/handler"
	"strconv"
	"strings"
	"time"
)

func getProjectRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			return wd, nil
		}

		parent := filepath.Dir(wd)
		if parent == wd {
			return "", fmt.Errorf("could not find project root (no go.mod file found)")
		}
		wd = parent
	}
}

func init() {
	projectRoot, err := getProjectRoot()
	if err != nil {
		fmt.Printf("Warning: Could not determine project root: %v\n", err)
		return
	}

	envPath := filepath.Join(projectRoot, ".env")
	if err := godotenv.Load(envPath); err != nil {
		fmt.Printf("Warning: Error loading .env file: %v\n", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string, sep string) []string {
	if value, exists := os.LookupEnv(key); exists {
		return strings.Split(value, sep)
	}
	return defaultValue
}

func LoadConfig() (*handler.Config, error) {
	// Required value check
	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		return nil, fmt.Errorf("REDIS_PASSWORD is required")
	}

	return &handler.Config{
		Port:           getEnvAsInt("SERVER_PORT", 8080),
		ReadTimeout:    getEnvAsDuration("SERVER_READ_TIMEOUT", 10*time.Second),
		WriteTimeout:   getEnvAsDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
		MaxHeaderBytes: getEnvAsInt("SERVER_MAX_HEADER_BYTES", 1<<20),

		// CORS configuration
		AllowedMethods: getEnvAsSlice("CORS_ALLOWED_METHODS",
			[]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, ","),
		AllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS",
			[]string{"*"}, ","),

		// Rate limit configuration
		RateLimit: handler.RateLimitConfig{
			Enabled: getEnvAsBool("RATE_LIMIT_ENABLED", true) &&
				getEnvAsBool("REDIS_ENABLED", true),
			RedisAddr:  getEnv("REDIS_ADDRESS", "localhost:6379"),
			RedisPass:  redisPassword,
			RatePerSec: uint(getEnvAsInt("RATE_LIMIT_PER_SEC", 5)),
			BurstSize:  getEnvAsInt("RATE_LIMIT_BURST_SIZE", 10),
			KeyFunc: func(c *gin.Context) string {
				return c.ClientIP()
			},
			ErrorHandler: func(c *gin.Context, info ratelimit.Info) {
				c.JSON(429, gin.H{
					"error": fmt.Sprintf("Too many requests. Try again in %v",
						time.Until(info.ResetTime)),
				})
			},
		},
	}, nil
}

func NewDefault() *handler.Config {
	cfg, err := LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load config from environment: %v, using defaults\n", err)
		return getDefaultConfig()
	}
	return cfg
}

func getDefaultConfig() *handler.Config {
	return &handler.Config{
		Port:           8080,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedOrigins: []string{"*"},
		RateLimit: handler.RateLimitConfig{
			Enabled:    true,
			RedisAddr:  "localhost:6379",
			RedisPass:  "",
			RatePerSec: 5,
			BurstSize:  10,
			KeyFunc: func(c *gin.Context) string {
				return c.ClientIP()
			},
			ErrorHandler: func(c *gin.Context, info ratelimit.Info) {
				c.JSON(429, gin.H{
					"error": fmt.Sprintf("Too many requests. Try again in %v",
						time.Until(info.ResetTime)),
				})
			},
		},
	}
}
