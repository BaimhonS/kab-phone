package validates

import (
	"time"

	"github.com/BaimhonS/kab-phone/models"
	"github.com/BaimhonS/kab-phone/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type (
	RequestRegisterUser struct {
		Username        string    `json:"username" validate:"required,min=5,max=20,alphanum"`
		FirstName       string    `json:"first_name" validate:"required,min=3,max=50,alpha"`
		LastName        string    `json:"last_name" validate:"required,min=3,max=50,alpha"`
		PhoneNumber     string    `json:"phone_number" validate:"required,min=9,max=10,numeric"`
		Password        string    `json:"password" validate:"required,password"`
		ConfirmPassword string    `json:"confirm_password" validate:"required,eqfield=Password"`
		LineID          string    `json:"line_id"`
		Address         string    `json:"address" validate:"required,min=10,max=100"`
		Age             int       `json:"age" validate:"required,min=1,max=150"`
		BirthDate       time.Time `json:"birth_date" validate:"required"`
	}

	RequestLoginUser struct {
		Username string `json:"username" validate:"required,min=5,max=20,alphanum"`
		Password string `json:"password" validate:"required"`
	}

	RequestUpdateUser struct {
		FirstName   string    `json:"first_name" validate:"required,min=3,max=50,alpha"`
		LastName    string    `json:"last_name" validate:"required,min=3,max=50,alpha"`
		PhoneNumber string    `json:"phone_number" validate:"required,min=9,max=10,numeric"`
		LineID      string    `json:"line_id"`
		Address     string    `json:"address" validate:"required,min=10,max=100"`
		Age         int       `json:"age" validate:"required,min=1,max=150"`
		BirthDate   time.Time `json:"birth_date" validate:"required"`
	}

	UserValidateImpl struct{}
)

type UserValidate interface {
	ValidateRegisterUser(c *fiber.Ctx) error
	ValidateLoginUser(c *fiber.Ctx) error
	ValidateUpdateUser(c *fiber.Ctx) error
	ValidateRoleAdmin(c *fiber.Ctx) error
}

func NewUserValidate() UserValidate {
	return &UserValidateImpl{}
}

func (v *UserValidateImpl) ValidateRegisterUser(c *fiber.Ctx) error {
	var req RequestRegisterUser
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Message: "body parser error",
			Error:   err,
		})
	}

	validate := validator.New()
	validate.RegisterValidation("password", utils.ValidatePassword)
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ValidateErrorResponse{
			Message: "request register invalid",
			Error:   utils.HanddleValidateError(err),
		})
	}

	c.Locals("req", req)
	return c.Next()
}

func (v *UserValidateImpl) ValidateLoginUser(c *fiber.Ctx) error {
	var req RequestLoginUser
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Message: "body parser error",
			Error:   err,
		})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ValidateErrorResponse{
			Message: "request login invalid",
			Error:   utils.HanddleValidateError(err),
		})
	}

	c.Locals("req", req)
	return c.Next()
}

func (v *UserValidateImpl) ValidateUpdateUser(c *fiber.Ctx) error {
	var req RequestUpdateUser
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Message: "body parser error",
			Error:   err,
		})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ValidateErrorResponse{
			Message: "request update user invalid",
			Error:   utils.HanddleValidateError(err),
		})
	}

	c.Locals("req", req)
	return c.Next()
}

func (v *UserValidateImpl) ValidateRoleAdmin(c *fiber.Ctx) error {
	userLocal, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse{
			Message: "user not found",
			Error:   nil,
		})
	}

	if userLocal.Role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(utils.ErrorResponse{
			Message: "role invalid",
			Error:   nil,
		})
	}

	return c.Next()
}
