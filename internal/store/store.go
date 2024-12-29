package store

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

// Store represents the in-memory news store.
type Store struct {
	l sync.Mutex
	n []*News
}

// New returns an instance of store.
func New() *Store {
	return &Store{
		l: sync.Mutex{},
		n: []*News{},
	}
}

// Create a news.
func (s *Store) Create(news *News) (*News, error) {
	s.l.Lock()
	defer s.l.Unlock()
	news.ID = uuid.New()
	s.n = append(s.n, news)
	return news, nil
}

// FindAll news.
func (s *Store) FindAll() ([]*News, error) {
	s.l.Lock()
	defer s.l.Unlock()
	return s.n, nil
}

// FindByID find a single news with its ID.
func (s *Store) FindByID(id uuid.UUID) (*News, error) {
	s.l.Lock()
	defer s.l.Unlock()
	for _, n := range s.n {
		if n.ID == id {
			return n, nil
		}
	}
	return nil, errors.New("news not found")
}

// DeleteByID delete a news by its ID.
func (s *Store) DeleteByID(id uuid.UUID) error {
	s.l.Lock()
	defer s.l.Unlock()

	idx := func(id uuid.UUID) int {
		for i, n := range s.n {
			if n.ID == id {
				return i
			}
		}
		return -1
	}(id)

	if idx == -1 {
		return errors.New("news not found")
	}
	s.n = append(s.n[:idx], s.n[idx+1:]...)
	return nil
}

// UpdateByID update a news by its ID.
func (s *Store) UpdateByID(news *News) error {
	s.l.Lock()
	defer s.l.Unlock()
	for idx, n := range s.n {
		if n.ID == news.ID {
			s.n[idx] = news
			return nil
		}
	}
	return errors.New("not found")
}
