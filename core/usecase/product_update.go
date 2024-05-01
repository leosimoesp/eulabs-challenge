package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/lbsti/eulabs-challenge/core/repository"
)

type ProductUpdate struct {
	repository repository.ProductRepository
}

func NewProductUpdate(productRepo repository.ProductRepository) *ProductUpdate {
	return &ProductUpdate{
		repository: productRepo,
	}
}

func (p *ProductUpdate) Execute(ctx context.Context, input ProductInputDTO) error {
	if err := validate(input); err != nil {
		return err
	}
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(ProductDefaultTimeout))
	defer cancel()

	productData, err := p.repository.GetByCode(ctx, input.Code)

	if err != nil {
		slog.Error("impossible to update product", slog.Any("msg", err))
		return err
	}

	return p.repository.Update(ctxWithTimeout, repository.ProductRepositoryInput{
		Title:        input.Title,
		Description:  input.Description,
		Code:         productData.Code,
		Reference:    input.Reference,
		PriceInCents: input.PriceInCents,
	})
}
