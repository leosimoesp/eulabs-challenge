package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/lbsti/eulabs-challenge/core/repository"
)

type ProductGet struct {
	repository repository.ProductRepository
}

type ProductGetOutputDTO struct {
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Code         string `json:"code"`
	Reference    string `json:"reference"`
	PriceInCents int64  `json:"priceInCents"`
	ID           int64  `json:"id"`
}

func NewProductGet(productRepo repository.ProductRepository) *ProductGet {
	return &ProductGet{
		repository: productRepo,
	}
}

func (p *ProductGet) Execute(ctx context.Context, code string) (ProductGetOutputDTO, error) {

	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(ProductDefaultTimeout))
	defer cancel()
	productData, err := p.repository.GetByCode(ctxWithTimeout, code)

	if err != nil {
		slog.Error("impossible to get product", slog.Any("msg", err))
		return ProductGetOutputDTO{}, err
	}
	return ProductGetOutputDTO{
		ID:           productData.ID,
		Title:        productData.Title,
		CreatedAt:    productData.CreatedAt,
		UpdatedAt:    productData.UpdatedAt,
		Reference:    productData.Reference,
		Code:         productData.Code,
		PriceInCents: productData.PriceInCents,
		Description:  productData.Description,
	}, nil
}
