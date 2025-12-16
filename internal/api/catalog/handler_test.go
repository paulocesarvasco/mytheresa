package catalogapi

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/internal/catalog"
	errorsapi "github.com/mytheresa/go-hiring-challenge/internal/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetProducts(t *testing.T) {
	tests := []struct {
		name           string
		fakeProducts   catalog.ProductPage
		fakeError      error
		queryParams    map[string]string
		expectedStatus int
		expectedCT     string
		expectedBody   any
	}{
		{
			name: "fetch products list ok",
			fakeProducts: catalog.ProductPage{
				Total: 2,
				Products: []catalog.ProductView{
					{Category: "foo"},
					{Category: "bar"},
				},
			},
			queryParams: map[string]string{
				"limit":         "2",
				"offset":        "0",
				"category_code": "0xFFFF",
				"max_price":     "9.99",
			},
			expectedStatus: http.StatusOK,
			expectedCT:     "application/json",
			expectedBody: catalog.ProductPage{
				Total: 2,
				Products: []catalog.ProductView{
					{Category: "foo"},
					{Category: "bar"},
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
			fakeProducts: catalog.ProductPage{
				Total: 2,
				Products: []catalog.ProductView{
					{Category: "foo"},
					{Category: "bar"},
				},
			},
			queryParams: map[string]string{
				"limit": "200",
			},
			expectedStatus: http.StatusOK,
			expectedCT:     "application/json",
			expectedBody: catalog.ProductPage{
				Total: 2,
				Products: []catalog.ProductView{
					{Category: "foo"},
					{Category: "bar"},
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
			name: "invalid max price",
			queryParams: map[string]string{
				"max_price": "0.00",
			},
			expectedStatus: http.StatusBadRequest,
			expectedCT:     "application/json",
			expectedBody:   map[string]any{"error": errorsapi.ErrInvalidMaxPriceParam.Error()},
		},
		{
			name:           "catalog service error",
			fakeProducts:   catalog.ProductPage{},
			fakeError:      errorsapi.ErrRepositoryFetchProducts,
			expectedStatus: http.StatusInternalServerError,
			expectedCT:     "application/json",
			expectedBody:   map[string]any{"error": errorsapi.ErrRepositoryFetchProducts.Error()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewFakeService()
			s.SetListProductsResponse(tt.fakeProducts, tt.fakeError)
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

func TestGetDetailProduct(t *testing.T) {
	tests := []struct {
		name           string
		productCode    string
		fakeDetails    catalog.ProductView
		fakeError      error
		expectedStatus int
		expectedCT     string
		expectedBody   any
	}{
		{
			name:           "retrieve product details ok",
			productCode:    "PROD001",
			fakeDetails:    catalog.ProductView{Code: "PROD001"},
			expectedStatus: http.StatusOK,
			expectedCT:     "application/json",
			expectedBody:   catalog.ProductView{Code: "PROD001"},
		},
		{
			name:           "invalid product code",
			productCode:    "PROD00A",
			expectedStatus: http.StatusBadRequest,
			expectedCT:     "application/json",
			expectedBody:   map[string]any{"error": errorsapi.ErrInvalidProductCode.Error()},
		},
		{
			name:           "product code not found",
			productCode:    "PROD001",
			fakeError:      errorsapi.ErrProductNotFound,
			expectedStatus: http.StatusNotFound,
			expectedCT:     "application/json",
			expectedBody:   map[string]any{"error": errorsapi.ErrProductNotFound.Error()},
		},
		{
			name:           "catalog service error",
			productCode:    "PROD001",
			fakeError:      errorsapi.ErrRepositoryFetchProduct,
			expectedStatus: http.StatusInternalServerError,
			expectedCT:     "application/json",
			expectedBody:   map[string]any{"error": errorsapi.ErrRepositoryFetchProduct.Error()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewFakeService()
			s.SetDetailProductResponse(tt.fakeDetails, tt.fakeError)
			h := New(s)
			r := Routes(h)

			ts := httptest.NewServer(r)
			defer ts.Close()

			res, err := ts.Client().Get(ts.URL + "/" + tt.productCode)
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
