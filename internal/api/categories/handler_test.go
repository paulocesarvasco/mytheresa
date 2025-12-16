package categoriesapi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/internal/categories"
	errorsapi "github.com/mytheresa/go-hiring-challenge/internal/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCategories(t *testing.T) {
	tests := []struct {
		name           string
		fakeCategories categories.CategoryPage
		fakeError      error
		queryParams    map[string]string
		expectedStatus int
		expectedCT     string
		expectedBody   any
	}{
		{
			name: "fetch categories list ok",
			fakeCategories: categories.CategoryPage{
				Total: 2,
				Categories: []categories.CategoryView{
					{Code: "FOO", Name: "foo"},
					{Code: "BAR", Name: "bar"},
				},
			},
			queryParams: map[string]string{
				"limit":  "2",
				"offset": "0",
			},
			expectedStatus: http.StatusOK,
			expectedCT:     "application/json",
			expectedBody: categories.CategoryPage{
				Total: 2,
				Categories: []categories.CategoryView{
					{Code: "FOO", Name: "foo"},
					{Code: "BAR", Name: "bar"},
				},
			},
		},
		{
			name: "invalid limit parameter",
			queryParams: map[string]string{
				"limit": "0x01",
			},
			expectedStatus: http.StatusBadRequest,
			expectedCT:     "application/json",
			expectedBody:   map[string]any{"error": errorsapi.ErrInvalidLimitParam.Error()},
		},
		{
			name: "max limit parameter",
			fakeCategories: categories.CategoryPage{
				Total: 2,
				Categories: []categories.CategoryView{
					{Code: "FOO", Name: "foo"},
					{Code: "BAR", Name: "bar"},
				},
			},
			queryParams: map[string]string{
				"limit": "200",
			},
			expectedStatus: http.StatusOK,
			expectedCT:     "application/json",
			expectedBody: categories.CategoryPage{
				Total: 2,
				Categories: []categories.CategoryView{
					{Code: "FOO", Name: "foo"},
					{Code: "BAR", Name: "bar"},
				},
			},
		},
		{
			name: "negative offset",
			queryParams: map[string]string{
				"offset": "-1",
			},
			expectedStatus: http.StatusBadRequest,
			expectedCT:     "application/json",
			expectedBody:   map[string]any{"error": errorsapi.ErrInvalidOffsetParam.Error()},
		},
		{
			name: "negative offset",
			queryParams: map[string]string{
				"offset": "-1",
			},
			expectedStatus: http.StatusBadRequest,
			expectedCT:     "application/json",
			expectedBody:   map[string]any{"error": errorsapi.ErrInvalidOffsetParam.Error()},
		},
		{
			name:           "catalog service error",
			fakeCategories: categories.CategoryPage{},
			fakeError:      errorsapi.ErrRepositoryFetchCategories,
			expectedStatus: http.StatusInternalServerError,
			expectedCT:     "application/json",
			expectedBody:   map[string]any{"error": errorsapi.ErrRepositoryFetchCategories.Error()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewFakeService()
			s.SetListCategoriesResponse(tt.fakeCategories, tt.fakeError)
			h := New(s)
			r := Routes(h)

			ts := httptest.NewServer(r)
			defer ts.Close()

			u, err := url.Parse(ts.URL)
			require.NoError(t, err)

			q := u.Query()
			for k, v := range tt.queryParams {
				q.Set(k, v)
			}
			u.RawQuery = q.Encode()

			res, err := ts.Client().Get(u.String())
			require.NoError(t, err)
			defer res.Body.Close()

			rawBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, res.StatusCode)
			assert.Contains(t, res.Header.Get("Content-Type"), tt.expectedCT)

			rawExpectedBody, err := json.Marshal(tt.expectedBody)
			require.NoError(t, err)

			assert.JSONEq(t, string(rawExpectedBody), string(rawBody))
		})
	}
}

func TestCreateCategory(t *testing.T) {
	tests := []struct {
		name           string
		fakeCategory   categories.CategoryView
		fakeError      error
		requestBody    map[string]any
		contentType    string
		expectedStatus int
		expectedCT     string
		expectedBody   any
	}{
		{
			name:         "create category success",
			fakeCategory: categories.CategoryView{Code: "foo", Name: "bar"},
			requestBody: map[string]any{
				"code": "foo",
				"name": "bar",
			},
			contentType:    "application/json",
			expectedStatus: http.StatusCreated,
			expectedCT:     "application/json",
			expectedBody:   categories.CategoryView{Code: "foo", Name: "bar"},
		},
		{
			name:           "invalid content type",
			contentType:    "text/plain",
			expectedStatus: http.StatusUnsupportedMediaType,
			expectedCT:     "application/json",
			expectedBody:   map[string]any{"error": errorsapi.ErrInvalidContentType.Error()},
		},
		{
			name:           "invalid json format",
			requestBody:    map[string]any{"key": `invalid"json`},
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
			expectedCT:     "application/json",
			expectedBody:   map[string]any{"error": errorsapi.ErrInvalidJSONBody.Error()},
		},
		{
			name: "json with extra fields",
			requestBody: map[string]any{
				"code": "foo",
				"name": "bar",
				"id":   1,
			},
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
			expectedCT:     "application/json",
			expectedBody:   map[string]any{"error": errorsapi.ErrInvalidJSONBody.Error()},
		},
		{
			name: "json without required values",
			requestBody: map[string]any{
				"code": "foo",
			},
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
			expectedCT:     "application/json",
			expectedBody:   map[string]any{"error": errorsapi.ErrInvalidRequestSchema.Error()},
		},
		{
			name:      "registers conflict",
			fakeError: errorsapi.ErrRepositoryCategoryAlreadyExists,
			requestBody: map[string]any{
				"code": "foo",
				"name": "bar",
			},
			contentType:    "application/json",
			expectedStatus: http.StatusConflict,
			expectedCT:     "application/json",
			expectedBody:   map[string]any{"error": errorsapi.ErrRepositoryCategoryAlreadyExists.Error()},
		},
		{
			name:      "failed to insert new register",
			fakeError: errorsapi.ErrRepositoryCreateCategory,
			requestBody: map[string]any{
				"code": "foo",
				"name": "bar",
			},
			contentType:    "application/json",
			expectedStatus: http.StatusInternalServerError,
			expectedCT:     "application/json",
			expectedBody:   map[string]any{"error": errorsapi.ErrRepositoryCreateCategory.Error()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewFakeService()
			s.SetCreateCategoryResponse(tt.fakeCategory, tt.fakeError)
			h := New(s)
			r := Routes(h)

			ts := httptest.NewServer(r)
			defer ts.Close()

			var body io.Reader

			if tt.requestBody != nil {
				rawBody, err := json.Marshal(tt.requestBody)
				require.NoError(t, err)
				body = bytes.NewReader(rawBody)
			}

			res, err := ts.Client().Post(ts.URL, tt.contentType, body)
			require.NoError(t, err)
			defer res.Body.Close()

			rawBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, res.StatusCode)
			assert.Contains(t, res.Header.Get("Content-Type"), tt.expectedCT)

			rawExpectedBody, err := json.Marshal(tt.expectedBody)
			require.NoError(t, err)

			assert.JSONEq(t, string(rawExpectedBody), string(rawBody))
		})
	}
}
