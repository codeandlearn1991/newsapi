package handler_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/codeandlearn1991/newsapi/internal/handler"
	mockshandler "github.com/codeandlearn1991/newsapi/internal/handler/mocks"
	"github.com/codeandlearn1991/newsapi/internal/news"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_PostNews(t *testing.T) {
	testCases := []struct {
		name           string
		body           io.Reader
		setup          func(tb testing.TB) *mockshandler.MockNewsStorer
		expectedStatus int
	}{
		{
			name: "invalid request body json",
			body: strings.NewReader(`{`),
			setup: func(tb testing.TB) *mockshandler.MockNewsStorer {
				tb.Helper()
				return mockshandler.NewMockNewsStorer(gomock.NewController(t))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid request body",
			body: strings.NewReader(`
			{ 
			"id" : "3b082d9d-1dc7-4d1f-907e-50d449a03d45", 
			"author": "code learn", 
			"title": "first news", 
			"summary": "first news post", 
			"created_at": "2024-04-07T05:13:27+00:00", 
			"source": "https://example.com"
			}`),
			setup: func(tb testing.TB) *mockshandler.MockNewsStorer {
				tb.Helper()
				return mockshandler.NewMockNewsStorer(gomock.NewController(t))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "db error",
			body: strings.NewReader(`
			{ 
			"id" : "3b082d9d-1dc7-4d1f-907e-50d449a03d45", 
			"author": "code learn", 
			"content": "news content",
			"title": "first news", 
			"summary": "first news post", 
			"created_at": "2024-04-07T05:13:27+00:00", 
			"source": "https://example.com",
			"tags": ["politics"]
			}`),
			setup: func(tb testing.TB) *mockshandler.MockNewsStorer {
				tb.Helper()
				ms := mockshandler.NewMockNewsStorer(gomock.NewController(t))
				ms.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, errors.New("db error"))
				return ms
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "db custom error",
			body: strings.NewReader(`
			{ 
			"id" : "3b082d9d-1dc7-4d1f-907e-50d449a03d45", 
			"author": "code learn", 
			"content": "news content",
			"title": "first news", 
			"summary": "first news post", 
			"created_at": "2024-04-07T05:13:27+00:00", 
			"source": "https://example.com",
			"tags": ["politics"]
			}`),
			setup: func(tb testing.TB) *mockshandler.MockNewsStorer {
				tb.Helper()
				ms := mockshandler.NewMockNewsStorer(gomock.NewController(t))
				ms.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, news.NewCustomError(errors.New("some error"), http.StatusBadRequest))
				return ms
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "success",
			body: strings.NewReader(`
			{ 
			"id" : "3b082d9d-1dc7-4d1f-907e-50d449a03d45", 
			"author": "code learn", 
						"content": "news content",

			"title": "first news", 
			"summary": "first news post", 
			"created_at": "2024-04-07T05:13:27+00:00", 
			"source": "https://example.com",
			"tags": ["politics"]
			}`),
			setup: func(tb testing.TB) *mockshandler.MockNewsStorer {
				tb.Helper()
				ms := mockshandler.NewMockNewsStorer(gomock.NewController(t))
				ms.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, nil)
				return ms
			},
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", tc.body)

			// Act
			handler.PostNews(tc.setup(t))(w, r)

			// Assert
			assert.Equal(t, tc.expectedStatus, w.Result().StatusCode)
		})
	}
}

func Test_GetAllNews(t *testing.T) {
	testCases := []struct {
		name           string
		setup          func(tb testing.TB) *mockshandler.MockNewsStorer
		expectedStatus int
	}{
		{
			name: "db error",
			setup: func(tb testing.TB) *mockshandler.MockNewsStorer {
				tb.Helper()
				ms := mockshandler.NewMockNewsStorer(gomock.NewController(t))
				ms.EXPECT().FindAll(gomock.Any()).Return(nil, errors.New("db error"))
				return ms
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "db customer error",
			setup: func(tb testing.TB) *mockshandler.MockNewsStorer {
				tb.Helper()
				ms := mockshandler.NewMockNewsStorer(gomock.NewController(t))
				ms.EXPECT().FindAll(gomock.Any()).Return(nil, news.NewCustomError(errors.New("some error"), http.StatusBadRequest))
				return ms
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "success",
			setup: func(tb testing.TB) *mockshandler.MockNewsStorer {
				tb.Helper()
				ms := mockshandler.NewMockNewsStorer(gomock.NewController(t))
				ms.EXPECT().FindAll(gomock.Any()).Return(nil, nil)
				return ms
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", http.NoBody)

			// Act
			handler.GetAllNews(tc.setup(t))(w, r)

			// Assert
			assert.Equal(t, tc.expectedStatus, w.Result().StatusCode)
		})
	}
}

func Test_GetNewsByID(t *testing.T) {
	testCases := []struct {
		name           string
		setup          func(tb testing.TB) *mockshandler.MockNewsStorer
		newsID         string
		expectedStatus int
	}{
		{
			name: "invalid news id",
			setup: func(tb testing.TB) *mockshandler.MockNewsStorer {
				tb.Helper()
				return mockshandler.NewMockNewsStorer(gomock.NewController(t))
			},
			newsID:         "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "db error",
			setup: func(tb testing.TB) *mockshandler.MockNewsStorer {
				tb.Helper()
				ms := mockshandler.NewMockNewsStorer(gomock.NewController(t))
				ms.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(nil, errors.New("db error"))
				return ms
			},
			newsID:         uuid.NewString(),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "db error",
			setup: func(tb testing.TB) *mockshandler.MockNewsStorer {
				tb.Helper()
				ms := mockshandler.NewMockNewsStorer(gomock.NewController(t))
				ms.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(nil, news.NewCustomError(errors.New("some error"), http.StatusBadRequest))
				return ms
			},
			newsID:         uuid.NewString(),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "success",
			setup: func(tb testing.TB) *mockshandler.MockNewsStorer {
				tb.Helper()
				ms := mockshandler.NewMockNewsStorer(gomock.NewController(t))
				ms.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(nil, nil)
				return ms
			},
			newsID:         uuid.NewString(),
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", http.NoBody)
			r.SetPathValue("news_id", tc.newsID)

			// Act
			handler.GetNewsByID(tc.setup(t))(w, r)

			// Assert
			assert.Equal(t, tc.expectedStatus, w.Result().StatusCode)
		})
	}
}

func Test_UpdateNewsByID(t *testing.T) {
	testCases := []struct {
		name           string
		body           io.Reader
		setup          func(tb testing.TB) *mockshandler.MockNewsStorer
		expectedStatus int
	}{
		{
			name: "invalid request json body",
			body: strings.NewReader(`{`),
			setup: func(tb testing.TB) *mockshandler.MockNewsStorer {
				tb.Helper()
				return mockshandler.NewMockNewsStorer(gomock.NewController(t))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid request body",
			body: strings.NewReader(`
			{ 
			"id" : "3b082d9d-1dc7-4d1f-907e-50d449a03d45", 
			"author": "code learn", 
			"title": "first news", 
			"summary": "first news post", 
			"created_at": "2024-04-07T05:13:27+00:00", 
			"source": "https://example.com"
			}`),
			setup: func(tb testing.TB) *mockshandler.MockNewsStorer {
				tb.Helper()
				return mockshandler.NewMockNewsStorer(gomock.NewController(t))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "db error",
			body: strings.NewReader(`
			{ 
			"id" : "3b082d9d-1dc7-4d1f-907e-50d449a03d45", 
			"author": "code learn", 
			"content": "news content",
			"title": "first news", 
			"summary": "first news post", 
			"created_at": "2024-04-07T05:13:27+00:00", 
			"source": "https://example.com",
			"tags": ["politics"]
			}`),
			setup: func(tb testing.TB) *mockshandler.MockNewsStorer {
				tb.Helper()
				ms := mockshandler.NewMockNewsStorer(gomock.NewController(t))
				ms.EXPECT().UpdateByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("db error"))
				return ms
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "db custom error",
			body: strings.NewReader(`
			{ 
			"id" : "3b082d9d-1dc7-4d1f-907e-50d449a03d45", 
			"author": "code learn", 
			"content": "news content",
			"title": "first news", 
			"summary": "first news post", 
			"created_at": "2024-04-07T05:13:27+00:00", 
			"source": "https://example.com",
			"tags": ["politics"]
			}`),
			setup: func(tb testing.TB) *mockshandler.MockNewsStorer {
				tb.Helper()
				ms := mockshandler.NewMockNewsStorer(gomock.NewController(t))
				ms.EXPECT().UpdateByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(news.NewCustomError(errors.New("some error"), http.StatusBadRequest))
				return ms
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "success",
			body: strings.NewReader(`
			{ 
			"id" : "3b082d9d-1dc7-4d1f-907e-50d449a03d45", 
			"author": "code learn", 
			"title": "first news", 
			"content": "news content",
			"summary": "first news post", 
			"created_at": "2024-04-07T05:13:27+00:00", 
			"source": "https://example.com",
			"tags": ["politics"]
			}`),
			setup: func(tb testing.TB) *mockshandler.MockNewsStorer {
				tb.Helper()
				ms := mockshandler.NewMockNewsStorer(gomock.NewController(t))
				ms.EXPECT().UpdateByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return ms
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", tc.body)

			// Act
			handler.UpdateNewsByID(tc.setup(t))(w, r)

			// Assert
			assert.Equal(t, tc.expectedStatus, w.Result().StatusCode)
		})
	}
}

func Test_DeleteNewsByID(t *testing.T) {
	testCases := []struct {
		name           string
		setup          func(testing.TB) *mockshandler.MockNewsStorer
		newsID         string
		expectedStatus int
	}{
		{
			name: "invalid news id",
			setup: func(tb testing.TB) *mockshandler.MockNewsStorer {
				tb.Helper()
				return mockshandler.NewMockNewsStorer(gomock.NewController(t))
			},
			newsID:         "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "db error",
			setup: func(tb testing.TB) *mockshandler.MockNewsStorer {
				tb.Helper()
				ms := mockshandler.NewMockNewsStorer(gomock.NewController(t))
				ms.EXPECT().DeleteByID(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
				return ms
			},
			newsID:         uuid.NewString(),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "db custom error",
			setup: func(tb testing.TB) *mockshandler.MockNewsStorer {
				tb.Helper()
				ms := mockshandler.NewMockNewsStorer(gomock.NewController(t))
				ms.EXPECT().DeleteByID(gomock.Any(), gomock.Any()).Return(news.NewCustomError(errors.New("some error"), http.StatusBadRequest))
				return ms
			},
			newsID:         uuid.NewString(),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "success",
			setup: func(tb testing.TB) *mockshandler.MockNewsStorer {
				tb.Helper()
				ms := mockshandler.NewMockNewsStorer(gomock.NewController(t))
				ms.EXPECT().DeleteByID(gomock.Any(), gomock.Any()).Return(nil)
				return ms
			},
			newsID:         uuid.NewString(),
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", http.NoBody)
			r.SetPathValue("news_id", tc.newsID)

			// Act
			handler.DeleteNewsByID(tc.setup(t))(w, r)

			// Assert
			assert.Equal(t, tc.expectedStatus, w.Result().StatusCode)
		})
	}
}
