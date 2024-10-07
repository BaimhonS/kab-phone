package services

import (
	"fmt"
	"time"

	"github.com/BaimhonS/kab-phone/configs"
	"github.com/BaimhonS/kab-phone/models"
	"github.com/BaimhonS/kab-phone/utils"
	"github.com/gofiber/fiber/v2"

	"gorm.io/gorm"
)

type OrderServiceImpl struct {
	DB *gorm.DB
}

type OrderService interface {
	ConfirmOrder(c *fiber.Ctx) error
	GetOrderByTrackingNumber(c *fiber.Ctx) error
	GetTrackingNumbers(c *fiber.Ctx) error
	GetBestAndWorstSellingPhones(c *fiber.Ctx) error
	GetTotalIncome(c *fiber.Ctx) error
	GetAllOrders(c *fiber.Ctx) error
}

func NewOrderService(configClients configs.ConfigClients) OrderService {
	return &OrderServiceImpl{
		DB: configClients.DB,
	}
}

func (s *OrderServiceImpl) ConfirmOrder(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Message: "local user not found",
		})
	}

	var cart models.Cart
	if err := s.DB.Model(&models.Cart{}).Where("user_id = ? AND status = ?", user.ID, "PENDING").Preload("Items").Preload("Items.Phone").First(&cart).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "cart not found",
			Error:   err,
		})
	}

	cart.Status = "CONFIRMED"

	tx := s.DB.Begin()

	if err := tx.Save(&cart).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "update cart failed",
			Error:   err,
		})
	}

	var totalPrice float32
	for _, item := range cart.Items {
		totalPrice += item.Phone.Price * float32(item.Amount)

		if err := tx.Model(&models.Phone{}).Where("id = ?", item.Phone.ID).Select("Amount").Updates(models.Phone{
			Amount: item.Phone.Amount - item.Amount,
		}).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
				Message: "update phone amount failed",
				Error:   err,
			})
		}
	}

	var order models.Order
	order.CartID = cart.ID
	order.TrackingNumber = fmt.Sprintf("TH-%s", utils.GenerateNumericString(10))
	order.TotalPrice = totalPrice

	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "create new order failed",
			Error:   err,
		})
	}

	if err := tx.Create(&models.Cart{
		UserId: user.ID,
		Status: "PENDING",
	}).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "create new cart failed",
			Error:   err,
		})
	}

	tx.Commit()

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse{
		Message: "confirm order success",
		Data:    nil,
	})
}

func (s *OrderServiceImpl) GetAllOrders(c *fiber.Ctx) error {
	var orders []models.Order
	// Fetch all orders, including related data such as Cart, Items, and Phones
	if err := s.DB.Preload("Cart").Preload("Cart.Items").Preload("Cart.Items.Phone").Preload("Cart.User").Find(&orders).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "failed to fetch all orders",
			Error:   err,
		})
	}

	// Return success response with all the orders
	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse{
		Message: "get all orders success",
		Data:    orders,
	})
}

func (s *OrderServiceImpl) GetOrderByTrackingNumber(c *fiber.Ctx) error {
	var order models.Order
	if err := s.DB.Model(&models.Order{}).Where("tracking_number = ?", c.Params("tracking_number")).Preload("Cart").Preload("Cart.Items").Preload("Cart.Items.Phone").First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse{
				Message: "order not found",
				Error:   err,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "get order by tracking number error",
			Error:   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse{
		Message: "get order by tracking number success",
		Data:    order,
	})
}

func (s *OrderServiceImpl) GetTrackingNumbers(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Message: "local user not found",
		})
	}

	var query utils.QueryPagination
	if err := c.QueryParser(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Message: "query parser error",
			Error:   err,
		})
	}

	queryOrders := s.DB.Model(&models.Order{}).Joins("JOIN carts ON orders.cart_id = carts.id").Where("carts.user_id = ?", user.ID)

	if query.Search != "" {
		queryOrders.Where("tracking_number LIKE ?", "%"+query.Search+"%")
	}

	var orders []models.Order
	if err := queryOrders.Offset(query.Page * query.PageSize).Limit(query.PageSize).Order("created_at DESC").Find(&orders).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "database get phones error",
			Error:   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse{
		Message: "get tracking numbers success",
		Data:    orders,
	})
}

func (s *OrderServiceImpl) GetBestAndWorstSellingPhones(c *fiber.Ctx) error {
	var bestPhoneDay models.Phone
	var bestPhoneWeek models.Phone
	var bestPhoneMonth models.Phone
	var bestPhoneYear models.Phone
	var worstPhoneDay models.Phone
	var worstPhoneWeek models.Phone
	var worstPhoneMonth models.Phone
	var worstPhoneYear models.Phone

	// best phone selling of the day
	if err := queryGetBestPhone(utils.GetStartOfDay(), utils.GetEndOfDay(), &bestPhoneDay, s.DB); err != nil && err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "get best phone selling of the day error",
			Error:   err,
		})
	}

	// best phone selling of the week
	if err := queryGetBestPhone(utils.GetStartOfDay().Add(-7*24*time.Hour), utils.GetEndOfDay(), &bestPhoneWeek, s.DB); err != nil && err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "get best phone selling of the day error",
			Error:   err,
		})
	}

	// best phone selling of the month
	if err := queryGetBestPhone(utils.GetStartOfDay().Add(-30*24*time.Hour), utils.GetEndOfDay(), &bestPhoneMonth, s.DB); err != nil && err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "get best phone selling of the day error",
			Error:   err,
		})
	}

	// best phone selling of the year
	if err := queryGetBestPhone(utils.GetStartOfDay().Add(-365*24*time.Hour), utils.GetEndOfDay(), &bestPhoneYear, s.DB); err != nil && err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "get best phone selling of the day error",
			Error:   err,
		})
	}

	// worst phone selling of the day
	if err := queryGetWorstPhone(utils.GetStartOfDay(), utils.GetEndOfDay(), &worstPhoneDay, s.DB); err != nil && err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "get best phone selling of the day error",
			Error:   err,
		})
	}

	// worst phone selling of the week
	if err := queryGetWorstPhone(utils.GetStartOfDay().Add(-7*24*time.Hour), utils.GetEndOfDay(), &worstPhoneWeek, s.DB); err != nil && err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "get best phone selling of the day error",
			Error:   err,
		})
	}

	// worst phone selling of the month
	if err := queryGetWorstPhone(utils.GetStartOfDay().Add(-30*24*time.Hour), utils.GetEndOfDay(), &worstPhoneMonth, s.DB); err != nil && err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "get best phone selling of the day error",
			Error:   err,
		})
	}

	// worst phone selling of the year
	if err := queryGetWorstPhone(utils.GetStartOfDay().Add(-365*24*time.Hour), utils.GetEndOfDay(), &worstPhoneYear, s.DB); err != nil && err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "get best phone selling of the day error",
			Error:   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse{
		Message: "get best and worst selling phones success",
		Data: fiber.Map{
			"best_phone_day": fiber.Map{
				"brand_name": bestPhoneDay.BrandName,
				"model_name": bestPhoneDay.ModelName,
			},
			"best_phone_week": fiber.Map{
				"brand_name": bestPhoneWeek.BrandName,
				"model_name": bestPhoneWeek.ModelName,
			},
			"best_phone_month": fiber.Map{
				"brand_name": bestPhoneMonth.BrandName,
				"model_name": bestPhoneMonth.ModelName,
			},
			"best_phone_year": fiber.Map{
				"brand_name": bestPhoneYear.BrandName,
				"model_name": bestPhoneYear.ModelName,
			},
			"worst_phone_day": fiber.Map{
				"brand_name": worstPhoneDay.BrandName,
				"model_name": worstPhoneDay.ModelName,
			},
			"worst_phone_week": fiber.Map{
				"brand_name": worstPhoneWeek.BrandName,
				"model_name": worstPhoneWeek.ModelName,
			},
			"worst_phone_month": fiber.Map{
				"brand_name": worstPhoneMonth.BrandName,
				"model_name": worstPhoneMonth.ModelName,
			},
			"worst_phone_year": fiber.Map{
				"brand_name": worstPhoneYear.BrandName,
				"model_name": worstPhoneYear.ModelName,
			},
		},
	})
}

func queryGetBestPhone(startDate time.Time, endDate time.Time, phone *models.Phone, db *gorm.DB) error {
	return db.Model(&models.Item{}).
		Select("SUM(items.amount) as amount_sell, phones.*").
		Joins("JOIN carts ON items.cart_id = carts.id").
		Joins("JOIN phones ON phones.id = items.phone_id").
		Joins("JOIN orders ON orders.cart_id = carts.id").
		Where("orders.created_at >= ? AND orders.created_at <= ?", startDate, endDate).
		Group("phones.id").
		Order("amount_sell DESC").
		Limit(1).
		Scan(&phone).Error
}

func queryGetWorstPhone(startDate time.Time, endDate time.Time, phone *models.Phone, db *gorm.DB) error {
	return db.Model(&models.Item{}).
		Select("SUM(items.amount) as amount_sell, phones.*").
		Joins("JOIN carts ON items.cart_id = carts.id").
		Joins("JOIN phones ON phones.id = items.phone_id").
		Joins("JOIN orders ON orders.cart_id = carts.id").
		Where("orders.created_at >= ? AND orders.created_at <= ?", startDate, endDate).
		Group("phones.id").
		Order("amount_sell ASC").
		Limit(1).
		Scan(&phone).Error
}

type TotalIncome struct {
	CreatedAt time.Time
	Amount    float64
	Price     float64
}

func (s *OrderServiceImpl) GetTotalIncome(c *fiber.Ctx) error {
	var totalIncomeDay []TotalIncome
	var totalIncomeWeek []TotalIncome
	var totalIncomeMonth []TotalIncome
	var totalIncomeYear []TotalIncome

	// get total price selling of the day
	if err := queryGetTotalIncome(utils.GetStartOfDay(), utils.GetEndOfDay(), &totalIncomeDay, s.DB); err != nil && err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "get best phone selling of the day error",
			Error:   err,
		})
	}

	// get total price selling of the week
	if err := queryGetTotalIncome(utils.GetStartOfDay().Add(-7*24*time.Hour), utils.GetEndOfDay(), &totalIncomeWeek, s.DB); err != nil && err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "get best phone selling of the day error",
			Error:   err,
		})
	}

	// get total price selling of the month
	if err := queryGetTotalIncome(utils.GetStartOfDay().Add(-30*24*time.Hour), utils.GetEndOfDay(), &totalIncomeMonth, s.DB); err != nil && err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "get best phone selling of the day error",
			Error:   err,
		})
	}

	// get total price selling of the year
	if err := queryGetTotalIncome(utils.GetStartOfDay().Add(-365*24*time.Hour), utils.GetEndOfDay(), &totalIncomeYear, s.DB); err != nil && err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Message: "get best phone selling of the day error",
			Error:   err,
		})
	}

	totalIncomeDayAmount := 0.0
	totalIncomeWeekAmount := 0.0
	totalIncomeMonthAmount := 0.0
	totalIncomeYearAmount := 0.0

	for _, income := range totalIncomeDay {
		totalIncomeDayAmount += income.Amount * income.Price
	}

	for _, income := range totalIncomeWeek {
		totalIncomeWeekAmount += income.Amount * income.Price
	}

	for _, income := range totalIncomeMonth {
		totalIncomeMonthAmount += income.Amount * income.Price
	}

	for _, income := range totalIncomeYear {
		totalIncomeYearAmount += income.Amount * income.Price
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse{
		Message: "get total income success",
		Data: fiber.Map{
			"total_income_day":   totalIncomeDayAmount,
			"total_income_week":  totalIncomeWeekAmount,
			"total_income_month": totalIncomeMonthAmount,
			"total_income_year":  totalIncomeYearAmount,
		},
	})
}

func queryGetTotalIncome(startDate time.Time, endDate time.Time, totalIncome *[]TotalIncome, db *gorm.DB) error {
	return db.Model(&models.Item{}).
		Select("items.amount, phones.price, orders.created_at").
		Joins("JOIN carts ON items.cart_id = carts.id").
		Joins("JOIN phones ON phones.id = items.phone_id").
		Joins("JOIN orders ON orders.cart_id = carts.id").
		Where("orders.created_at >= ? AND orders.created_at <= ?", startDate, endDate).
		Where("carts.status = ?", "CONFIRMED").
		Find(&totalIncome).Error
}
