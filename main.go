package main

import (
	"god-dev/database"
	"god-dev/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())

	database.DB_Init()
	database.RD_Init()

	routes.Routes(app)

	app.Listen("127.0.0.1:3000")
}
