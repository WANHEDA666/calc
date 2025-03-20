package models

type Currency struct {
	Code     string  `json:"code"`
	Exchange float64 `json:"exchange"`
}
