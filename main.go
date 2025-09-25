package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	"tivix-performance-tracker-backend/database"
	"tivix-performance-tracker-backend/middleware"
	"tivix-performance-tracker-backend/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	database.Connect()

	database.Migrate()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			message := "Erro interno do servidor"

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				if code >= 400 && code < 500 {
					message = e.Message
				}
			}

			log.Printf("Error %d: %v", code, err)

			return ctx.Status(code).JSON(fiber.Map{
				"error":   true,
				"message": message,
			})
		},
		DisableStartupMessage: os.Getenv("ENVIRONMENT") == "production",
		ServerHeader:          "TivixAPI",
		AppName:               "Tivix Performance Tracker API",
	})

	var allowedOrigins []string

	if os.Getenv("ENVIRONMENT") == "development" {
		allowedOrigins = []string{
			"http://localhost:3000",
			"http://localhost:5173",
			"http://127.0.0.1:5173",
		}
	} else if os.Getenv("ENVIRONMENT") == "production" {
		allowedOrigins = []string{
			"https://performancetracker.tivix.com.br",
			"https://performance.valiantgroup.com.br",
		}
	} else {
		corsOrigin := os.Getenv("CORS_ORIGIN")
		if corsOrigin != "" {
			allowedOrigins = []string{corsOrigin}
		} else {
			allowedOrigins = []string{"http://localhost:5173"}
		}
	}

	var finalOrigins []string
	seen := make(map[string]bool)
	for _, origin := range allowedOrigins {
		if origin != "" && !seen[origin] {
			finalOrigins = append(finalOrigins, origin)
			seen[origin] = true
		}
	}

	log.Printf("🔧 CORS Origins permitidas: %v", finalOrigins)

	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(finalOrigins, ","),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With,X-CSRF-Token",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length,Content-Range",
		MaxAge:           86400,
	}))

	app.Options("/*", func(c *fiber.Ctx) error {
		origin := c.Get("Origin")
		for _, allowedOrigin := range finalOrigins {
			if origin == allowedOrigin {
				c.Set("Access-Control-Allow-Origin", origin)
				c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS,PATCH")
				c.Set("Access-Control-Allow-Headers", "Origin,Content-Type,Accept,Authorization,X-Requested-With,X-CSRF-Token")
				c.Set("Access-Control-Allow-Credentials", "true")
				c.Set("Access-Control-Max-Age", "86400")
				break
			}
		}
		return c.SendStatus(fiber.StatusOK)
	})

	app.Use(logger.New())

	app.Use(middleware.InputSizeLimit(10 * 1024 * 1024))

	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   true,
				"message": "Muitas requisições. Tente novamente em alguns instantes.",
			})
		},
	}))

	app.Use("/api/v1/auth/login", limiter.New(limiter.Config{
		Max:        5,
		Expiration: 15 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   true,
				"message": "Muitas tentativas de login. Tente novamente em 15 minutos.",
			})
		},
	}))

	routes.SetupRoutes(app)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Tivix Performance Tracker API is running",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
