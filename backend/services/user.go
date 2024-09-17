package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/BaimhonS/kab-phone/configs"
	"github.com/BaimhonS/kab-phone/models"
	"github.com/BaimhonS/kab-phone/utils"
	"github.com/BaimhonS/kab-phone/validates"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserServiceImpl struct {
	DB    *gorm.DB
	Redis *redis.Client
}

type UserService interface {
	GetProfileUser(c *fiber.Ctx) error
	RegisterUser(c *fiber.Ctx) error
	LoginUser(c *fiber.Ctx) error
	LogoutUser(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
}

func NewUserService(configClients configs.ConfigClients) UserService {
	return &UserServiceImpl{
		DB:    configClients.DB,
		Redis: configClients.Redis,
	}
}

func (s *UserServiceImpl) GetProfileUser(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Message: "local user not found",
			Error:   nil,
		})
	}

	var userInfo models.User
	if err := s.DB.Where("id = ?", user.ID).First(&userInfo).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "database get user error",
			Error:   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse{
		Message: "get profile user success",
		Data:    userInfo,
	})
}

func (s *UserServiceImpl) RegisterUser(c *fiber.Ctx) error {
	req, ok := c.Locals("req").(validates.RequestRegisterUser)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Message: "local register request not found",
			Error:   nil,
		})
	}

	var hasUser int64
	if err := s.DB.Model(&models.User{}).Where("username = ?", req.Username).Count(&hasUser).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "database check user error",
			Error:   err,
		})
	}

	if hasUser > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Message: "username already exists",
			Error:   nil,
		})
	}

	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "password hash error",
			Error:   err,
		})
	}

	user := models.User{
		Username:    req.Username,
		Password:    string(passwordHashed),
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		LineID:      req.LineID,
		Address:     req.Address,
		Role:        "guess",
		Age:         req.Age,
		BirthDate:   req.BirthDate,
	}

	tx := s.DB.Begin()
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "database create user error",
			Error:   err,
		})
	}

	cart := models.Cart{
		UserId: user.ID,
		Status: "PENDING",
	}

	if err := tx.Create(&cart).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "database create cart error",
			Error:   err,
		})
	}

	tx.Commit()

	user.Password = ""

	return c.Status(fiber.StatusCreated).JSON(utils.SuccessResponse{
		Message: "register user success",
		Data:    user,
	})
}

func (s *UserServiceImpl) LoginUser(c *fiber.Ctx) error {
	req, ok := c.Locals("req").(validates.RequestLoginUser)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Message: "local login request not found",
			Error:   nil,
		})
	}

	var user models.User
	if err := s.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse{
				Message: "username not found",
				Error:   nil,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "database get user error",
			Error:   err,
		})
	}

	if user.Role == "guess" {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse{
				Message: "password not match",
				Error:   nil,
			})
		}
	}

	if user.Role == "admin" {
		if req.Password != user.Password {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse{
				Message: "password not match",
				Error:   nil,
			})
		}
	}

	userClaim := map[string]interface{}{
		"id":           user.ID,
		"first_name":   user.FirstName,
		"last_name":    user.LastName,
		"username":     user.Username,
		"age":          user.Age,
		"birth_date":   user.BirthDate.Unix(),
		"phone_number": user.PhoneNumber,
		"line_id":      user.LineID,
		"address":      user.Address,
		"role":         user.Role,
	}

	accesstoken, err := utils.GenerateToken(userClaim, time.Now().Add(30*time.Minute).Unix())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "generate token error",
			Error:   err,
		})
	}

	user.Password = ""

	jsonUser, err := json.Marshal(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "marshal user error",
			Error:   err,
		})
	}

	if err := s.Redis.Set(c.Context(), fmt.Sprintf("user_id:%v", user.ID), jsonUser, 30*time.Minute).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "set redis error",
			Error:   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(utils.LoginResponse{
		AccesssToken: accesstoken,
	})
}

func (s *UserServiceImpl) LogoutUser(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Message: "local user not found",
			Error:   nil,
		})
	}

	if err := s.Redis.Del(c.Context(), fmt.Sprintf("user_id:%v", user.ID)).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "delete redis error",
			Error:   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse{
		Message: "logout user success",
		Data:    nil,
	})
}

func (s *UserServiceImpl) UpdateUser(c *fiber.Ctx) error {
	req, ok := c.Locals("req").(validates.RequestUpdateUser)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Message: "local update request not found",
			Error:   nil,
		})
	}

	userLocal, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Message: "local user not found",
			Error:   nil,
		})
	}

	var user models.User
	if err := s.DB.Where("id = ?", userLocal.ID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "database get user error",
			Error:   err,
		})
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.PhoneNumber = req.PhoneNumber
	user.LineID = req.LineID
	user.Address = req.Address
	user.Age = req.Age
	user.BirthDate = req.BirthDate

	if err := s.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "database update user error",
			Error:   err,
		})
	}

	user.Password = ""

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse{
		Message: "update user success",
		Data:    user,
	})
}
