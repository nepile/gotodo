package main

import (
	"github.com/nepile/gotodo/internal/handler"
	"github.com/nepile/gotodo/internal/middleware"

	"github.com/nepile/gotodo/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Koneksi dan Migrasi Database
	config.ConnectDB()

	app := fiber.New()

	// Middleware untuk melihat log request di terminal
	app.Use(logger.New())

	api := app.Group("/api")

	// Routes Auth (Tidak diproteksi)
	auth := api.Group("/auth")
	auth.Post("/register", handler.Register)
	auth.Post("/login", handler.Login)

	// Routes Todo (Diproteksi oleh JWT Middleware)
	todos := api.Group("/todos", middleware.Protected())
	todos.Get("/", handler.GetTodos)
	todos.Post("/", handler.CreateTodo)
	todos.Put("/:id", handler.UpdateTodo)
	todos.Delete("/:id", handler.DeleteTodo)

	// Jalankan server di port 8080
	app.Listen(":8080")
}
