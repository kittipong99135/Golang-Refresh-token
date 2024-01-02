package routes // model : go-be

import (
	"go-be/controllers"
	"go-be/middleware"

	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
	auth := app.Group("/auth")
	auth.Post("/regis", controllers.Regis)
	auth.Post("/login", controllers.Login)
	auth.Get("/readed",
		middleware.ReqAuth(),
		middleware.RefAuth(),
		controllers.LoginReaded,
	)
}
