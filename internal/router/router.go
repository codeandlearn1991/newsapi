package router

import (
	"net/http"

	"github.com/codeandlearn1991/newsapi/internal/handler"
)

// New creates a new router with all the handlers configured.
func New(ns handler.NewsStorer) *http.ServeMux {
	r := http.NewServeMux()

	// Create news route.
	r.HandleFunc("POST /news", handler.PostNews(ns))
	// Get all news.
	r.HandleFunc("GET /news", handler.GetAllNews(ns))
	// Get news by ID.
	r.HandleFunc("GET /news/{news_id}", handler.GetNewsByID(ns))
	// Update news by ID.
	r.HandleFunc("PUT /news/{news_id}", handler.UpdateNewsByID(ns))
	// Delete news by ID.
	r.HandleFunc("DELETE /news/{news_id}", handler.DeleteNewsByID(ns))

	return r
}
