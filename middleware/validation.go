package middleware

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	validate.RegisterValidation("no_html", ValidateNoHTML)
	validate.RegisterValidation("safe_string", ValidateSafeString)
	validate.RegisterValidation("uuid_or_empty", ValidateUUIDOrEmpty)
}

func ValidateNoHTML(fl validator.FieldLevel) bool {
	htmlRegex := regexp.MustCompile(`<[^>]+>`)
	return !htmlRegex.MatchString(fl.Field().String())
}

func ValidateSafeString(fl validator.FieldLevel) bool {
	dangerousChars := regexp.MustCompile(`[<>'"&;(){}[\]\\]`)
	return !dangerousChars.MatchString(fl.Field().String())
}

func ValidateUUIDOrEmpty(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true
	}
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	return uuidRegex.MatchString(value)
}

func ValidateStruct() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body interface{}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": "Dados inválidos no corpo da requisição",
			})
		}

		if body == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": "Corpo da requisição obrigatório",
			})
		}

		return c.Next()
	}
}

func SanitizeInput(input string) string {
	htmlRegex := regexp.MustCompile(`<[^>]*>`)
	sanitized := htmlRegex.ReplaceAllString(input, "")

	scriptRegex := regexp.MustCompile(`(?i)(javascript:|data:|vbscript:|onclick|onerror|onload)`)
	sanitized = scriptRegex.ReplaceAllString(sanitized, "")

	sanitized = strings.TrimSpace(sanitized)

	return sanitized
}

func ValidateAndSanitize(s interface{}) error {
	sanitizeStruct(s)

	if err := validate.Struct(s); err != nil {
		return err
	}

	return nil
}

func sanitizeStruct(s interface{}) {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.String && field.CanSet() {
			sanitized := SanitizeInput(field.String())
			field.SetString(sanitized)
		} else if field.Kind() == reflect.Ptr && !field.IsNil() {
			if field.Elem().Kind() == reflect.String && field.Elem().CanSet() {
				sanitized := SanitizeInput(field.Elem().String())
				field.Elem().SetString(sanitized)
			}
		}
	}
}

func InputSizeLimit(maxSize int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if len(c.Body()) > maxSize {
			return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{
				"error":   true,
				"message": "Requisição muito grande",
			})
		}
		return c.Next()
	}
}
