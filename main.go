package main

import (
	"calc/internal/handlers"
	"calc/internal/repository"
	"fmt"
	"net/http"
	"os"
)

func main() {
	currencyRepository := repository.NewCurrencyRepository()
	handler := handlers.NewCurrencyHandler(currencyRepository)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/currencies", handler.GetCurrencies)

	err := http.ListenAndServe("0.0.0.0:"+port, nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}
