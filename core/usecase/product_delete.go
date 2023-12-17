package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/lbsti/eulabs-challenge/core/repository"
)

type ProductDelete struct {
	repository repository.ProductRepository
}

func NewProductDelete(productRepo repository.ProductRepository) *ProductDelete {
	return &ProductDelete{
		repository: productRepo,
	}
}

func (p *ProductDelete) Execute(ctx context.Context, code string) (bool, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(ProductDefaultTimeout))
	defer cancel()
	isDeleted, err := p.repository.DeleteByCode(ctxWithTimeout, code)

	if err != nil {
		slog.Error("impossible to delete product: %v", err)
		return false, err
	}
	return isDeleted, nil
}
