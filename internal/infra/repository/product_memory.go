package repository

import (
	"context"
	"math/rand"
	"time"

	"github.com/lbsti/eulabs-challenge/internal/core/repository"
)

type ProductRepositoryInMemory struct{}

func NewProductRepositoryInMemory() repository.ProductRepository {
	return ProductRepositoryInMemory{}
}

func (r ProductRepositoryInMemory) Insert(ctx context.Context,
	in repository.ProductRepositoryInput) (repository.ProductRepositoryData, error) {
	now := time.Now().UTC()

	return repository.ProductRepositoryData{
		Reference: in.Reference,
		CreatedAt: now.Format("2006-01-02 15:04:05"),
		ID:        rand.Int63n(10000) + 1,
	}, nil
}

func (r ProductRepositoryInMemory) GetByCode(ctx context.Context, code string) (repository.ProductRepositoryData, error) {
	now := time.Now().UTC()

	return repository.ProductRepositoryData{
		Reference:    "RF009-pods74",
		CreatedAt:    now.Format("2006-01-02 15:04:05"),
		ID:           rand.Int63n(10000) + 1,
		Title:        "Toy",
		Description:  "Blahahhs",
		Code:         "XSZ-000741",
		PriceInCents: int64(51400),
	}, nil
}

func (ProductRepositoryInMemory) DeleteByCode(ctx context.Context, code string) (bool, error) {
	return true, nil
}

func (ProductRepositoryInMemory) Update(ctx context.Context, in repository.ProductRepositoryInput) error {
	return nil
}

type ProductRepositoryInMemorySpy struct {
	ExpectedError error
	ExpectedData  repository.ProductRepositoryData
}

func (spyRepo ProductRepositoryInMemorySpy) Insert(ctx context.Context,
	in repository.ProductRepositoryInput) (repository.ProductRepositoryData, error) {
	return spyRepo.ExpectedData, spyRepo.ExpectedError
}

func (spyRepo ProductRepositoryInMemorySpy) GetByCode(ctx context.Context,
	code string) (repository.ProductRepositoryData, error) {
	return spyRepo.ExpectedData, spyRepo.ExpectedError
}

func (spyRepo ProductRepositoryInMemorySpy) DeleteByCode(ctx context.Context,
	code string) (bool, error) {
	return spyRepo.ExpectedError == nil, spyRepo.ExpectedError
}

func (spyRepo ProductRepositoryInMemorySpy) Update(ctx context.Context,
	in repository.ProductRepositoryInput) error {
	return spyRepo.ExpectedError
}
