package models

type Currency struct {
	Name     string  `json:"name"`
	Code     string  `json:"code"`
	Exchange float64 `json:"exchange"`
}
