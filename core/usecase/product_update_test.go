package usecase_test

import (
	"context"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/lbsti/eulabs-challenge/core/entity"
	"github.com/lbsti/eulabs-challenge/core/usecase"
	"github.com/lbsti/eulabs-challenge/infra/repository"
	"github.com/stretchr/testify/assert"
)

func TestProductUpdate_Execute(t *testing.T) {
	t.Run("Should update a product with success", productUpdateSuccess)
	t.Run("Should results an error if product code is empty", productUpdateCodeEmptyErr)
	t.Run("Should results an error if product is invalid", productUpdateInvalidErr)
	t.Run("Should results an error if product does not exist", productUpdateNotFoundErr)
}

func productUpdateSuccess(t *testing.T) {
	productRepoInMemory := repository.NewProductRepositoryInMemory()
	productUpdate := usecase.NewProductUpdate(productRepoInMemory)
	reference := uuid.FromStringOrNil("6ba7b810-9dad-11d1-80b4-00c04fd430c8")

	e := productUpdate.Execute(context.TODO(), usecase.ProductInputDTO{
		Title: "exacqVision® VMS",
		Description: `The exacqVision® VMS (Video Management System) software
	installs on commercial off-the-shelf (COTS) servers running
	Windows or Linux operating systems to create an advanced
	security solution, providing recording of the latest, state-of-
	the-art IP video surveillance cameras.`,
		Code:         "0001-DEF-UDSE-14587",
		PriceInCents: int64(10000),
		Reference:    reference.String(),
	})
	assert.Nil(t, e)
}

func productUpdateCodeEmptyErr(t *testing.T) {
	productRepoInMemory := repository.NewProductRepositoryInMemory()
	productUpdate := usecase.NewProductUpdate(productRepoInMemory)
	reference := uuid.FromStringOrNil("6ba7b810-9dad-11d1-80b4-00c04fd430c8")

	e := productUpdate.Execute(context.TODO(), usecase.ProductInputDTO{
		Title: "exacqVision® VMS",
		Description: `The exacqVision® VMS (Video Management System) software
	installs on commercial off-the-shelf (COTS) servers running
	Windows or Linux operating systems to create an advanced
	security solution, providing recording of the latest, state-of-
	the-art IP video surveillance cameras.`,
		Code:         " ",
		PriceInCents: int64(10000),
		Reference:    reference.String(),
	})
	assert.NotNil(t, e)
	assert.EqualError(t, e, entity.InvalidCodeErr.Error())
}

func productUpdateInvalidErr(t *testing.T) {
	productRepoInMemory := repository.NewProductRepositoryInMemory()
	productUpdate := usecase.NewProductUpdate(productRepoInMemory)
	reference := uuid.FromStringOrNil("6ba7b810-9dad-11d1-80b4-00c04fd430c8")

	e := productUpdate.Execute(context.TODO(), usecase.ProductInputDTO{
		Title:        "exacqVision® VMS",
		Description:  ``,
		Code:         "0001-DEF-UDSE-14587",
		PriceInCents: int64(10000),
		Reference:    reference.String(),
	})
	assert.NotNil(t, e)
	assert.EqualError(t, e, entity.RequiredDescriptionErr.Error())
}

func productUpdateNotFoundErr(t *testing.T) {
	productRepoInMemory := repository.ProductRepositoryInMemorySpy{
		ExpectedError: entity.ProductNotFoundErr,
	}
	productUpdate := usecase.NewProductUpdate(productRepoInMemory)
	reference := uuid.FromStringOrNil("6ba7b810-9dad-11d1-80b4-00c04fd430c8")

	e := productUpdate.Execute(context.TODO(), usecase.ProductInputDTO{
		Title: "exacqVision® VMS",
		Description: `The exacqVision® VMS (Video Management System) software
	installs on commercial off-the-shelf (COTS) servers running
	Windows or Linux operating systems to create an advanced
	security solution, providing recording of the latest, state-of-
	the-art IP video surveillance cameras.`,
		Code:         "0001-DEF-UDSE-0000",
		PriceInCents: int64(10000),
		Reference:    reference.String(),
	})
	assert.NotNil(t, e)
	assert.EqualError(t, e, entity.ProductNotFoundErr.Error())
}
