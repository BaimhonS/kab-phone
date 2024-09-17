package validates

import (
	"github.com/BaimhonS/kab-phone/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type (
	RequestCreatePhone struct {
		ModelName string  `form:"model_name" validate:"required,min=3,max=50"`
		BrandName string  `form:"brand_name" validate:"required,min=3,max=50"`
		OS        string  `form:"os" validate:"required,min=3,max=50,alpha"`
		Price     float32 `form:"price" validate:"required,min=0"`
		Amount    int     `form:"amount" validate:"required,min=0"`
		Image     []byte  `form:"image"`
	}

	RequestUpdatePhone struct {
		Price  float32 `form:"price" validate:"min=0"`
		Amount int     `form:"amount" validate:"min=0"`
		Image  []byte  `form:"image"`
	}

	PhoneValidateImpl struct{}
)

type PhoneValidate interface {
	ValidateCreatePhone(c *fiber.Ctx) error
	ValidateUpdatePhone(c *fiber.Ctx) error
}

func NewPhoneValidate() PhoneValidate {
	return &PhoneValidateImpl{}
}

func (v *PhoneValidateImpl) ValidateCreatePhone(c *fiber.Ctx) error {
	var req RequestCreatePhone
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Message: "body parser error",
			Error:   err,
		})
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Message: "form file error",
			Error:   err,
		})
	}

	if file == nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Message: "image is required",
			Error:   nil,
		})
	}

	if err := utils.ValidateImageFile(file); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Message: "validate file error",
			Error:   err,
		})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ValidateErrorResponse{
			Message: "validate create phone error",
			Error:   utils.HanddleValidateError(err),
		})
	}

	c.Locals("req", req)
	c.Locals("file", file)
	return c.Next()
}

func (v *PhoneValidateImpl) ValidateUpdatePhone(c *fiber.Ctx) error {
	var req RequestUpdatePhone
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Message: "body parser error",
			Error:   err,
		})
	}

	file, err := c.FormFile("image")
	if err == nil {
		if err := utils.ValidateImageFile(file); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
				Message: "validate file error",
				Error:   err,
			})
		}

		c.Locals("file", file)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ValidateErrorResponse{
			Message: "validate update phone error",
			Error:   utils.HanddleValidateError(err),
		})
	}

	c.Locals("req", req)
	return c.Next()
}
