package repository

import (
	"math/rand"
	"time"

	"github.com/lbsti/eulabs-challenge/core/repository"
)

type ProductRepositoryInMemory struct{}

func NewProductRepositoryInMemory() repository.ProductRepository {
	return ProductRepositoryInMemory{}
}

func (r ProductRepositoryInMemory) Insert(in repository.ProductRepositoryInput) (repository.ProductRepositoryData, error) {
	now := time.Now().UTC()

	return repository.ProductRepositoryData{
		Reference: in.Reference,
		CreatedAt: now.Format("2006-01-02 15:04:05"),
		ID:        rand.Int63n(10000) + 1,
	}, nil
}
