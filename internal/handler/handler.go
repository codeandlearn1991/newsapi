package handler

import (
	"encoding/json"
	"net/http"

	"github.com/codeandlearn1991/newsapi/internal/logger"
	"github.com/codeandlearn1991/newsapi/internal/store"
	"github.com/google/uuid"
)

// NewsStorer represents the news store opertions.
type NewsStorer interface {
	// Create news from post request body.
	Create(*store.News) (*store.News, error)
	// FindByID news by its ID.
	FindByID(uuid.UUID) (*store.News, error)
	// FindAll returns all news in the store.
	FindAll() ([]*store.News, error)
	// DeleteByID deletes a news item by its ID.
	DeleteByID(uuid.UUID) error
	// UpdateByID updates a news resource by its ID.
	UpdateByID(*store.News) error
}

// PostNews handler.
func PostNews(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromContext(r.Context())
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

		if _, err := ns.Create(&n); err != nil {
			log.Error("error creating news", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

// GetAllNews handler.
func GetAllNews(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromContext(r.Context())
		log.Info("request received")
		news, err := ns.FindAll()
		if err != nil {
			log.Error("failed to fetch all news", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		allNewsResponse := AllNewsResponse{News: news}
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
		log := logger.FromContext(r.Context())
		log.Info("request received")
		newsID := r.PathValue("news_id")
		newsUUID, err := uuid.Parse(newsID)
		if err != nil {
			log.Error("news id not a valid uuid", "newsId", newsID, "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		news, err := ns.FindByID(newsUUID)
		if err != nil {
			log.Error("news not found", "newsId", newsID)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(&news); err != nil {
			log.Error("failed to encode", "newsId", newsID, "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

// UpdateNewsByID handler.
func UpdateNewsByID(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromContext(r.Context())
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

		if err := ns.UpdateByID(&n); err != nil {
			log.Error("error updating news", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

// DeleteNewsByID handler.
func DeleteNewsByID(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromContext(r.Context())
		newsID := r.PathValue("news_id")
		newsUUID, err := uuid.Parse(newsID)
		if err != nil {
			log.Error("news id not a valid uuid", "newsId", newsID, "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := ns.DeleteByID(newsUUID); err != nil {
			log.Error("news not found", "newsId", newsID, "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
