package article

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// ErrNotFound is returned when no article exists for the requested id.
var ErrNotFound = errors.New("article not found")

// Repository abstracts persistence so the service layer stays storage-agnostic.
type Repository interface {
	Create(ctx context.Context, req UpsertRequest) (Article, error)
	List(ctx context.Context, limit, offset int, status string) ([]Article, error)
	Count(ctx context.Context, status string) (int64, error)
	GetByID(ctx context.Context, id int64) (Article, error)
	Update(ctx context.Context, id int64, req UpsertRequest) (Article, error)
	Delete(ctx context.Context, id int64) error
}

type mysqlRepository struct {
	db *sql.DB
}

// NewRepository creates the MySQL-backed article repository.
func NewRepository(db *sql.DB) Repository {
	return &mysqlRepository{db: db}
}

const selectColumns = "id, title, content, category, created_date, updated_date, status"

func (r *mysqlRepository) Create(ctx context.Context, req UpsertRequest) (Article, error) {
	result, err := r.db.ExecContext(ctx,
		"INSERT INTO posts (title, content, category, status) VALUES (?, ?, ?, ?)",
		req.Title, req.Content, req.Category, req.Status,
	)
	if err != nil {
		return Article{}, fmt.Errorf("insert post: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Article{}, fmt.Errorf("read insert id: %w", err)
	}
	return r.GetByID(ctx, id)
}

func (r *mysqlRepository) List(ctx context.Context, limit, offset int, status string) ([]Article, error) {
	query := "SELECT " + selectColumns + " FROM posts"
	args := make([]any, 0, 3)
	if status != "" {
		query += " WHERE status = ?"
		args = append(args, status)
	}
	query += " ORDER BY updated_date DESC, id DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("select posts: %w", err)
	}
	defer rows.Close()

	articles := make([]Article, 0)
	for rows.Next() {
		article, err := scanArticle(rows)
		if err != nil {
			return nil, fmt.Errorf("scan post: %w", err)
		}
		articles = append(articles, article)
	}
	return articles, rows.Err()
}

func (r *mysqlRepository) Count(ctx context.Context, status string) (int64, error) {
	query := "SELECT COUNT(*) FROM posts"
	args := make([]any, 0, 1)
	if status != "" {
		query += " WHERE status = ?"
		args = append(args, status)
	}

	var total int64
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&total); err != nil {
		return 0, fmt.Errorf("count posts: %w", err)
	}
	return total, nil
}

func (r *mysqlRepository) GetByID(ctx context.Context, id int64) (Article, error) {
	row := r.db.QueryRowContext(ctx, "SELECT "+selectColumns+" FROM posts WHERE id = ?", id)

	article, err := scanArticle(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Article{}, ErrNotFound
	}
	if err != nil {
		return Article{}, fmt.Errorf("select post %d: %w", id, err)
	}
	return article, nil
}

func (r *mysqlRepository) Update(ctx context.Context, id int64, req UpsertRequest) (Article, error) {
	// Look up first: MySQL reports zero affected rows both for a missing id
	// and for an update with unchanged values, so RowsAffected alone cannot
	// distinguish "not found" from "no change".
	if _, err := r.GetByID(ctx, id); err != nil {
		return Article{}, err
	}

	_, err := r.db.ExecContext(ctx,
		"UPDATE posts SET title = ?, content = ?, category = ?, status = ? WHERE id = ?",
		req.Title, req.Content, req.Category, req.Status, id,
	)
	if err != nil {
		return Article{}, fmt.Errorf("update post %d: %w", id, err)
	}
	return r.GetByID(ctx, id)
}

func (r *mysqlRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM posts WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete post %d: %w", id, err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read affected rows: %w", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

// scanArticle reads one row from either *sql.Row or *sql.Rows.
func scanArticle(row interface{ Scan(dest ...any) error }) (Article, error) {
	var a Article
	err := row.Scan(&a.ID, &a.Title, &a.Content, &a.Category, &a.CreatedDate, &a.UpdatedDate, &a.Status)
	return a, err
}
