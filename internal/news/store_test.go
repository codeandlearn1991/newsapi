package news_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/codeandlearn1991/newsapi/internal/news"
	"github.com/codeandlearn1991/newsapi/internal/postgres"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	pgtc "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/uptrace/bun"
)

var db *bun.DB

func TestMain(m *testing.M) {
	ctx := context.Background()
	pdb, cf, err := createTestDB(ctx)
	if err != nil {
		panic(err)
	}

	db = pdb
	code := m.Run()

	if err := cf(ctx); err != nil {
		panic(err)
	}

	os.Exit(code)
}

func TestStore_Create(t *testing.T) {
	testCases := []struct {
		name               string
		news               *news.Record
		expectedErr        string
		expectedStatusCode int
	}{
		{
			name: "missing author",
			news: &news.Record{
				Title:   "test-title",
				Summary: "test-summary",
				Content: "test-content",
				Source:  "https://example.com",
				Tags:    []string{"tag1", "tag2"},
			},
			expectedErr:        "not-null",
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "success",
			news: &news.Record{
				Author:  "test-author",
				Title:   "test-title",
				Summary: "test-summary",
				Content: "test-content",
				Source:  "https://www.example.com",
				Tags:    []string{"tag1", "tag2"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := news.NewStore(db)
			createdNews, err := s.Create(context.Background(), tc.news)

			if tc.expectedErr != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectedErr)
				var storeErr *news.CustomError
				assert.ErrorAs(t, err, &storeErr)
				assert.Equal(t, tc.expectedStatusCode, storeErr.HTTPStatusCode())
			} else {
				assert.NoError(t, err)
				assertOnNews(t, tc.news, createdNews)
				err = s.DeleteByID(context.Background(), createdNews.ID)
				assert.NoError(t, err)
			}
		})
	}
}

func TestStore_FindByID(t *testing.T) {
	testCases := []struct {
		name               string
		id                 uuid.UUID
		expectedNews       *news.Record
		expectedStatusCode int
		expectedErr        string
	}{
		{
			name: "found",
			id:   uuid.MustParse("17628bea-9d11-47f9-986e-16703a87e451"),
			expectedNews: &news.Record{
				Author:  "Batman",
				Title:   "Breaking News",
				Summary: "A brief summary of the news",
				Content: "Full content of the news article",
				Source:  "https://www.example.com",
				Tags:    []string{"tag1", "tag2"},
			},
		},
		{
			name:               "not found",
			id:                 uuid.New(),
			expectedErr:        "no rows",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "soft deleted",
			id:                 uuid.MustParse("f710bc79-9ad3-4e0f-8dab-e43d94b42fbb"),
			expectedErr:        "no rows",
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := news.NewStore(db)

			n, err := s.FindByID(context.Background(), tc.id)

			if tc.expectedErr != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectedErr)
				var storeErr *news.CustomError
				assert.ErrorAs(t, err, &storeErr)
				assert.Equal(t, tc.expectedStatusCode, storeErr.HTTPStatusCode())
			} else {
				assert.NoError(t, err)
				assertOnNews(t, tc.expectedNews, n)
			}
		})
	}
}

func TestStore_FindAll(t *testing.T) {
	testCases := []struct {
		name         string
		expectedNews []*news.Record
	}{
		{
			name: "found all",
			expectedNews: []*news.Record{
				{
					Author:  "Batman",
					Title:   "Breaking News",
					Summary: "A brief summary of the news",
					Content: "Full content of the news article",
					Source:  "https://www.example.com",
					Tags:    []string{"tag1", "tag2"},
				},
				{
					Author:  "Superman",
					Title:   "Breaking News",
					Summary: "A brief summary of the news",
					Content: "Full content of the news article",
					Source:  "https://www.example.com",
					Tags:    []string{"tag1", "tag2"},
				},
			},
		},
	}

	for _, tc := range testCases {
		s := news.NewStore(db)
		allNews, err := s.FindAll(context.Background())

		assert.NoError(t, err)
		assert.Len(t, allNews, len(tc.expectedNews))
		for idx, n := range allNews {
			assertOnNews(t, tc.expectedNews[idx], n)
		}
	}
}

func TestStore_DeleteByID(t *testing.T) {
	testCases := []struct {
		name string
		id   uuid.UUID
	}{
		{
			name: "deleted",
			id:   uuid.MustParse("17628bea-9d11-47f9-986e-16703a87e451"),
		},
		{
			name: "not found",
			id:   uuid.MustParse("f710bc79-9ad3-4e0f-8dab-e43d94b42fbb"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := news.NewStore(db)

			err := s.DeleteByID(context.Background(), tc.id)
			assert.NoError(t, err)
		})
	}
}

func TestStore_UpdatedByID(t *testing.T) {
	testCases := []struct {
		name           string
		news           *news.Record
		expectedStatus int
	}{
		{
			name: "updated",
			news: &news.Record{
				ID:        uuid.MustParse("bde0c593-0df6-4eba-9326-3f00be67aade"),
				Author:    "Wolverine",
				Title:     "Breaking News",
				Summary:   "A brief summary of the news",
				Content:   "Full content of the news article",
				Source:    "https://www.example.com",
				Tags:      []string{"tag1", "tag2"},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		{
			name: "not found",
			news: &news.Record{
				ID: uuid.MustParse("6a3483c7-e28e-442e-b603-b06ff60eeeb4"),
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := news.NewStore(db)

			err := s.UpdateByID(context.Background(), tc.news.ID, tc.news)

			if tc.expectedStatus != 0 {
				assert.Error(t, err)
				var storeErr *news.CustomError
				assert.ErrorAs(t, err, &storeErr)
				assert.Equal(t, tc.expectedStatus, storeErr.HTTPStatusCode())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func assertOnNews(tb testing.TB, expected, got *news.Record) {
	tb.Helper()
	assert.Equal(tb, expected.Author, got.Author)
	assert.Equal(tb, expected.Title, got.Title)
	assert.Equal(tb, expected.Content, got.Content)
	assert.Equal(tb, expected.Summary, got.Summary)
	assert.Equal(tb, expected.Source, got.Source)
	assert.Equal(tb, expected.Tags, got.Tags)
	assert.NotEqual(tb, time.Time{}, got.CreatedAt)
	assert.NotEqual(tb, time.Time{}, got.UpdatedAt)
	assert.Equal(tb, time.Time{}, got.DeletedAt)
}

func createTestContainer(ctx context.Context) (ctr *pgtc.PostgresContainer, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return ctr, fmt.Errorf("working dir: %w", err)
	}

	sqlScripts := wd + "/testdata/sql/store.sql"

	ctr, err = pgtc.Run(
		ctx,
		"postgres:16-alpine",
		pgtc.WithInitScripts(sqlScripts),
		pgtc.WithDatabase("postgres"),
		pgtc.WithUsername("postgres"),
		pgtc.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		return ctr, fmt.Errorf("run container: %w", err)
	}
	return ctr, nil
}

type DBCleanupFunc func(ctx context.Context) error

func createTestDB(ctx context.Context) (*bun.DB, DBCleanupFunc, error) {
	ctr, err := createTestContainer(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("create test container: %w", err)
	}

	p, err := ctr.MappedPort(ctx, nat.Port("5432/tcp"))
	if err != nil {
		return nil, nil, fmt.Errorf("mapped port: %w", err)
	}

	db, err := postgres.NewDB(&postgres.Config{
		Host:     "localhost",
		Debug:    true,
		DBName:   "postgres",
		User:     "postgres",
		Password: "postgres",
		Port:     p.Port(),
		SSLMode:  "disable",
	})
	if err != nil {
		return nil, nil, fmt.Errorf("new db: %w", err)
	}

	cf := func(ctx context.Context) error {
		if err := db.Close(); err != nil {
			return fmt.Errorf("db close: %w", err)
		}
		if err := ctr.Terminate(ctx); err != nil {
			return fmt.Errorf("container terminate: %w", err)
		}
		return nil
	}

	return db, cf, nil
}
