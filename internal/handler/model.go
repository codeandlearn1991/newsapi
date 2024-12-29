package handler

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/codeandlearn1991/newsapi/internal/news"
	"github.com/google/uuid"
)

// NewsPostReqBody represents the post request body for the news.
type NewsPostReqBody struct {
	ID        uuid.UUID `json:"id"`
	Author    string    `json:"author"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	CreatedAt string    `json:"created_at"`
	Content   string    `json:"content"`
	Source    string    `json:"source"`
	Tags      []string  `json:"tags"`
}

// Validate the incoming request.
func (n *NewsPostReqBody) Validate() (record *news.Record, errs error) {
	if n.Author == "" {
		errs = errors.Join(errs, fmt.Errorf("author is empty: %s", n.Author))
	}
	if n.Title == "" {
		errs = errors.Join(errs, fmt.Errorf("title is empty: %s", n.Title))
	}
	if n.Content == "" {
		errs = errors.Join(errs, fmt.Errorf("content is empty: %s", n.Content))
	}
	if n.Summary == "" {
		errs = errors.Join(errs, fmt.Errorf("summary is empty: %s", n.Summary))
	}
	t, err := time.Parse(time.RFC3339, n.CreatedAt)
	if err != nil {
		errs = errors.Join(errs, err)
	}
	if n.Source == "" {
		errs = errors.Join(errs, fmt.Errorf("source is empty: %s", n.Source))
	}
	parsedURL, err := url.Parse(n.Source)
	if err != nil {
		errs = errors.Join(errs, err)
	}
	if len(n.Tags) == 0 {
		errs = errors.Join(errs, errors.New("tags cannot be empty"))
	}

	if errs != nil {
		return record, errs
	}
	return &news.Record{
		ID:        n.ID,
		Author:    n.Author,
		Title:     n.Title,
		Content:   n.Content,
		Summary:   n.Summary,
		CreatedAt: t,
		Source:    parsedURL.String(),
		Tags:      n.Tags,
	}, nil
}

// AllNewsResponse represents the all news response.
type AllNewsResponse struct {
	News []*news.Record `json:"news"`
}
