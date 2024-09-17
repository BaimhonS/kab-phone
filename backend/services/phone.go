package services

import (
	"io"
	"mime/multipart"

	"github.com/BaimhonS/kab-phone/configs"
	"github.com/BaimhonS/kab-phone/models"
	"github.com/BaimhonS/kab-phone/utils"
	"github.com/BaimhonS/kab-phone/validates"
	"github.com/gofiber/fiber/v2"

	"gorm.io/gorm"
)

type PhoneServiceImpl struct {
	DB *gorm.DB
}

type PhoneService interface {
	GetPhones(c *fiber.Ctx) error
	GetPhoneImageByID(c *fiber.Ctx) error
	CreatePhone(c *fiber.Ctx) error
	UpdatePhone(c *fiber.Ctx) error
	DeletePhone(c *fiber.Ctx) error
}

func NewPhoneService(configClients configs.ConfigClients) PhoneService {
	return &PhoneServiceImpl{
		DB: configClients.DB,
	}
}

func (s *PhoneServiceImpl) GetPhones(c *fiber.Ctx) error {
	var query utils.QueryPagination
	if err := c.QueryParser(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Message: "query parser error",
			Error:   err,
		})
	}

	var phones []models.Phone

	queryPhones := s.DB

	if query.Search != "" {
		queryPhones = queryPhones.Where("model_name LIKE ?", "%"+query.Search+"%")
	}

	if err := queryPhones.Offset(query.Page * query.PageSize).Limit(query.PageSize).Find(&phones).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "database get phones error",
			Error:   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessPaginationResponse{
		Message: "get phones success",
		Data:    phones,
		Total:   len(phones),
	})
}

func (s *PhoneServiceImpl) GetPhoneImageByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var phone models.Phone
	if err := s.DB.First(&phone, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "database get phone error",
			Error:   err,
		})
	}

	c.Set("Content-Type", "image/jpeg")

	return c.Send(phone.Image)
}

func (s *PhoneServiceImpl) CreatePhone(c *fiber.Ctx) error {
	req, ok := c.Locals("req").(validates.RequestCreatePhone)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Message: "locals req error",
			Error:   nil,
		})
	}

	file, ok := c.Locals("file").(*multipart.FileHeader)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Message: "locals file error",
			Error:   nil,
		})
	}

	fileData, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "file open error",
			Error:   err,
		})
	}
	defer fileData.Close()

	imgBytes, err := io.ReadAll(fileData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "file read error",
			Error:   err,
		})
	}

	phone := models.Phone{
		ModelName: req.ModelName,
		BrandName: req.BrandName,
		OS:        req.OS,
		Price:     req.Price,
		Amount:    req.Amount,
		Image:     imgBytes,
	}

	if err := s.DB.Create(&phone).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "database create phone error",
			Error:   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse{
		Message: "create phone success",
		Data:    phone,
	})
}

func (s *PhoneServiceImpl) UpdatePhone(c *fiber.Ctx) error {
	req, ok := c.Locals("req").(validates.RequestUpdatePhone)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Message: "locals req error",
			Error:   nil,
		})
	}

	file, ok := c.Locals("file").(*multipart.FileHeader)
	var imgBytes []byte
	if ok {
		fileData, err := file.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
				Message: "file open error",
				Error:   err,
			})
		}
		defer fileData.Close()

		imgBytes, err = io.ReadAll(fileData)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
				Message: "file read error",
				Error:   err,
			})
		}
	} else {
		var phone models.Phone
		if err := s.DB.Model(&models.Phone{}).Where("id = ?", c.Params("id")).First(&phone).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
				Message: "phone not found error",
				Error:   err,
			})
		}
		imgBytes = phone.Image
	}

	if err := s.DB.Model(&models.Phone{}).Where("id = ?", c.Params("id")).Updates(models.Phone{
		Price:  req.Price,
		Amount: req.Amount,
		Image:  imgBytes,
	}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "database update phone error",
			Error:   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse{
		Message: "update phone success",
		Data:    nil,
	})
}

func (s *PhoneServiceImpl) DeletePhone(c *fiber.Ctx) error {
	if err := s.DB.Delete(&models.Phone{}, c.Params("id")).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "database delete phone error",
			Error:   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse{
		Message: "delete phone success",
		Data:    nil,
	})
}
