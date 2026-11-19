package repository

import (
	"context"
	"errors"
	"fmt"

	"phone-book/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("contact not found")

type ContactRepository interface {
	GetAll(ctx context.Context) ([]model.Contact, error)
	GetByID(ctx context.Context, id int) (*model.Contact, error)
	Create(ctx context.Context, contact *model.Contact) error
	Update(ctx context.Context, contact *model.Contact) error
	Delete(ctx context.Context, id int) error
}

type postgresContactRepo struct {
	pool *pgxpool.Pool
}

func NewContactRepository(pool *pgxpool.Pool) ContactRepository {
	return &postgresContactRepo{pool: pool}
}

func (r *postgresContactRepo) GetAll(ctx context.Context) ([]model.Contact, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT id, username, email, mobile, home FROM contacts ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var contacts []model.Contact
	for rows.Next() {
		var c model.Contact
		if err := rows.Scan(&c.ID, &c.Username, &c.Email, &c.Mobile, &c.Home); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		contacts = append(contacts, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}
	if contacts == nil {
		contacts = []model.Contact{}
	}
	return contacts, nil
}

func (r *postgresContactRepo) GetByID(ctx context.Context, id int) (*model.Contact, error) {
	var c model.Contact
	err := r.pool.QueryRow(ctx,
		"SELECT id, username, email, mobile, home FROM contacts WHERE id = $1", id).
		Scan(&c.ID, &c.Username, &c.Email, &c.Mobile, &c.Home)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query row: %w", err)
	}
	return &c, nil
}

func (r *postgresContactRepo) Create(ctx context.Context, contact *model.Contact) error {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO contacts (username, email, mobile, home)
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		contact.Username, contact.Email, contact.Mobile, contact.Home,
	).Scan(&contact.ID)
	if err != nil {
		return fmt.Errorf("insert: %w", err)
	}
	return nil
}

func (r *postgresContactRepo) Update(ctx context.Context, contact *model.Contact) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE contacts SET username=$1, email=$2, mobile=$3, home=$4 WHERE id=$5`,
		contact.Username, contact.Email, contact.Mobile, contact.Home, contact.ID,
	)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *postgresContactRepo) Delete(ctx context.Context, id int) error {
	tag, err := r.pool.Exec(ctx, "DELETE FROM contacts WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}