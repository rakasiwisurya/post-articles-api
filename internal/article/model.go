package article

import "time"

// Allowed post statuses as defined by the test spec.
const (
	StatusPublish = "publish"
	StatusDraft   = "draft"
	StatusThrash  = "thrash"
)

// Article maps to a row in the posts table.
type Article struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Category    string    `json:"category"`
	CreatedDate time.Time `json:"created_date"`
	UpdatedDate time.Time `json:"updated_date"`
	Status      string    `json:"status"`
}

// UpsertRequest is the JSON payload for creating or updating an article.
type UpsertRequest struct {
	Title    string `json:"title" validate:"required,min=20"`
	Content  string `json:"content" validate:"required,min=200"`
	Category string `json:"category" validate:"required,min=3"`
	Status   string `json:"status" validate:"required,oneof=publish draft thrash"`
}
