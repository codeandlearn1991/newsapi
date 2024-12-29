package news

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// Store is a wrapper around bun.DB.
type Store struct {
	db bun.IDB
}

// NewStore returns an instance of news store.
func NewStore(db bun.IDB) *Store {
	return &Store{
		db: db,
	}
}

// Create news record.
func (s Store) Create(ctx context.Context, news *Record) (*Record, error) {
	news.ID = uuid.New()
	if err := s.db.NewInsert().Model(news).Returning("*").Scan(ctx, news); err != nil {
		return nil, NewCustomError(err, http.StatusInternalServerError)
	}
	return news, nil
}

// FindByID finds a news record with the provided id.
func (s Store) FindByID(ctx context.Context, id uuid.UUID) (*Record, error) {
	var news Record
	if err := s.db.NewSelect().Model(&news).Where("id = ?", id).Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NewCustomError(err, http.StatusNotFound)
		}
		return nil, NewCustomError(err, http.StatusInternalServerError)
	}
	return &news, nil
}

// FindAll returns all news store in the database.
func (s Store) FindAll(ctx context.Context) ([]*Record, error) {
	var news []*Record
	if err := s.db.NewSelect().Model(&Record{}).Scan(ctx, &news); err != nil {
		return nil, NewCustomError(err, http.StatusInternalServerError)
	}
	return news, nil
}

// DeleteByID deletes a news by its ID.
func (s Store) DeleteByID(ctx context.Context, id uuid.UUID) (err error) {
	_, err = s.db.NewDelete().Model(&Record{}).Where("id = ?", id).Returning("NULL").Exec(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return NewCustomError(err, http.StatusInternalServerError)
	}
	return nil
}

// UpdateByID update news by it's ID.
func (s Store) UpdateByID(ctx context.Context, id uuid.UUID, news *Record) (err error) {
	r, err := s.db.NewUpdate().Model(news).Where("id = ?", id).Returning("NULL").Exec(ctx)
	if err != nil {
		return NewCustomError(err, http.StatusInternalServerError)
	}

	rowsAffected, err := r.RowsAffected()
	if err != nil {
		return NewCustomError(err, http.StatusInternalServerError)
	}
	if rowsAffected == 0 {
		return NewCustomError(err, http.StatusNotFound)
	}
	return nil
}
