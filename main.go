package main

import (
	"calc/internal/handlers"
	"calc/internal/repository"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	currencyRepository := repository.NewCurrencyRepository()
	handler := handlers.NewCurrencyHandler(currencyRepository)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/currencies", handler.GetCurrencies)
	http.HandleFunc("/convert/", handler.ConvertCurrency)

	startCurrencyFetcher()

	err := http.ListenAndServe("0.0.0.0:"+port, nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}

func startCurrencyFetcher() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			fetchCurrencies()
		}
	}()
}

func fetchCurrencies() {
	url := "https://calc-j5oi.onrender.com/currencies"
	resp, _ := http.Get(url)
	defer resp.Body.Close()
}
