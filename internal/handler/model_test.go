package handler_test

import (
	"testing"

	"github.com/codeandlearn1991/newsapi/internal/handler"
)

func TestNewsPostReqBody_Validate(t *testing.T) {
	testCases := []struct {
		name        string
		req         handler.NewsPostReqBody
		expectedErr bool
	}{
		{
			name:        "author empty",
			req:         handler.NewsPostReqBody{},
			expectedErr: true,
		},
		{
			name: "title empty",
			req: handler.NewsPostReqBody{
				Author: "test-author",
			},
			expectedErr: true,
		},
		{
			name: "summary empty",
			req: handler.NewsPostReqBody{
				Author: "test-author",
				Title:  "test-title",
			},
			expectedErr: true,
		},
		{
			name: "time invalid",
			req: handler.NewsPostReqBody{
				Author:    "test-author",
				Title:     "test-title",
				Summary:   "test-summary",
				CreatedAt: "invalid",
			},
			expectedErr: true,
		},
		{
			name: "source invalid",
			req: handler.NewsPostReqBody{
				Author:    "test-author",
				Title:     "test-title",
				Summary:   "test-summary",
				CreatedAt: "2024-04-07T05:13:27+00:00",
			},
			expectedErr: true,
		},
		{
			name: "tags empty",
			req: handler.NewsPostReqBody{
				Author:    "test-author",
				Title:     "test-title",
				Summary:   "test-summary",
				CreatedAt: "2024-04-07T05:13:27+00:00",
				Source:    "https://test-news.com",
			},
			expectedErr: true,
		},
		{
			name: "validate",
			req: handler.NewsPostReqBody{
				Author:    "test-author",
				Title:     "test-title",
				Summary:   "test-summary",
				CreatedAt: "2024-04-07T05:13:27+00:00",
				Source:    "https://test-news.com",
				Tags:      []string{"test-tag"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.req.Validate()

			if tc.expectedErr && err == nil {
				t.Fatal("expected error but got nil")
			}

			if !tc.expectedErr && err != nil {
				t.Fatalf("expected nil but got error: %v", err)
			}
		})
	}
}
