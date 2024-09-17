package middlewares

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/BaimhonS/kab-phone/models"
	"github.com/BaimhonS/kab-phone/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

var bypassPaths = map[string][]string{
	"POST": {
		"/api/users/login",
		"/api/users/register",
	},
	"GET": {
		"/api/phones",
		"/api/phones/images/*",
	},
}

func AuthToken(redisClient *redis.Client) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if checkByPassPath(c.Method(), c.Path()) {
			return c.Next()
		}

		tokenString := c.Get("Authorization")

		if !strings.Contains(tokenString, "Bearer") {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse{
				Message: "token type invalid",
				Error:   nil,
			})
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "unexpected signing method")
			}

			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse{
				Message: "token parse error",
				Error:   err,
			})
		}

		if !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse{
				Message: "token invalid",
				Error:   nil,
			})
		}

		claims := token.Claims.(jwt.MapClaims)

		var userClaim models.User
		if err := utils.ParseClaimToUserModel(claims, &userClaim); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
				Message: "parse claim to model user error",
				Error:   err,
			})
		}

		var userRedis models.User
		userRawRedis, err := redisClient.Get(c.Context(), fmt.Sprintf("user_id:%v", userClaim.ID)).Result()
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse{
				Message: "get user from redis error",
				Error:   err,
			})
		}

		if err := json.Unmarshal([]byte(userRawRedis), &userRedis); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
				Message: "unmarshal user from redis error",
				Error:   err,
			})
		}

		c.Locals("user", userRedis)

		return c.Next()
	}
}

func checkByPassPath(method, path string) bool {
	if paths, ok := bypassPaths[method]; ok {
		for _, p := range paths {
			if strings.Contains(p, "*") {
				bypassPathCustom := strings.Split(p, "/")
				pathCustom := strings.Split(path, "/")
				for i, _ := range bypassPathCustom {
					if bypassPathCustom[i] == "*" && len(bypassPathCustom) == len(pathCustom) {
						bypassPathCustom = append(bypassPathCustom[:i], bypassPathCustom[i+1:]...)
						pathCustom = append(pathCustom[:i], pathCustom[i+1:]...)
					}
				}

				joinByPassPath := strings.Join(bypassPathCustom, "/")
				joinPath := strings.Join(pathCustom, "/")

				if joinByPassPath == joinPath {
					return true
				}
			} else {
				if p == path {
					return true
				}
			}
		}
	}

	return false
}
