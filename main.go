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

	mux := http.NewServeMux()
	mux.HandleFunc("/currencies", handler.GetCurrencies)
	mux.HandleFunc("/convert/", handler.ConvertCurrency)

	startCurrencyFetcher()

	handlerWithCORS := corsMiddleware(mux)
	err := http.ListenAndServe("0.0.0.0:"+port, handlerWithCORS)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
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
