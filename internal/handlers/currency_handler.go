package handlers

import (
	"calc/internal/repository"
	"encoding/json"
	"net/http"
	"strings"
)

type CurrencyHandler struct {
	repository *repository.CurrencyRepository
}

func NewCurrencyHandler(repository *repository.CurrencyRepository) *CurrencyHandler {
	return &CurrencyHandler{repository: repository}
}

func (h *CurrencyHandler) GetCurrencies(writer http.ResponseWriter, _ *http.Request) {
	currencies, err := h.repository.GetAll()
	if err != nil {
		http.Error(writer, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(currencies)
}

func (h *CurrencyHandler) ConvertCurrency(writer http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	parts := strings.Split(strings.TrimPrefix(path, "/convert/"), "-")
	if len(parts) != 2 {
		http.Error(writer, `{"error": "Invalid currency format. Use USD-RUB"}`, http.StatusBadRequest)
		return
	}

	currencies, err := h.repository.GetAll()
	if err != nil {
		http.Error(writer, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	from := strings.ToUpper(parts[0])
	to := strings.ToUpper(parts[1])

	if from == "RUB" {
		for _, currency := range currencies {
			if currency.Code == to {
				writer.Header().Set("Content-Type", "application/json")
				json.NewEncoder(writer).Encode(1.0 / currency.Exchange)
				return
			}
		}
		http.Error(writer, `{"error": "One or both currencies not found"}`, http.StatusNotFound)
		return
	} else if to == "RUB" {
		for _, currency := range currencies {
			if currency.Code == from {
				writer.Header().Set("Content-Type", "application/json")
				json.NewEncoder(writer).Encode(currency.Exchange)
				return
			}
		}
		http.Error(writer, `{"error": "One or both currencies not found"}`, http.StatusNotFound)
		return
	}

	var fromRate, toRate float64
	foundFrom, foundTo := false, false

	for _, currency := range currencies {
		if currency.Code == from {
			fromRate = currency.Exchange
			foundFrom = true
		}
		if currency.Code == to {
			toRate = currency.Exchange
			foundTo = true
		}
	}

	if !foundFrom || !foundTo {
		http.Error(writer, `{"error": "One or both currencies not found"}`, http.StatusNotFound)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(fromRate / toRate)
}
