package restapi_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"test-assigment-ltp/internal/restapi"
)

type MockKrakenService struct {
	mock.Mock
}

func (m *MockKrakenService) GetLastTradedPrice(pair string) (float64, error) {
	args := m.Called(pair)
	return args.Get(0).(float64), args.Error(1)
}

func TestGetLTP(t *testing.T) {
	e := echo.New()
	mockService := new(MockKrakenService)
	handler := restapi.NewLTPHanlder(mockService)

	t.Run("Missing 'pair' parameter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ltp", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetLTP(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		expected := `{"error":"Missing 'pair' parameter"}`
		assert.JSONEq(t, expected, rec.Body.String())
	})

	t.Run("Invalid pair format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ltp?pair=INVALID", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetLTP(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusPartialContent, rec.Code)
		expected := `{
			"ltp": [],
			"errors": ["Invalid pair format: INVALID"]
		}`
		assert.JSONEq(t, expected, rec.Body.String())
	})

	t.Run("Successful LTP retrieval", func(t *testing.T) {
		mockService.On("GetLastTradedPrice", "BTC/USD").Return(50000.0, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/ltp?pair=BTC/USD", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetLTP(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		expected := `{
			"ltp": [{"pair": "BTC/USD", "amount": 50000}]
		}`
		assert.JSONEq(t, expected, rec.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("Partial failure", func(t *testing.T) {
		mockService.On("GetLastTradedPrice", "BTC/USD").Return(50000.0, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/ltp?pair=BTC/USD&pair=INVALID", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetLTP(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusPartialContent, rec.Code)

		expected := `{
			"ltp": [{"pair": "BTC/USD", "amount": 50000}],
			"errors": ["Invalid pair format: INVALID"]
		}`
		assert.JSONEq(t, expected, rec.Body.String())
		mockService.AssertExpectations(t)
	})
}
