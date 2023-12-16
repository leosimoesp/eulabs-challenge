package entity

import "strings"

type Product struct {
	Title        string
	Description  string
	Code         string
	Reference    string
	ID           int64
	PriceInCents int64
}

func NewProduct() *Product {
	return &Product{}
}

func (p Product) IsValid() error {
	if isEmpty := isEmpty(p.Code); isEmpty {
		return InvalidCodeErr
	}
	if isEmpty := isEmpty(p.Description); isEmpty {
		return RequiredDescriptionErr
	}
	if strings.Contains(p.Code, " ") {
		return InvalidCodeErr
	}
	if isEmpty := isEmpty(p.Title); isEmpty {
		return RequiredTitleErr
	}
	if isEmpty := isEmpty(p.Reference); isEmpty {
		return RequiredReferenceErr
	}
	if p.PriceInCents <= 0 {
		return InvalidPriceErr
	}
	return nil

}

func isEmpty(value string) bool {
	rawValue := strings.ReplaceAll(value, " ", "")
	return rawValue == ""
}
