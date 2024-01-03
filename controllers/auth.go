package controllers

import (
	"context"
	"god-dev/database"
	"god-dev/models"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// Functions register - Authen routher.
func Regis(c *fiber.Ctx) error { // Routes -> http://127.0.0.1:3000/api/auth/register

	// Connect Database.
	db := database.DBConn

	// Recive body parser register routes.
	var regisBody models.User
	err := c.BodyParser(&regisBody)
	if err != nil { // Case : Input json data invalid.
		return c.Status(503).JSON(fiber.Map{
			"status":  "error",
			"message": "Error : Regid body invalid.",
			"error":   err.Error(),
		})
	}

	// Checked user exists.
	var userExists models.User
	result := db.Find(&userExists, "email = ?", strings.TrimSpace(regisBody.Email))
	if result.RowsAffected != 0 {
		return c.Status(503).JSON(fiber.Map{ // Case : Input email id exists.
			"status":  "error",
			"message": "Error : Email exists.",
		})
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(regisBody.Password), 10)
	if err != nil {
		c.Status(503).JSON(fiber.Map{ // Case : Hash password invalid.
			"status":  "error",
			"message": "Error : Invalid password hashing.",
			"error":   err.Error(),
		})
	}

	// Create register data.
	userRegisted := models.User{
		Email:    regisBody.Email,
		Password: "secretpass:" + string(hash),
		Name:     regisBody.Name,
		Phone:    regisBody.Phone,
		Age:      regisBody.Age,
		Rank:     regisBody.Rank,
		Status:   "nactive",
		Role:     "user",
	}

	// Insert register data into database.
	db.Create(&userRegisted)

	// Return Status200, Json data
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"massage": "User register success",
		"detail":  userRegisted,
	})
}

// Functions register - Authen login.
func Login(c *fiber.Ctx) error { // Routes -> http://127.0.0.1:3000/api/auth/login

	// Connect Database.
	db := database.DBConn

	// Recive body parser register routes.
	var loginBody models.RequestLogin
	err := c.BodyParser(&loginBody)
	if err != nil { // Case : Input json data invalid.
		return c.Status(503).JSON(fiber.Map{
			"status":  "error",
			"message": "Error : Logging body invalid.",
			"error":   err.Error(),
		})
	}

	// Checked email address.
	var user models.User
	result := db.Find(&user, "email = ?", strings.TrimSpace(loginBody.Email))
	if result.RowsAffected == 0 { // Case : can't find email in database.
		return c.Status(503).JSON(fiber.Map{
			"status":  "error",
			"message": "Error : Email invalid.",
		})
	}

	// Compare password.
	splitPass := strings.Split(user.Password, ":")[1:][0]
	err = bcrypt.CompareHashAndPassword([]byte(splitPass), []byte(loginBody.Password))

	if err != nil { // Case : Input password invalid.
		return c.Status(503).JSON(fiber.Map{
			"status":  "error",
			"message": "Error : Compare password error.",
			"error":   err.Error(),
		})
	}

	// Convert to string user.ID.
	udid := strconv.Itoa(int(user.ID))

	// Create access token then save "JWT_SECRET" in .env file
	acc_token, err := CreateToken(udid, "JWT_SECRET")
	if err != nil {
		c.Status(503).JSON(fiber.Map{ // Case : Create access token invalid.
			"status":  "error",
			"message": "Error : Create access token  error.",
			"error":   err.Error(),
		})
	}
	// Save access token to redis server -> key:access_toekn:udid, val:acc_token
	SetAccessToken("access_token:"+udid, acc_token)

	// Create refresh token then save "JWT_REFRESH" in .env file
	rfh_token, err := CreateToken(udid, "JWT_REFRESH")
	if err != nil {
		c.Status(503).JSON(fiber.Map{ // Case : Create refresh token invalid.
			"status":  "error",
			"message": "Error : Create refresh token error.",
			"error":   err.Error(),
		})
	}
	// Save refresh token to redis server -> key:refresh_token:udid, val:rfh_token
	SetRefreshToken("refresh_token:"+udid, rfh_token)

	// Return Status200, json data
	return c.Status(200).JSON(fiber.Map{
		"status":        "success",
		"message":       "Success : Logging in success.",
		"user":          udid,
		"token_access":  GetToken("access_token:" + udid),
		"token_refresh": GetToken("refresh_token:" + udid),
	})
}

func CreateToken(udid string, env string) (string, error) {
	cliams := jwt.MapClaims{"uid": udid}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cliams)
	return token.SignedString([]byte(os.Getenv("env")))
}

func SetAccessToken(key string, token string) {
	rd := database.RDConn
	ctx := context.Background()
	rd.Set(ctx, key, token, time.Hour*2)
}

func SetRefreshToken(key string, token string) {
	rd := database.RDConn
	ctx := context.Background()
	rd.Set(ctx, key, token, 0)
}

func GetToken(key string) string {
	rd := database.RDConn
	ctx := context.Background()
	val, _ := rd.Get(ctx, key).Result()
	return val
}
