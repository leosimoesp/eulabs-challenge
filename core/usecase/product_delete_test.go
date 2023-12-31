package usecase_test

import (
	"context"
	"testing"

	"github.com/lbsti/eulabs-challenge/core/entity"
	"github.com/lbsti/eulabs-challenge/core/usecase"
	"github.com/lbsti/eulabs-challenge/infra/repository"
	"github.com/stretchr/testify/assert"
)

func TestProductDelete_Execute(t *testing.T) {
	t.Run("Should delete a product with success", productDeleteSuccess)
	t.Run("Should results an error if product code was not found", productDeleteNotFoundErr)
}

func productDeleteSuccess(t *testing.T) {
	productRepoInMemory := repository.NewProductRepositoryInMemory()
	productDelete := usecase.NewProductDelete(productRepoInMemory)

	isDeleted, e := productDelete.Execute(context.TODO(), "XSZ-000741")
	assert.Nil(t, e)
	assert.True(t, isDeleted)
}

func productDeleteNotFoundErr(t *testing.T) {
	productRepoInMemory := repository.ProductRepositoryInMemorySpy{
		ExpectedError: entity.ProductNotFoundErr,
	}
	productDelete := usecase.NewProductDelete(productRepoInMemory)

	isDeleted, e := productDelete.Execute(context.TODO(), "XSZ-000741")
	assert.NotNil(t, e)
	assert.False(t, isDeleted)
}
