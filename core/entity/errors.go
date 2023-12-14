package entity

import "fmt"

var (
	InvalidCodeErr       = fmt.Errorf("code is invalid")
	RequiredTitleErr     = fmt.Errorf("title is required")
	RequiredReferenceErr = fmt.Errorf("reference is required")
	InvalidPriceErr      = fmt.Errorf("price is invalid")
)
