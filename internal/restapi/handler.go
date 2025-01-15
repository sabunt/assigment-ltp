package restapi

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"test-assigment-ltp/internal/models"
	"test-assigment-ltp/internal/services"
	"test-assigment-ltp/pkg/utils"
)

type LTPHandler struct {
	krakenService  services.KrakenServiceInterface
	maxConcurrency int
}

func NewLTPHanlder(krakenService services.KrakenServiceInterface, maxConcurrency int) *LTPHandler {
	return &LTPHandler{krakenService: krakenService, maxConcurrency: maxConcurrency}
}

func (h *LTPHandler) GetLTP(c echo.Context) error {
	pairs := c.QueryParams()["pair"]
	if len(pairs) == 0 {
		logrus.Warn("Missing 'pair' parameter in request")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing 'pair' parameter"})
	}

	var (
		ltpList []models.LastTradedPrice
		errors  []string
		wg      sync.WaitGroup
		mu      sync.Mutex
	)

	semaphore := make(chan struct{}, h.maxConcurrency)

	for _, pair := range pairs {
		if err := utils.ValidatePair(pair); err != nil {
			logrus.Errorf("Invalid pair format: %v", err)
			errors = append(errors, fmt.Sprintf("Invalid pair format: %s", pair))
			continue
		}

		wg.Add(1)

		semaphore <- struct{}{}

		go func(pair string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			amount, err := h.krakenService.GetLastTradedPrice(pair)
			if err != nil {
				logrus.Errorf("Failed to fetch LTP for pair %s: %v", pair, err)
				mu.Lock()
				errors = append(errors, fmt.Sprintf("Failed to fetch LTP for pair %s", pair))
				mu.Unlock()
				return
			}

			mu.Lock()
			ltpList = append(ltpList, models.LastTradedPrice{
				Pair:   pair,
				Amount: amount,
			})
			mu.Unlock()
		}(pair)
	}

	wg.Wait()

	if len(errors) > 0 {
		return c.JSON(http.StatusPartialContent, models.PartialLastTradedPriceResponse{
			Ltp:    utils.ConvertToPointerSlice(ltpList),
			Errors: errors,
		})
	}

	response := models.LastTradedPriceResponse{Ltp: utils.ConvertToPointerSlice(ltpList)}
	return c.JSON(http.StatusOK, response)
}
