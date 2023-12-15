package usecase_test

import (
	"errors"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/lbsti/eulabs-challenge/core/entity"
	productrepo "github.com/lbsti/eulabs-challenge/core/repository"
	"github.com/lbsti/eulabs-challenge/core/usecase"
	"github.com/lbsti/eulabs-challenge/infra/repository"
	"github.com/stretchr/testify/assert"
)

type ProductRepositoryInMemorySpy struct {
	ExpectedError error
	ExpectedData  productrepo.ProductRepositoryData
}

func (spyRepo ProductRepositoryInMemorySpy) Insert(in productrepo.ProductRepositoryInput) (productrepo.ProductRepositoryData, error) {
	return spyRepo.ExpectedData, spyRepo.ExpectedError
}

func TestProductCreate_Create(t *testing.T) {
	t.Run("Should create a product with success", productCreateSuccess)
	t.Run("Should result error if product is invalid", productCreateInvalid)
	t.Run("Should result error if repository timeout", productRepositoryTimeout)
}

func productCreateSuccess(t *testing.T) {
	productRepoInMemory := repository.NewProductRepositoryInMemory()
	productCreate := usecase.NewProductCreate(productRepoInMemory)
	reference := uuid.FromStringOrNil("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	productOutDTO, e := productCreate.Execute(usecase.ProductInputDTO{
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
	assert.GreaterOrEqual(t, productOutDTO.ID, int64(1))
}

func productCreateInvalid(t *testing.T) {
	productRepoInMemory := repository.NewProductRepositoryInMemory()
	productCreate := usecase.NewProductCreate(productRepoInMemory)
	reference := uuid.FromStringOrNil("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	productOutDTO, e := productCreate.Execute(usecase.ProductInputDTO{
		Title: "exacqVision® VMS",
		Description: `The exacqVision® VMS (Video Management System) software
		installs on commercial off-the-shelf (COTS) servers running
		Windows or Linux operating systems to create an advanced
		security solution, providing recording of the latest, state-of-
		the-art IP video surveillance cameras.`,
		PriceInCents: int64(10000),
		Reference:    reference.String(),
	})
	assert.NotNil(t, e)
	assert.EqualError(t, e, entity.InvalidCodeErr.Error())
	assert.Equal(t, productOutDTO, usecase.ProductOutputDTO{})
}

func productRepositoryTimeout(t *testing.T) {
	timeoutErr := errors.New("timeout")
	productRepoInMemory := ProductRepositoryInMemorySpy{ExpectedError: timeoutErr}
	productCreate := usecase.NewProductCreate(productRepoInMemory)
	reference := uuid.FromStringOrNil("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	productOutDTO, e := productCreate.Execute(usecase.ProductInputDTO{
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
	assert.NotNil(t, e)
	assert.EqualError(t, e, timeoutErr.Error())
	assert.Equal(t, productOutDTO, usecase.ProductOutputDTO{})
}
