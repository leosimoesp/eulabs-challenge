package usecase_test

import (
	"context"
	"testing"

	"github.com/lbsti/eulabs-challenge/internal/core/entity"
	"github.com/lbsti/eulabs-challenge/internal/core/usecase"
	"github.com/lbsti/eulabs-challenge/internal/infra/repository"
	"github.com/stretchr/testify/assert"
)

func TestProductDelete_Execute(t *testing.T) {
	t.Run("Should delete a product with success", productDeleteSuccess)
	t.Run("Should results an error if product code was not found", productDeleteNotFoundErr)
}

func productDeleteSuccess(t *testing.T) {
	productRepoInMemory := repository.NewProductRepositoryInMemory()
	productDelete := usecase.NewProductDelete(productRepoInMemory)

	isDeleted, err := productDelete.Execute(context.TODO(), "XSZ-000741")
	assert.Nil(t, err)
	assert.True(t, isDeleted)
}

func productDeleteNotFoundErr(t *testing.T) {
	productRepoInMemory := repository.ProductRepositoryInMemorySpy{
		ExpectedError: entity.ProductNotFoundErr,
	}
	productDelete := usecase.NewProductDelete(productRepoInMemory)

	isDeleted, err := productDelete.Execute(context.TODO(), "XSZ-000741")
	assert.NotNil(t, err)
	assert.False(t, isDeleted)
}
