package usecase_test

import (
	"context"
	"testing"

	"github.com/lbsti/eulabs-challenge/core/entity"
	"github.com/lbsti/eulabs-challenge/core/usecase"
	"github.com/lbsti/eulabs-challenge/infra/repository"
	"github.com/stretchr/testify/assert"
)

func TestProductGet_Execute(t *testing.T) {
	t.Run("Should get a product by code with success", productGetByCodeSuccess)
	t.Run("Should results a error if product not found", productGetByCodeNotFoundErr)
}

func productGetByCodeSuccess(t *testing.T) {
	productRepoInMemory := repository.NewProductRepositoryInMemory()
	productGet := usecase.NewProductGet(productRepoInMemory)
	productGetOutputDTO, e := productGet.Execute(context.TODO(), "XSZ-000741")
	assert.Nil(t, e)
	assert.GreaterOrEqual(t, productGetOutputDTO.ID, int64(1))
}

func productGetByCodeNotFoundErr(t *testing.T) {
	productRepoInMemory := repository.ProductRepositoryInMemorySpy{
		ExpectedError: entity.ProductNotFoundErr,
	}
	productGet := usecase.NewProductGet(productRepoInMemory)
	productGetOutputDTO, e := productGet.Execute(context.TODO(), "XSZ-000741")
	assert.NotNil(t, e)
	assert.Equal(t, productGetOutputDTO.ID, int64(0))
}
