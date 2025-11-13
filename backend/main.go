package main

import (
	"log"
	"os"

	"github.com/ZiplEix/test-better-auth/backend/database"
	"github.com/ZiplEix/test-better-auth/backend/handlers"
	mdw "github.com/ZiplEix/test-better-auth/backend/middleware"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	godotenv.Load()

	database.InitDB()

	authURL := os.Getenv("BETTER_AUTH_URL")
	if err := mdw.InitJWKS(authURL); err != nil {
		log.Fatalf("Failed to initialize JWKS: %v", err)
	}

	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	g := e.Group("/api")
	g.Use(mdw.JWTMiddleware)

	g.GET("/todos", handlers.GetTodos)
	g.POST("/todos", handlers.CreateTodo)
	g.PATCH("/todos/:id/toggle", handlers.ToggleTodo)
	g.DELETE("/todos/:id", handlers.DeleteTodo)

	log.Println("Server running on http://localhost:8080")
	e.Logger.Fatal(e.Start(":8080"))
}
