package lead

import (
	"context"

	"github.com/jackc/pgx/v5"
		lead "lead-automation/internal/lead/model"

)

type Repository struct {
	db *pgx.Conn
}

func NewRepository(db *pgx.Conn) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, l *lead.Lead) error {
	query := `
		INSERT INTO leads (
			name,
			email,
			company,
			message,
			budget,
			timeline
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		l.Name,
		l.Email,
		l.Company,
		l.Message,
		l.Budget,
		l.Timeline,
	).Scan(
		&l.ID,
		&l.CreatedAt,
	)
}