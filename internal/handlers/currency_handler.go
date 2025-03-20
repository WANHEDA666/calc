package handlers

import (
	"calc/internal/repository"
	"encoding/json"
	"net/http"
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
