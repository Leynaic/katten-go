package handlers

import (
	"github.com/Leynaic/katten-go/database"
	"github.com/Leynaic/katten-go/helpers"
	"github.com/Leynaic/katten-go/models"
	"github.com/gofiber/fiber/v2"
)

func GetCats(c *fiber.Ctx) error {
	filter := c.Query("filter", "none")
	currentCat := helpers.GetCurrentCat(c)
	db := database.GetInstance()
	likeIds := db.Table("cat_likes").Where("cat_id = ?", currentCat.ID).Select("like_id")
	dislikeIds := db.Table("cat_dislikes").Where("cat_id = ?", currentCat.ID).Select("dislike_id")
	var cats []models.Cat
	switch filter {
	case "likes":
		db.Model(&models.Cat{}).
			Where("id IN (?)", likeIds).
			Find(&cats)
	case "dislikes":
		db.Model(&models.Cat{}).
			Where("id IN (?)", dislikeIds).
			Find(&cats)
	default:
		db.Model(&models.Cat{}).
			Where("id != ?", currentCat.ID).
			Where("id NOT IN (?)", likeIds).
			Where("id NOT IN (?)", dislikeIds).
			Find(&cats)
	}

	return c.Status(fiber.StatusOK).JSON(cats)
}

func LikeCat(c *fiber.Ctx) error {
	currentCat := helpers.GetCurrentCat(c)
	if catToLikeID, err := c.ParamsInt("id", 0); err == nil && int(currentCat.ID) != catToLikeID {
		db := database.GetInstance()
		if db.Exec("INSERT INTO cat_likes VALUES (?, ?)", currentCat.ID, catToLikeID).Error == nil {
			return c.Status(fiber.StatusOK).JSON(true)
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(false)
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(false)
	}
}

func CancelLikeCat(c *fiber.Ctx) error {
	currentCat := helpers.GetCurrentCat(c)
	if catToLikeID, err := c.ParamsInt("id", 0); err == nil && int(currentCat.ID) != catToLikeID {
		db := database.GetInstance()
		if db.Exec("DELETE FROM cat_likes WHERE cat_id = ? AND like_id = ?", currentCat.ID, catToLikeID).Error == nil {
			return c.Status(fiber.StatusOK).JSON(true)
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(false)
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(false)
	}
}

func DislikeCat(c *fiber.Ctx) error {
	currentCat := helpers.GetCurrentCat(c)
	if catToDislikeID, err := c.ParamsInt("id", 0); err == nil && int(currentCat.ID) != catToDislikeID {
		db := database.GetInstance()
		if db.Exec("INSERT INTO cat_dislikes VALUES (?, ?)", currentCat.ID, catToDislikeID).Error == nil {
			return c.Status(fiber.StatusOK).JSON(true)
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(false)
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(false)
	}
}
