package main

import (
	"github.com/Leynaic/katten-go/handlers"
	jwtware "github.com/gofiber/jwt/v3"
	"log"
	"os"

	"github.com/Leynaic/katten-go/database"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	if err := godotenv.Load("kat.env"); err != nil {
		log.Println("The .env file is not found, continue without it...")
		log.Println(err.Error())
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("CORS_ORIGINS"),
	}))

	database.New(
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_NAME"),
		os.Getenv("POSTGRES_PORT"),
	)

	// authGroup := app.Group("")
	app.Post("/login", handlers.Login)
	app.Post("/register", handlers.Register)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRETKEY")),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "You must provide a valid token.",
				"error":   err.Error(),
			})
		},
	}))

	app.Get("/cats", handlers.GetCats)
	app.Post("/cats/like/:id", handlers.LikeCat)
	app.Delete("/cats/like/:id", handlers.CancelLikeCat)
	app.Post("/cats/dislike/:id", handlers.DislikeCat)

	app.Patch("/update/password", handlers.UpdatePassword)
	app.Patch("/update/description", handlers.UpdateDescription)

	err := app.Listen(":8080")
	if err != nil {
		log.Println(err)
	}
}
