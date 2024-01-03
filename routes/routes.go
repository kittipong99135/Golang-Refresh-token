package routes

import (
	"god-dev/controllers"
	"god-dev/middleware"

	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
	api := app.Group("/api")

	//Group authentication routers
	auth := api.Group("/auth")

	//Authentication routers function
	auth.Post("/register", controllers.Regis)
	auth.Post("login", controllers.Login)

	//Group user routers
	user := api.Group("/user", middleware.RequestAuth(), middleware.RefreshAuth())

	//Users routers function
	user.Post("/logout", controllers.UserLogout)
	user.Get("/params/dashboard", controllers.UserParams)
	user.Get("/", controllers.UserList)
	user.Get("/:id", controllers.UserRead)
	user.Put("/:id", controllers.UserUpdate)
	user.Put("/active/:id", controllers.UserActive)
	user.Delete("/:id", controllers.UserRemove)
}
