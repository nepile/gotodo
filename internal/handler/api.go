package handler

import (
	"time"

	"github.com/nepile/gotodo/internal/config"
	"github.com/nepile/gotodo/internal/middleware"
	"github.com/nepile/gotodo/internal/model"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input tidak valid"})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	user := model.User{
		Username: data["username"],
		Password: string(password),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Username mungkin sudah terpakai"})
	}

	return c.JSON(fiber.Map{"message": "Register berhasil", "user": user})
}

func Login(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input tidak valid"})
	}

	var user model.User
	config.DB.Where("username = ?", data["username"]).First(&user)

	if user.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Password salah"})
	}

	// Generate JWT
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token valid 24 jam
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(middleware.JwtSecret)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membuat token"})
	}

	return c.JSON(fiber.Map{"message": "Login berhasil", "token": t})
}

// --- TODO HANDLERS ---

func GetTodos(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint) // Ambil ID dari token JWT
	var todos []model.Todo
	config.DB.Where("user_id = ?", userID).Find(&todos)
	return c.JSON(todos)
}

func CreateTodo(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	todo := new(model.Todo)

	if err := c.BodyParser(todo); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input tidak valid"})
	}

	todo.UserID = userID
	config.DB.Create(&todo)
	return c.JSON(todo)
}

func UpdateTodo(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	id := c.Params("id")

	var todo model.Todo
	// Pastikan hanya bisa update todo miliknya sendiri
	if err := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&todo).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Todo tidak ditemukan"})
	}

	type UpdateInput struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		IsCompleted bool   `json:"is_completed"`
	}
	var input UpdateInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input tidak valid"})
	}

	todo.Title = input.Title
	todo.Description = input.Description
	todo.IsCompleted = input.IsCompleted

	config.DB.Save(&todo)
	return c.JSON(todo)
}

func DeleteTodo(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	id := c.Params("id")

	var todo model.Todo
	if err := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&todo).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Todo tidak ditemukan"})
	}

	config.DB.Delete(&todo)
	return c.JSON(fiber.Map{"message": "Todo berhasil dihapus"})
}
