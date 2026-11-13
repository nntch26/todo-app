// Package handler is the HTTP transport layer: it parses requests and writes responses.
package handler

import "github.com/gofiber/fiber/v3"

// Health reports that the API process is up. Public: no auth required.
func Health(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}
