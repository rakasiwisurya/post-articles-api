package article

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidationError carries field-level messages for a 400 response.
type ValidationError struct {
	Fields map[string]string
}

func (e *ValidationError) Error() string {
	messages := make([]string, 0, len(e.Fields))
	for field, message := range e.Fields {
		messages = append(messages, field+": "+message)
	}
	return "validation failed: " + strings.Join(messages, "; ")
}

// Service holds the article business rules on top of the repository.
type Service interface {
	Create(ctx context.Context, req UpsertRequest) (Article, error)
	List(ctx context.Context, limit, offset int, status string) ([]Article, error)
	Count(ctx context.Context, status string) (int64, error)
	Get(ctx context.Context, id int64) (Article, error)
	Update(ctx context.Context, id int64, req UpsertRequest) (Article, error)
	Delete(ctx context.Context, id int64) error
}

type service struct {
	repo     Repository
	validate *validator.Validate
}

// NewService creates the article service with the spec validation rules.
func NewService(repo Repository) Service {
	return &service{
		repo:     repo,
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (s *service) Create(ctx context.Context, req UpsertRequest) (Article, error) {
	if err := s.validateRequest(req); err != nil {
		return Article{}, err
	}
	return s.repo.Create(ctx, req)
}

func (s *service) List(ctx context.Context, limit, offset int, status string) ([]Article, error) {
	return s.repo.List(ctx, limit, offset, status)
}

func (s *service) Count(ctx context.Context, status string) (int64, error) {
	return s.repo.Count(ctx, status)
}

func (s *service) Get(ctx context.Context, id int64) (Article, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) Update(ctx context.Context, id int64, req UpsertRequest) (Article, error) {
	if err := s.validateRequest(req); err != nil {
		return Article{}, err
	}
	return s.repo.Update(ctx, id, req)
}

func (s *service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// validateRequest applies the spec rules and converts validator output into
// readable, field-keyed messages.
func (s *service) validateRequest(req UpsertRequest) error {
	err := s.validate.Struct(req)
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	fields := make(map[string]string, len(validationErrors))
	for _, fieldError := range validationErrors {
		field := strings.ToLower(fieldError.Field())
		switch fieldError.Tag() {
		case "required":
			fields[field] = field + " is required"
		case "min":
			fields[field] = fmt.Sprintf("%s must be at least %s characters", field, fieldError.Param())
		case "oneof":
			fields[field] = field + " must be one of: publish, draft, thrash"
		default:
			fields[field] = field + " is invalid"
		}
	}
	return &ValidationError{Fields: fields}
}
