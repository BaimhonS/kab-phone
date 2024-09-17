package utils

import (
	"errors"
	"os"
	"time"

	"github.com/BaimhonS/kab-phone/models"
	"github.com/golang-jwt/jwt/v5"
)

type (
	UserClaim struct {
		ID          string  `json:"id"`
		FirstName   string  `json:"first_name"`
		LastName    string  `json:"last_name"`
		Username    string  `json:"username"`
		Age         float64 `json:"age"`
		BirthDate   float64 `json:"birth_date"`
		PhoneNumber string  `json:"phone_number"`
		LineID      string  `json:"line_id"`
		Address     string  `json:"address"`
		Role        string  `json:"role"`
	}
)

func GenerateToken(data map[string]interface{}, exp int64) (string, error) {
	claim := jwt.MapClaims(data)
	claim["exp"] = exp

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseClaimToUserModel(claim jwt.MapClaims, user *models.User) error {
	username, ok := claim["username"].(string)
	if !ok {
		return errors.New("username not found")
	}

	firstName, ok := claim["first_name"].(string)
	if !ok {
		return errors.New("first_name not found")
	}

	lastName, ok := claim["last_name"].(string)
	if !ok {
		return errors.New("last_name not found")
	}

	phoneNumber, ok := claim["phone_number"].(string)
	if !ok {
		return errors.New("phone_number not found")
	}

	lineID, ok := claim["line_id"].(string)
	if !ok {
		return errors.New("line_id not found")
	}

	address, ok := claim["address"].(string)
	if !ok {
		return errors.New("address not found")
	}

	age, ok := claim["age"].(float64)
	if !ok {
		return errors.New("age not found")
	}

	birthDate, ok := claim["birth_date"].(float64)
	if !ok {
		return errors.New("birth_date not found")
	}

	role, ok := claim["role"].(string)
	if !ok {
		return errors.New("role not found")
	}

	userID, ok := claim["id"].(float64)
	if !ok {
		return errors.New("id not found")
	}

	user.Username = username
	user.FirstName = firstName
	user.LastName = lastName
	user.PhoneNumber = phoneNumber
	user.LineID = lineID
	user.Address = address
	user.Age = int(age)
	user.BirthDate = time.Unix(int64(birthDate), 0)
	user.Role = role
	user.ID = uint(userID)

	return nil
}
