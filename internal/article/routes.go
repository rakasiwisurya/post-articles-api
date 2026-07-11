package article

import "github.com/gofiber/fiber/v3"

// RegisterRoutes mounts the article endpoints required by the test spec.
func RegisterRoutes(app *fiber.App, h *Handler) {
	group := app.Group("/article")

	group.Post("/", h.Create)
	group.Get("/:limit/:offset", h.List)
	group.Get("/:id", h.GetByID)
	group.Put("/:id", h.Update)
	group.Patch("/:id", h.Update)
	group.Post("/:id", h.Update) // spec allows POST, PUT or PATCH for updates
	group.Delete("/:id", h.Delete)
}
