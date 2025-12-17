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
		fakeCategories []categories.Category
		fakeTotal      int64
		fakeError      error
		queryParams    map[string]string
		expectedStatus int
		expectedCT     string
		expectedBody   any
	}{
		{
			name: "fetch categories list ok",
			fakeCategories: []categories.Category{
				{Code: "FOO", Name: "foo"},
				{Code: "BAR", Name: "bar"},
			},
			fakeTotal: 2,
			queryParams: map[string]string{
				"limit":  "2",
				"offset": "0",
			},
			expectedStatus: http.StatusOK,
			expectedCT:     "application/json",
			expectedBody: CategoryPage{
				Total: 2,
				Categories: []CategoryView{
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
			fakeCategories: []categories.Category{
				{Code: "FOO", Name: "foo"},
				{Code: "BAR", Name: "bar"},
			},
			fakeTotal: 2,
			queryParams: map[string]string{
				"limit": "200",
			},
			expectedStatus: http.StatusOK,
			expectedCT:     "application/json",
			expectedBody: CategoryPage{
				Total: 2,
				Categories: []CategoryView{
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
			fakeError:      errorsapi.ErrRepositoryFetchCategories,
			expectedStatus: http.StatusInternalServerError,
			expectedCT:     "application/json",
			expectedBody:   map[string]any{"error": errorsapi.ErrRepositoryFetchCategories.Error()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewFakeService()
			s.SetListCategoriesResponse(tt.fakeCategories, tt.fakeTotal, tt.fakeError)
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
		fakeCategory   categories.Category
		fakeError      error
		requestBody    any
		contentType    string
		expectedStatus int
		expectedCT     string
		expectedBody   any
	}{
		{
			name: "creates single category successfully",
			requestBody: map[string]any{
				"code": "foo",
				"name": "bar",
			},
			contentType:    "application/json",
			expectedStatus: http.StatusCreated,
			expectedCT:     "application/json",
			expectedBody: []map[string]any{
				{"code": "foo", "name": "bar"},
			},
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
			requestBody:    `{"key": "invalid"json}`,
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
		{
			name: "successfully creates a batch of categories",
			requestBody: []map[string]any{
				{"code": "FOO", "name": "foo"},
				{"code": "BAR", "name": "bar"},
			},
			contentType:    "application/json",
			expectedStatus: http.StatusCreated,
			expectedCT:     "application/json",
			expectedBody: []map[string]any{
				{"code": "FOO", "name": "foo"},
				{"code": "BAR", "name": "bar"},
			},
		},
		{
			name:           "request body is not object or list",
			requestBody:    `[{"code":"FOO","name":"foo"},{"code":"BAR","name":"bar"}]`,
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
			expectedCT:     "application/json",
			expectedBody:   map[string]any{"error": errorsapi.ErrInvalidJSONBody.Error()},
		},
		{
			name:      "conflict to insert batch registers",
			fakeError: errorsapi.ErrRepositoryCategoryAlreadyExists,
			requestBody: []categories.CreateCategoryInput{
				{Code: "FOO", Name: "foo"},
				{Code: "BAR", Name: "bar"},
			},
			contentType:    "application/json",
			expectedStatus: http.StatusConflict,
			expectedCT:     "application/json",
			expectedBody:   map[string]any{"error": errorsapi.ErrRepositoryCategoryAlreadyExists.Error()},
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
			assert.Equal(t, res.Header.Get("Content-Type"), tt.expectedCT)

			rawExpectedBody, err := json.Marshal(tt.expectedBody)
			require.NoError(t, err)

			assert.JSONEq(t, string(rawExpectedBody), string(rawBody))
		})
	}
}
