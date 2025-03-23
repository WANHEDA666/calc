package repository

import (
	"calc/internal/models"
	"encoding/xml"
	"errors"
	"fmt"
	"golang.org/x/net/html/charset"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ValCurs struct {
	Currencies []Currency `xml:"Valute"`
	Date       string     `xml:"Date,attr"`
}

type Currency struct {
	CharCode string `xml:"CharCode"`
	Nominal  int    `xml:"Nominal"`
	Name     string `xml:"Name"`
	Value    string `xml:"Value"`
}

type CurrencyRepository struct {
	currencies map[int][]models.Currency
}

func NewCurrencyRepository() *CurrencyRepository {
	return &CurrencyRepository{
		currencies: make(map[int][]models.Currency),
	}
}

func (r *CurrencyRepository) GetAll() ([]models.Currency, error) {
	if _, ok := r.currencies[time.Now().Day()]; ok {
		return r.currencies[time.Now().Day()], nil
	}

	var result []models.Currency

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.cbr.ru/scripts/XML_daily.asp?date_req=%s", time.Now().Format("02/01/2006")), nil)
	if err != nil {
		return nil, errors.New("Ошибка при создании запроса: " + err.Error())
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("Ошибка при выполнении запроса: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Ошибка: сервер ЦБ вернул статус-код %d", resp.StatusCode))
	}

	decoder := xml.NewDecoder(resp.Body)
	decoder.CharsetReader = charset.NewReaderLabel

	var valCurs ValCurs
	err = decoder.Decode(&valCurs)
	if err != nil {
		return nil, errors.New("Ошибка при парсинге XML: " + err.Error())
	}

	for _, currency := range valCurs.Currencies {
		exchangeRate, err := parseExchangeRate(currency.Value)
		if err != nil {
			continue
		}
		result = append(result, models.Currency{
			Code:     currency.CharCode,
			Exchange: exchangeRate,
		})
	}
	result = append(result, models.Currency{
		Code:     "RUB",
		Exchange: 1.0,
	})
	sort.Slice(result, func(i, j int) bool {
		return result[i].Code < result[j].Code
	})
	return result, nil
}

func parseExchangeRate(valueStr string) (float64, error) {
	valueStr = strings.Replace(valueStr, ",", ".", -1)
	return strconv.ParseFloat(valueStr, 64)
}
