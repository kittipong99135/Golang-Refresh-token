package middleware

import (
	"fmt"
	"god-dev/controllers"
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func RequestAuth() func(*fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
		ErrorHandler: errNext,
	})
}

func RefreshAuth() func(*fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey:     jwtware.SigningKey{Key: []byte(os.Getenv("JWT_REFRESH"))},
		SuccessHandler: resendToken,
	})
}

func errNext(c *fiber.Ctx, err error) error {
	c.Next()
	return nil
}

func resendToken(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	uid := fmt.Sprintf("%v", claims["uid"])
	getAccessToken := controllers.GetToken("access_token:" + uid)

	if getAccessToken != "" {
		fmt.Println("have token: \n" + getAccessToken + "\n")
		c.Next()
		return nil
	}

	acc_token, err := controllers.CreateToken(uid, "JWT_SECRET")
	if err != nil {
		return c.Status(503).JSON(fiber.Map{
			"status":  "error",
			"message": "Error : Refresh password invalid.",
			"error":   err.Error(),
		})
	}
	fmt.Println("have token: \n" + acc_token + "\n")
	controllers.SetAccessToken("access_token:"+uid, acc_token)

	c.Next()
	return c.Status(200).JSON(fiber.Map{
		"status":  "warning",
		"message": "Warning : Refresh token.",
		"token":   controllers.GetToken("access_token:" + uid),
	})

}
