package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/mytheresa/go-hiring-challenge/internal/logs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateJSON(t *testing.T) {
	type TestRequestPayload struct {
		Info string `json:"info" validate:"required"`
	}
	tests := []struct {
		name           string
		testBody       any
		testCT         string
		fakeBody       any
		expectedStatus int
		expectedCT     string
		expectedBody   any
	}{
		{
			name:           "valid JSON",
			testBody:       map[string]string{"info": "foo"},
			testCT:         "application/json",
			fakeBody:       TestRequestPayload{Info: "foo"},
			expectedStatus: http.StatusOK,
			expectedCT:     "application/json",
			expectedBody:   map[string]string{"info": "foo"},
		},
		{
			name:           "invalid content type",
			testBody:       `{"info": "foo"}`,
			testCT:         "text/plain",
			expectedStatus: http.StatusUnsupportedMediaType,
			expectedCT:     "application/json",
			expectedBody:   map[string]string{"error": "content type must be application/json"},
		},
		{
			name:           "invalid JSON",
			testBody:       `{info: invalid_json}`,
			testCT:         "application/json",
			expectedStatus: http.StatusBadRequest,
			expectedCT:     "application/json",
			expectedBody:   map[string]string{"error": "invalid json body"},
		},
		{
			name:           "invalid schema",
			testBody:       TestRequestPayload{},
			testCT:         "application/json",
			expectedStatus: http.StatusBadRequest,
			expectedCT:     "application/json",
			expectedBody:   map[string]string{"error": "invalid request schema"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := ValidateJSON[TestRequestPayload](logs.Logger())

			r := chi.NewRouter()
			h := NewFakeHandler(http.StatusOK, tt.fakeBody)

			r.With(m).Post("/", h.WriteJSON)

			ts := httptest.NewServer(r)
			defer ts.Close()

			rawTestBody, err := json.Marshal(tt.testBody)
			require.NoError(t, err)

			res, err := ts.Client().Post(ts.URL, tt.testCT, bytes.NewReader(rawTestBody))
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
