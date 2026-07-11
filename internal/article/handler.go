package article

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// Handler exposes the article service over HTTP.
type Handler struct {
	service Service
}

// NewHandler creates the article HTTP handler.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Create handles POST /article/
func (h *Handler) Create(c fiber.Ctx) error {
	var req UpsertRequest
	if err := c.Bind().Body(&req); err != nil {
		return badRequest(c, "request body must be valid JSON")
	}

	created, err := h.service.Create(c.Context(), req)
	if err != nil {
		return respondError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(created)
}

// List handles GET /article/:limit/:offset with an optional ?status= filter.
func (h *Handler) List(c fiber.Ctx) error {
	limit, err := strconv.Atoi(c.Params("limit"))
	if err != nil || limit < 1 {
		return badRequest(c, "limit must be a positive integer")
	}
	offset, err := strconv.Atoi(c.Params("offset"))
	if err != nil || offset < 0 {
		return badRequest(c, "offset must be a non-negative integer")
	}

	status := c.Query("status")
	if status != "" && status != StatusPublish && status != StatusDraft && status != StatusThrash {
		return badRequest(c, "status must be one of: publish, draft, thrash")
	}

	articles, err := h.service.List(c.Context(), limit, offset, status)
	if err != nil {
		return respondError(c, err)
	}

	total, err := h.service.Count(c.Context(), status)
	if err != nil {
		return respondError(c, err)
	}

	// Total exposed as a header so the response body stays the plain array
	// required by the test spec while clients can still paginate properly.
	c.Set("X-Total-Count", strconv.FormatInt(total, 10))
	return c.JSON(articles)
}

// GetByID handles GET /article/:id
func (h *Handler) GetByID(c fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return badRequest(c, err.Error())
	}

	article, err := h.service.Get(c.Context(), id)
	if err != nil {
		return respondError(c, err)
	}
	return c.JSON(article)
}

// Update handles POST/PUT/PATCH /article/:id
func (h *Handler) Update(c fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return badRequest(c, err.Error())
	}

	var req UpsertRequest
	if err := c.Bind().Body(&req); err != nil {
		return badRequest(c, "request body must be valid JSON")
	}

	updated, err := h.service.Update(c.Context(), id, req)
	if err != nil {
		return respondError(c, err)
	}
	return c.JSON(updated)
}

// Delete handles DELETE /article/:id
func (h *Handler) Delete(c fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return badRequest(c, err.Error())
	}

	if err := h.service.Delete(c.Context(), id); err != nil {
		return respondError(c, err)
	}
	return c.JSON(fiber.Map{"message": "article deleted"})
}

func parseID(c fiber.Ctx) (int64, error) {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("id must be a positive integer")
	}
	return id, nil
}

func badRequest(c fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": message})
}

// respondError maps domain errors to HTTP responses.
func respondError(c fiber.Ctx, err error) error {
	var validationErr *ValidationError
	switch {
	case errors.As(err, &validationErr):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "validation failed",
			"errors":  validationErr.Fields,
		})
	case errors.Is(err, ErrNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "article not found"})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "internal server error"})
	}
}
