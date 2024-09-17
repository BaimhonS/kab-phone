package utils

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"regexp"
	"strconv"

	"github.com/go-playground/validator/v10"
)

var AllowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
}

func ValidatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	lowercaseLetter := regexp.MustCompile(`[a-z]`)
	uppercaseLetter := regexp.MustCompile(`[A-Z]`)
	digit := regexp.MustCompile(`[\d]`)
	specialCharacter := regexp.MustCompile(`[!\"#$%&'()*+,\-./:;<=>?@[\\\]^_` + "`" + `{|}~]`)

	return lowercaseLetter.MatchString(password) &&
		uppercaseLetter.MatchString(password) &&
		digit.MatchString(password) &&
		specialCharacter.MatchString(password)
}

func HanddleValidateError(err error) (errs []*ValidateError) {
	for _, err := range err.(validator.ValidationErrors) {
		errs = append(errs, &ValidateError{
			Field: err.Field(),
			Tag:   err.Tag(),
			Value: err.Param(),
		})
	}

	return errs
}

func ValidateImageFile(fileHeader *multipart.FileHeader) error {
	maxFileSize, err := strconv.Atoi(os.Getenv("MAX_IMAGE_SIZE"))
	if err != nil {
		return err
	}
	fmt.Println("test1")

	if fileHeader.Size > int64(maxFileSize*1024*1024) {
		return errors.New("file size exceeds the 5MB limit")
	}
	fmt.Println("test2")

	isAllowed, ok := AllowedImageTypes[fileHeader.Header.Get("Content-Type")]

	if !ok || !isAllowed {
		return errors.New("file type is not allowed")
	}

	fmt.Println("test3")

	return nil
}
