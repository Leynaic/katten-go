package handlers

import (
	"github.com/Leynaic/katten-go/database"
	"github.com/Leynaic/katten-go/helpers"
	"github.com/Leynaic/katten-go/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type UserAuth struct {
	Username string `json:"username" xml:"username" form:"username"`
	Password string `json:"password" xml:"password" form:"password"`
}

func GetProfile(c *fiber.Ctx) error {
	currentCat := helpers.GetCurrentCat(c)
	return c.Status(fiber.StatusOK).JSON(currentCat)
}

func UpdatePassword(c *fiber.Ctx) error {
	newPassword := c.Body()
	currentCat := helpers.GetCurrentCat(c)
	if len(newPassword) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Your password must be at least 8 characters",
		})
	}

	if bytes, err := bcrypt.GenerateFromPassword(newPassword, 14); err == nil {
		currentCat.Password = string(bytes)
		database.GetInstance().Save(&currentCat)
		return c.JSON(currentCat)
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
}

func UpdateDescription(c *fiber.Ctx) error {
	description := c.Body()
	currentCat := helpers.GetCurrentCat(c)
	if len(description) < 20 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Your description must be at least 20 characters",
		})
	}

	currentCat.Description = string(description)
	database.GetInstance().Save(&currentCat)
	return c.JSON(currentCat)
}

func Login(c *fiber.Ctx) error {
	auth := new(UserAuth)

	if err := c.BodyParser(auth); err != nil {
		return invalid(c)
	} else {
		user := new(models.Cat)

		db := database.GetInstance()
		db = db.Where(
			"username = ?",
			auth.Username,
		).First(user)

		if db.Error != nil {
			return invalid(c)
		} else if checkPasswordHash(auth.Password, user.Password) {
			if token, err := createSignedToken(user.ID, user.Username); err != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			} else {
				return c.JSON(fiber.Map{
					"token": token,
				})
			}
		} else {
			return invalid(c)
		}
	}
}

func Register(c *fiber.Ctx) error {
	auth := new(UserAuth)

	if err := c.BodyParser(auth); err != nil {
		return invalid(c)
	} else if bytes, err := bcrypt.GenerateFromPassword([]byte(auth.Password), 14); err == nil {
		user := models.Cat{
			Username: auth.Username,
			Password: string(bytes),
		}

		var count int64

		db := database.GetInstance()
		db.Model(&models.Cat{}).Where("username = ?", auth.Username).Count(&count)

		if count == 0 {
			db = db.Create(&user)
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"message": "Username already exists",
			})
		}

		if db.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(db.Error)
		} else {
			return c.JSON(user)
		}
	} else {
		return invalid(c)
	}
}

func invalid(c *fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"message": "Invalid credentials",
	})
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func createSignedToken(id uint, username string) (string, error) {
	expirationTime := time.Now().AddDate(0, 0, 1).Unix()

	// Create the Claims
	claims := jwt.MapClaims{
		"userId":   id,
		"username": username,
		"exp":      expirationTime,
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRETKEY")))

	return signedToken, err
}
