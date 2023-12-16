package entity

import "fmt"

var (
	InvalidCodeErr           = fmt.Errorf("code is invalid")
	RequiredTitleErr         = fmt.Errorf("title is required")
	RequiredReferenceErr     = fmt.Errorf("reference is required")
	RequiredDescriptionErr   = fmt.Errorf("description is required")
	InvalidPriceErr          = fmt.Errorf("price is invalid")
	DuplicatedProductCodeErr = fmt.Errorf("a product with this code already exists")
	ProductNotFoundErr       = fmt.Errorf("product doesn't exists")
)
