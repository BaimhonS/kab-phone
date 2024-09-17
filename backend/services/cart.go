package services

import (
	"fmt"

	"github.com/BaimhonS/kab-phone/configs"
	"github.com/BaimhonS/kab-phone/models"
	"github.com/BaimhonS/kab-phone/utils"
	"github.com/BaimhonS/kab-phone/validates"
	"github.com/gofiber/fiber/v2"

	"gorm.io/gorm"
)

type CartServiceImpl struct {
	DB *gorm.DB
}

type CartService interface {
	AddItemToCart(c *fiber.Ctx) error
	GetCart(c *fiber.Ctx) error
	RemoveItemFromCart(c *fiber.Ctx) error
	UpdateItemFromCart(c *fiber.Ctx) error
}

func NewCartService(configClients configs.ConfigClients) CartService {
	return &CartServiceImpl{
		DB: configClients.DB,
	}
}

func (s *CartServiceImpl) AddItemToCart(c *fiber.Ctx) error {
	req, ok := c.Locals("req").(validates.RequestAddItemToCart)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Message: "local req not found",
		})
	}

	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Message: "local user not found",
		})
	}

	var cart models.Cart
	if err := s.DB.Model(&models.Cart{}).Where("user_id = ? AND status = ?", user.ID, "PENDING").Preload("Items").First(&cart).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Message: "cart not found",
			Error:   err,
		})
	}

	isExist := false
	for i, item := range cart.Items {
		if item.PhoneID == req.Item.PhoneID {
			cart.Items[i].Amount += req.Item.Amount

			if err := s.DB.Model(&models.Item{}).Where("id = ?", item.ID).Update("amount", cart.Items[i].Amount).Error; err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
					Message: "update product to cart failed",
					Error:   err,
				})
			}

			isExist = true
			break
		}
	}

	if !isExist {
		cart.Items = append(cart.Items, req.Item)

		if err := s.DB.Save(&cart).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
				Message: "add product to cart failed",
				Error:   err,
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse{
		Message: "add or update product to cart success",
		Data:    cart,
	})
}

func (s *CartServiceImpl) GetCart(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Message: "local user not found",
		})
	}

	var cart models.Cart
	if err := s.DB.Model(&models.Cart{}).Where("user_id = ? AND status = ?", user.ID, "PENDING").Preload("Items").Preload("Items.Phone").First(&cart).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Message: "cart not found",
			Error:   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse{
		Message: "get cart success",
		Data:    cart,
	})
}

func (s *CartServiceImpl) RemoveItemFromCart(c *fiber.Ctx) error {
	if err := s.DB.Delete(&models.Item{}, c.Params("id")).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Message: "remove product from cart failed",
			Error:   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse{
		Message: "remove product from cart success",
		Data:    nil,
	})
}

func (s *CartServiceImpl) UpdateItemFromCart(c *fiber.Ctx) error {
	req, ok := c.Locals("req").(validates.RequestUpdateItemFromCart)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Message: "local req not found",
		})
	}

	var Item models.Item
	if err := s.DB.Model(&models.Item{}).Where("id = ?", c.Params("id")).Preload("Phone").First(&Item).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Message: "product not found",
			Error:   err,
		})
	}

	if Item.Phone.Amount < req.Amount {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Message: fmt.Sprintf("you can only buy max %v in one checkout", Item.Phone.Amount),
			Error:   nil,
		})
	}

	Item.Amount = req.Amount

	if err := s.DB.Save(&Item).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Message: "update product from cart failed",
			Error:   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse{
		Message: "update product from cart success",
		Data:    nil,
	})
}
