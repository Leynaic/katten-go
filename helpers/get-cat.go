package helpers

import (
	"fmt"

	"github.com/Leynaic/katten-go/database"
	"github.com/Leynaic/katten-go/models"
	"github.com/Leynaic/katten-go/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func GetCurrentCat(c *fiber.Ctx) models.Cat {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := uint(claims["userId"].(float64))
	fmt.Println("ID : ", userId)
	db := database.GetInstance()
	var currentCat models.Cat
	db.Where(models.Cat{ID: userId}).First(&currentCat)
	if avatar, err := utils.GetUrl(currentCat.Avatar); err == nil {
		currentCat.Avatar = avatar.String()
	}
	return currentCat
}
