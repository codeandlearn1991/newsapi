package handler_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/codeandlearn1991/newsapi/internal/handler"
	"github.com/codeandlearn1991/newsapi/internal/news"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewsPostReqBody_Validate(t *testing.T) {
	type expectations struct {
		err  string
		news *news.Record
	}
	testCases := []struct {
		name         string
		req          handler.NewsPostReqBody
		expectations expectations
	}{
		{
			name: "author empty",
			req:  handler.NewsPostReqBody{},
			expectations: expectations{
				err: "author is empty",
			},
		},
		{
			name: "title empty",
			req: handler.NewsPostReqBody{
				Author: "test-author",
			},
			expectations: expectations{
				err: "title is empty",
			},
		},
		{
			name: "content empty",
			req: handler.NewsPostReqBody{
				Author: "test-author",
				Title:  "test-title",
			},
			expectations: expectations{
				err: "content is empty",
			},
		},
		{
			name: "summary empty",
			req: handler.NewsPostReqBody{
				Author:  "test-author",
				Title:   "test-title",
				Content: "test-content",
			},
			expectations: expectations{
				err: "summary is empty",
			},
		},
		{
			name: "time invalid",
			req: handler.NewsPostReqBody{
				Author:    "test-author",
				Title:     "test-title",
				Summary:   "test-summary",
				Content:   "test-content",
				CreatedAt: "invalid",
			},
			expectations: expectations{
				err: `parsing time "invalid"`,
			},
		},
		{
			name: "source empty",
			req: handler.NewsPostReqBody{
				Author:    "test-author",
				Title:     "test-title",
				Summary:   "test-summary",
				CreatedAt: "2024-04-07T05:13:27+00:00",
				Content:   "test-content",
			},
			expectations: expectations{
				err: "source is empty",
			},
		},
		{
			name: "souce invalid url",
			req: handler.NewsPostReqBody{
				Author:    "test-author",
				Title:     "test-title",
				Summary:   "test-summary",
				CreatedAt: "2024-04-07T05:13:27+00:00",
				Source:    "https://xyz:abc",
				Content:   "test-content",
			},
			expectations: expectations{
				err: "invalid port",
			},
		},
		{
			name: "tags empty",
			req: handler.NewsPostReqBody{
				Author:    "test-author",
				Title:     "test-title",
				Summary:   "test-summary",
				CreatedAt: "2024-04-07T05:13:27+00:00",
				Source:    "https://test-news.com",
				Content:   "test-content",
			},
			expectations: expectations{
				err: "tags cannot be empty",
			},
		},
		{
			name: "validate",
			req: handler.NewsPostReqBody{
				Author:    "test-author",
				Title:     "test-title",
				Content:   "test-content",
				Summary:   "test-summary",
				CreatedAt: "2024-04-07T05:13:27+00:00",
				Source:    "https://test-news.com",
				Tags:      []string{"test-tag"},
			},
			expectations: expectations{
				news: &news.Record{
					Author:  "test-author",
					Title:   "test-title",
					Content: "test-content",
					Summary: "test-summary",
					Tags:    []string{"test-tag"},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			n, err := tc.req.Validate()

			// Assert
			if tc.expectations.err != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectations.err)
			} else {
				assert.NoError(t, err)

				parsedTime, parseErr := time.Parse(time.RFC3339, tc.req.CreatedAt)
				require.NoError(t, parseErr)
				tc.expectations.news.CreatedAt = parsedTime

				parsedSource, err := url.Parse(tc.req.Source)
				require.NoError(t, err)
				tc.expectations.news.Source = parsedSource.String()

				assert.Equal(t, tc.expectations.news, n)
			}
		})
	}
}
