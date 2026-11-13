package handler

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

var validate = newValidator()

// newValidator reports field names using their json tag ("email", not "Email")
func newValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterTagNameFunc(func(f reflect.StructField) string {
		return strings.Split(f.Tag.Get("json"), ",")[0]
	})
	return v
}

// parseBody reads JSON into out and validates it; any problem becomes a 400
func parseBody(c fiber.Ctx, out any) error {
	if err := c.Bind().Body(out); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid JSON body")
	}

	var verrs validator.ValidationErrors
	if err := validate.Struct(out); errors.As(err, &verrs) {
		msgs := make([]string, 0, len(verrs))
		for _, fe := range verrs {
			msgs = append(msgs, fieldMessage(fe))
		}
		return fiber.NewError(fiber.StatusBadRequest, strings.Join(msgs, "; "))
	} else if err != nil {
		return err
	}
	return nil
}

// fieldMessage turns a validator error into a user-friendly sentence
func fieldMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email", fe.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", fe.Field(), fe.Param())
	default:
		return fmt.Sprintf("%s is invalid", fe.Field())
	}
}
