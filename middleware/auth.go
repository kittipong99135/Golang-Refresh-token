package middleware

import (
	"fmt"
	"os"

	"go-be/controllers"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var udid string

func ReqAuth() func(*fiber.Ctx) error {
	fmt.Println("init id" + udid)
	return jwtware.New(jwtware.Config{
		SigningKey:     jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
		SuccessHandler: validatorRedis,
		ErrorHandler:   errNext,
	})
}

func RefAuth() func(*fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey:     jwtware.SigningKey{Key: []byte(os.Getenv("JWT_REFRESH"))},
		SuccessHandler: tokenRefresh,
	})
}

func errNext(c *fiber.Ctx, err error) error {
	c.Next()
	return nil
}

func tokenRefresh(c *fiber.Ctx) error {
	getToken := controllers.GetTokenRedis("access_token:" + udid)
	if getToken != "" {
		fmt.Println("have token:" + getToken)
		c.Next()
		return nil
	} else {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		id := fmt.Sprintf("%v", claims["id"])

		claims = jwt.MapClaims{
			"id": id,
		}
		access_token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		act, err := access_token.SignedString([]byte(os.Getenv("JWT_SECRET")))

		controllers.SetActTokenRedis("access_token:"+udid, act)

		if err != nil {
			return c.Status(503).JSON(fiber.Map{
				"message": "Error : Refresh password error.",
				"status":  "error",
				"error":   err.Error(),
			})
		}
		fmt.Println("refresh token:" + act)
		c.Next()
		return c.Status(503).JSON(fiber.Map{
			"message": "Warning : Token is refresh..",
			"status":  "warning",
			"token":   act,
		})

	}
}

func validatorRedis(c *fiber.Ctx) error {
	type redisId struct {
		Id string `json:"id"`
	}
	var id redisId
	c.BodyParser(&id)
	udid = id.Id
	fmt.Println(udid)
	c.Next()
	return nil

}
