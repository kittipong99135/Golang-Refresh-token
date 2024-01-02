package controllers //model : go-be

import (
	"context"
	"go-be/database"
	m "go-be/models"
	"os"
	"time"

	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func Regis(c *fiber.Ctx) error {
	db := database.DBConn
	var regisBody m.User

	err := c.BodyParser(&regisBody)
	if err != nil {
		return c.Status(503).JSON(fiber.Map{
			"message": "Error : Register body has invalid.",
			"status":  "error",
			"error":   err.Error(),
		})
	}

	var userExists m.User
	result := db.Find(&userExists, "email = ?", strings.TrimSpace(regisBody.Email))
	if result.RowsAffected != 0 {
		return c.Status(503).JSON(fiber.Map{
			"massage": "Error : Email exists.",
			"status":  "error",
		})
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(regisBody.Password), 10)
	if err != nil {
		c.Status(503).JSON(fiber.Map{
			"message": "Error : Password has invalid.",
			"status":  "error",
			"error":   err.Error(),
		})
	}

	userRegisted := m.User{
		Email:    regisBody.Email,
		Password: "secretpass:" + string(passHash),
		Status:   "active",
		Role:     "admin",
	}

	db.Create(&userRegisted)
	return c.Status(200).JSON(fiber.Map{
		"massage": "Seccess : Register success.",
		"status":  "seccess",
		"newuser": userRegisted,
	})
}

func Login(c *fiber.Ctx) error {
	db := database.DBConn

	type loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var loginBody loginRequest
	err := c.BodyParser(&loginBody)
	if err != nil {
		return c.Status(503).JSON(fiber.Map{
			"message": "Error : Logging body has some in invalid.",
			"status":  "error",
			"error":   err.Error(),
		})
	}

	var userSelect m.User
	selectResult := db.Find(&userSelect, "email = ?", strings.TrimSpace(loginBody.Email))
	if selectResult.RowsAffected == 0 {
		return c.Status(503).JSON(fiber.Map{
			"message": "Error : Email has invalid.",
			"status":  "error",
		})
	}

	dwarfsPass := strings.Split(userSelect.Password, ":")[1:][0]
	err = bcrypt.CompareHashAndPassword([]byte(dwarfsPass), []byte(loginBody.Password))
	if err != nil {
		return c.Status(503).JSON(fiber.Map{
			"message": "Error : Compare password error.",
			"status":  "error",
			"error":   err.Error(),
		})
	}

	udid := strconv.Itoa(int(userSelect.ID))

	claims := jwt.MapClaims{
		"id":     userSelect.ID,
		"email":  userSelect.Email,
		"status": userSelect.Status,
		"role":   userSelect.Role,
		"exp":    time.Now().Add(time.Minute * 1).Unix(),
	}
	access_token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	act, err := access_token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	SetActTokenRedis("access_token:"+udid, act)
	SetActTokenRedis("id_user:"+udid, udid)
	SetActTokenRedis("email_user:"+udid, userSelect.Email)
	SetActTokenRedis("status_user:"+udid, userSelect.Status)
	SetActTokenRedis("role_user:"+udid, userSelect.Role)

	claims = jwt.MapClaims{
		"id":  userSelect.ID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	}
	refresh_token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	rft, err := refresh_token.SignedString([]byte(os.Getenv("JWT_REFRESH")))
	SetRftTokenRedis("refresh_token:"+udid, rft)

	return c.Status(200).JSON(fiber.Map{
		"message":       "Success : Login success.",
		"status":        "success",
		"id":            udid,
		"access_token":  GetTokenRedis("access_token:" + udid),
		"refresh_token": GetTokenRedis("refresh_token:" + udid),
	})
}

func SetActTokenRedis(key string, token string) {
	rd := database.RDConn
	ctx := context.Background()
	rd.Set(ctx, key, token, time.Minute*1).Err()
}

func SetRftTokenRedis(key string, token string) {
	rd := database.RDConn
	ctx := context.Background()
	rd.Set(ctx, key, token, 0).Err()
}

func GetTokenRedis(key string) string {
	rd := database.RDConn
	ctx := context.Background()
	val, _ := rd.Get(ctx, key).Result()
	return val
}
