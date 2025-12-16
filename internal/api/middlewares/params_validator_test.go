package middlewares

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/mytheresa/go-hiring-challenge/internal/logs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseQueryParameters(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    map[string]string
		expectedStatus int
		expectedCT     string
		expectedBody   any
	}{
		{
			name: "valid query parameters",
			queryParams: map[string]string{
				"limit":         "10",
				"offset":        "0",
				"category_code": "0xFFFF",
			},
			expectedStatus: http.StatusOK,
			expectedCT:     "application/json",
			expectedBody:   "{}",
		},
		{
			name: "query parameter limit exceeds maximum allowed value",
			queryParams: map[string]string{
				"limit": "1000",
			},
			expectedStatus: http.StatusOK,
			expectedCT:     "application/json",
			expectedBody:   "{}",
		},
		{
			name: "query parameter limit null value",
			queryParams: map[string]string{
				"limit": "0",
			},
			expectedStatus: http.StatusBadRequest,
			expectedCT:     "application/json",
			expectedBody:   map[string]string{"error": "invalid limit parameter"},
		},
		{
			name: "query parameter limit invalid value",
			queryParams: map[string]string{
				"limit": "\n",
			},
			expectedStatus: http.StatusBadRequest,
			expectedCT:     "application/json",
			expectedBody:   map[string]string{"error": "invalid limit parameter"},
		},
		{
			name: "query parameter offset negative value",
			queryParams: map[string]string{
				"offset": "-1",
			},
			expectedStatus: http.StatusBadRequest,
			expectedCT:     "application/json",
			expectedBody:   map[string]string{"error": "invalid offset parameter"},
		},
		{
			name: "query parameter offset invalid value",
			queryParams: map[string]string{
				"offset": "A",
			},
			expectedStatus: http.StatusBadRequest,
			expectedCT:     "application/json",
			expectedBody:   map[string]string{"error": "invalid offset parameter"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := ParseQueryParameters(logs.Logger())

			r := chi.NewRouter()
			h := NewFakeHandler(http.StatusOK, "{}")

			r.With(m).Get("/", h.WriteJSON)

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

			assert.Equal(t, tt.expectedStatus, res.StatusCode)
			assert.Equal(t, res.Header.Get("Content-Type"), tt.expectedCT)

			rawBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			rawExpectedBody, err := json.Marshal(tt.expectedBody)
			require.NoError(t, err)

			assert.JSONEq(t, string(rawExpectedBody), string(rawBody))
		})
	}
}

func TestParseMaxPrice(t *testing.T) {
	tests := []struct {
		name               string
		maxPriceQueryValue string
		expectedStatus     int
		expectedCT         string
		expectedBody       any
	}{
		{
			name:               "valid max price value",
			maxPriceQueryValue: "10.00",
			expectedStatus:     http.StatusOK,
			expectedCT:         "application/json",
			expectedBody:       "{}",
		},
		{
			name:           "request without max price filter",
			expectedStatus: http.StatusOK,
			expectedCT:     "application/json",
			expectedBody:   "{}",
		},
		{
			name:               "max price invalid value",
			maxPriceQueryValue: "USD10",
			expectedStatus:     http.StatusBadRequest,
			expectedCT:         "application/json",
			expectedBody:       map[string]string{"error": "invalid max_price parameter"},
		},
		{
			name:               "max price negative value",
			maxPriceQueryValue: "-1.00",
			expectedStatus:     http.StatusBadRequest,
			expectedCT:         "application/json",
			expectedBody:       map[string]string{"error": "invalid max_price parameter"},
		},
		{
			name:               "max price null value",
			maxPriceQueryValue: "-1.00",
			expectedStatus:     http.StatusBadRequest,
			expectedCT:         "application/json",
			expectedBody:       map[string]string{"error": "invalid max_price parameter"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := ParseMaxPrice(logs.Logger())

			r := chi.NewRouter()
			h := NewFakeHandler(http.StatusOK, "{}")

			r.With(m).Get("/", h.WriteJSON)

			ts := httptest.NewServer(r)
			defer ts.Close()

			u, err := url.Parse(ts.URL)
			require.NoError(t, err)

			q := u.Query()
			q.Set("max_price", tt.maxPriceQueryValue)

			u.RawQuery = q.Encode()

			res, err := ts.Client().Get(u.String())
			require.NoError(t, err)
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)
			assert.Equal(t, res.Header.Get("Content-Type"), tt.expectedCT)

			rawBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			rawExpectedBody, err := json.Marshal(tt.expectedBody)
			require.NoError(t, err)

			assert.JSONEq(t, string(rawExpectedBody), string(rawBody))
		})
	}
}

func TestValidateProductCode(t *testing.T) {
	tests := []struct {
		name           string
		productCode    string
		expectedStatus int
		expectedCT     string
		expectedBody   any
	}{
		{
			name:           "valid product code format",
			productCode:    "PROD001",
			expectedStatus: http.StatusOK,
			expectedCT:     "application/json",
			expectedBody:   "{}",
		},
		{
			name:           "invalid product code format",
			productCode:    "PRODFFF",
			expectedStatus: http.StatusBadRequest,
			expectedCT:     "application/json",
			expectedBody:   map[string]string{"error": "invalid product code"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := ValidateProductCode(logs.Logger())

			r := chi.NewRouter()
			h := NewFakeHandler(http.StatusOK, "{}")

			r.With(m).Get("/{code}", h.WriteJSON)

			ts := httptest.NewServer(r)
			defer ts.Close()

			res, err := ts.Client().Get(fmt.Sprintf("%s/%s", ts.URL, tt.productCode))
			require.NoError(t, err)
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)
			assert.Equal(t, res.Header.Get("Content-Type"), tt.expectedCT)

			rawBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			rawExpectedBody, err := json.Marshal(tt.expectedBody)
			require.NoError(t, err)

			assert.JSONEq(t, string(rawExpectedBody), string(rawBody))
		})
	}
}
