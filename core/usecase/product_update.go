package usecase

import (
	"context"
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
	if e := validate(input); e != nil {
		return e
	}
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(ProductDefaultTimeout))
	defer cancel()

	return p.repository.Update(ctxWithTimeout, repository.ProductRepositoryInput{
		Title:        input.Title,
		Description:  input.Description,
		Code:         input.Code,
		Reference:    input.Reference,
		PriceInCents: input.PriceInCents,
	})
}
