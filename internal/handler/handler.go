package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/codeandlearn1991/newsapi/internal/logger"
	"github.com/codeandlearn1991/newsapi/internal/news"
	"github.com/google/uuid"
)

//go:generate mockgen -source=handler.go -destination=mocks/handler.go -package=mockshandler

// NewsStorer represents the news store opertions.
type NewsStorer interface {
	// Create news from post request body.
	Create(context.Context, *news.Record) (*news.Record, error)
	// FindByID news by its ID.
	FindByID(context.Context, uuid.UUID) (*news.Record, error)
	// FindAll returns all news in the store.
	FindAll(context.Context) ([]*news.Record, error)
	// DeleteByID deletes a news item by its ID.
	DeleteByID(context.Context, uuid.UUID) error
	// UpdateByID updates a news resource by its ID.
	UpdateByID(context.Context, uuid.UUID, *news.Record) error
}

// PostNews handler.
func PostNews(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.FromContext(ctx)
		log.Info("request received")

		var newsRequestBody NewsPostReqBody
		if err := json.NewDecoder(r.Body).Decode(&newsRequestBody); err != nil {
			log.Error("failed to decode the request", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		n, err := newsRequestBody.Validate()
		if err != nil {
			log.Error("request validation failed", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			if _, wrErr := w.Write([]byte(err.Error())); wrErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		if _, err := ns.Create(ctx, n); err != nil {
			log.Error("error creating news", "error", err)
			var dbErr *news.CustomError
			if errors.As(err, &dbErr) {
				w.WriteHeader(dbErr.HTTPStatusCode())
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

// GetAllNews handler.
func GetAllNews(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.FromContext(ctx)
		log.Info("request received")
		n, err := ns.FindAll(ctx)
		if err != nil {
			log.Error("failed to fetch all news", "error", err)
			var dbErr *news.CustomError
			if errors.As(err, &dbErr) {
				w.WriteHeader(dbErr.HTTPStatusCode())
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		allNewsResponse := AllNewsResponse{News: n}
		if err := json.NewEncoder(w).Encode(allNewsResponse); err != nil {
			log.Error("failed to write response", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

// GetNewsByID handler.
func GetNewsByID(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.FromContext(ctx)
		log.Info("request received")
		newsID := r.PathValue("news_id")
		newsUUID, err := uuid.Parse(newsID)
		if err != nil {
			log.Error("news id not a valid uuid", "newsId", newsID, "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		n, err := ns.FindByID(ctx, newsUUID)
		if err != nil {
			log.Error("news not found", "newsId", newsID)
			var dbErr *news.CustomError
			if errors.As(err, &dbErr) {
				w.WriteHeader(dbErr.HTTPStatusCode())
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(&n); err != nil {
			log.Error("failed to encode", "newsId", newsID, "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

// UpdateNewsByID handler.
func UpdateNewsByID(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.FromContext(ctx)
		log.Info("request received")

		var newsRequestBody NewsPostReqBody
		if err := json.NewDecoder(r.Body).Decode(&newsRequestBody); err != nil {
			log.Error("failed to decode the request", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		n, err := newsRequestBody.Validate()
		if err != nil {
			log.Error("request validation failed", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			if _, wrErr := w.Write([]byte(err.Error())); wrErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		if err := ns.UpdateByID(ctx, n.ID, n); err != nil {
			log.Error("error updating news", "error", err)
			var dbErr *news.CustomError
			if errors.As(err, &dbErr) {
				w.WriteHeader(dbErr.HTTPStatusCode())
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

// DeleteNewsByID handler.
func DeleteNewsByID(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.FromContext(ctx)
		newsID := r.PathValue("news_id")
		newsUUID, err := uuid.Parse(newsID)
		if err != nil {
			log.Error("news id not a valid uuid", "newsId", newsID, "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := ns.DeleteByID(ctx, newsUUID); err != nil {
			log.Error("news not found", "newsId", newsID, "error", err)
			var dbErr *news.CustomError
			if errors.As(err, &dbErr) {
				w.WriteHeader(dbErr.HTTPStatusCode())
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
