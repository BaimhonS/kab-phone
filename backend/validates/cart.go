package validates

import (
	"github.com/BaimhonS/kab-phone/models"
	"github.com/BaimhonS/kab-phone/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type (
	RequestAddItemToCart struct {
		Item models.Item `json:"item" validate:"required"`
	}

	RequestUpdateItemFromCart struct {
		Amount int `json:"amount" validate:"required"`
	}

	CartValidateImpl struct{}
)

type CartValidate interface {
	ValidateAddItemToCart(c *fiber.Ctx) error
	ValidateUpdateitemFromCart(c *fiber.Ctx) error
}

func NewCartValidate() CartValidate {
	return &CartValidateImpl{}
}

func (v *CartValidateImpl) ValidateAddItemToCart(c *fiber.Ctx) error {
	var req RequestAddItemToCart
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Message: "body parser error",
			Error:   err,
		})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ValidateErrorResponse{
			Message: "validate add product to cart error",
			Error:   utils.HanddleValidateError(err),
		})
	}

	c.Locals("req", req)
	return c.Next()
}

func (v *CartValidateImpl) ValidateUpdateitemFromCart(c *fiber.Ctx) error {
	var req RequestUpdateItemFromCart
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Message: "body parser error",
			Error:   err,
		})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ValidateErrorResponse{
			Message: "validate update product from cart error",
			Error:   utils.HanddleValidateError(err),
		})
	}

	c.Locals("req", req)
	return c.Next()
}
