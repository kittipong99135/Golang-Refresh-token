package controllers

import (
	"context"
	"fmt"
	"god-dev/database"
	"god-dev/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Functions Userparams after login - User routher.
func UserParams(c *fiber.Ctx) error { // Routes -> http://127.0.0.1:3000/api/user/params/dashboard
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	fmt.Print(claims["uid"])
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Success : Read user success",
		"id":      claims["uid"],
	})
}

// Functions List all users - User routher .
func UserList(c *fiber.Ctx) error { // Routes -> http://127.0.0.1:3000/api/user/
	db := database.DBConn

	var listUser []models.User
	resultUser := db.Find(&listUser)
	if resultUser.RowsAffected == 0 {
		return c.Status(503).JSON(fiber.Map{
			"status":  "warning",
			"message": "Warning : Can't find user",
			"result":  "No result.",
		})
	}

	return c.Status(500).JSON(fiber.Map{
		"status":  "success",
		"message": "Success : Find all user success",
		"result":  listUser,
	})
}

func UserRead(c *fiber.Ctx) error {
	db := database.DBConn

	id := c.Params("id")

	var readUser models.User

	resultUser := db.Find(&readUser, "id = ?", id)
	if resultUser.RowsAffected == 0 {
		return c.Status(503).JSON(fiber.Map{
			"status":  "warning",
			"message": "Warning : Can't find user",
			"result":  "No result.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Success : Find user success",
		"result":  readUser,
	})
}

func UserUpdate(c *fiber.Ctx) error {
	db := database.DBConn
	id := c.Params("id")

	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(503).JSON(fiber.Map{
			"status":  "error",
			"message": "Error : Update body parser error.",
			"error":   err.Error(),
		})
	}

	updateUser := models.User{
		Email: user.Email,
		Name:  user.Name,
		Phone: user.Phone,
		Age:   user.Age,
		Rank:  user.Rank,
	}

	db.Where("id = ?", id).Updates(&updateUser)
	return c.Status(503).JSON(fiber.Map{
		"status":  "error",
		"message": "Error : Update body parser error.",
		"result":  updateUser,
	})
}

func UserRemove(c *fiber.Ctx) error {
	db := database.DBConn
	id := c.Params("id")

	var user models.User

	db.Delete(&user, id)
	return c.Status(503).JSON(fiber.Map{
		"status":  "success",
		"message": "Success : Remove user success",
	})
}

func UserActive(c *fiber.Ctx) error {
	db := database.DBConn

	id := c.Params("id")

	var user models.User

	result := db.First(&user, "id = ?", id)
	if result.RowsAffected == 0 {
		return c.Status(503).JSON(fiber.Map{
			"status":  "warning",
			"message": "Warning : Can't find user",
			"result":  "No result.",
		})
	}

	activeUser := models.User{
		Email:    user.Email,
		Password: user.Password,
		Name:     user.Name,
		Phone:    user.Phone,
		Age:      user.Age,
		Rank:     user.Rank,
		Status:   "active",
		Role:     user.Role,
	}

	db.Where("id = ?", id).Updates(&activeUser)

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Success : Active user success.",
		"result":  activeUser,
	})
}

func UserLogout(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	uid := fmt.Sprintf("%v", claims["uid"])
	result, err := DeleteFromRedis("access_token:" + uid)
	if err != nil {
		return c.Status(503).JSON(fiber.Map{
			"status":  "error",
			"message": "Error : Delete access token from redis server.",
			"error":   err.Error(),
		})
	}
	fmt.Println("Delete access toeken : " + result + " | Success")

	result, err = DeleteFromRedis("refresh_token:" + uid)
	if err != nil {
		return c.Status(503).JSON(fiber.Map{
			"status":  "error",
			"message": "Error : Delete refresh  token from redis server.",
			"error":   err.Error(),
		})
	}
	fmt.Println("Delete access toeken : " + result + " | Success")

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User Logout success.",
	})

}

func DeleteFromRedis(key string) (string, error) {
	rd := database.RDConn
	ctx := context.Background()
	val, err := rd.GetDel(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}
