package entity_test

import (
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/lbsti/eulabs-challenge/core/entity"
	"github.com/stretchr/testify/assert"
)

func TestProduct_IsValid(t *testing.T) {
	t.Run("Should validate with success a product", productIsValid)
	t.Run("Should results error if code is empty", productCodeEmpty)
	t.Run("Should results error if code is empty", productDescriptionEmpty)
	t.Run("Should results error if code is empty", productCodeWithSpace)
	t.Run("Should results error if title is empty", productTitleEmpty)
	t.Run("Should results error if reference is empty", productReferenceEmpty)
	t.Run("Should results error if price is less than zero", productPriceInvalid)
}

func productIsValid(t *testing.T) {
	product := entity.NewProduct()
	product.Code = "0001-DEF-UDSE-14587"
	product.PriceInCents = int64(10000)
	product.Title = "exacqVision® VMS"
	product.Description = `The exacqVision® VMS (Video Management System) software
	installs on commercial off-the-shelf (COTS) servers running
	Windows or Linux operating systems to create an advanced
	security solution, providing recording of the latest, state-of-
	the-art IP video surveillance cameras.`
	reference := uuid.FromStringOrNil("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	product.Reference = reference.String()
	assert.NoError(t, product.IsValid())
}

func productCodeEmpty(t *testing.T) {
	product := entity.NewProduct()
	product.Title = "exacqVision® VMS"
	product.PriceInCents = int64(10000)
	product.Description = `The exacqVision® VMS (Video Management System) software
	installs on commercial off-the-shelf (COTS) servers running
	Windows or Linux operating systems to create an advanced
	security solution, providing recording of the latest, state-of-
	the-art IP video surveillance cameras.`
	reference := uuid.FromStringOrNil("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	product.Reference = reference.String()
	assert.EqualError(t, product.IsValid(), entity.InvalidCodeErr.Error())
}

func productDescriptionEmpty(t *testing.T) {
	product := entity.NewProduct()
	product.Title = "exacqVision® VMS"
	product.Code = "0001-DEF-UDSE-14587"
	product.PriceInCents = int64(10000)
	reference := uuid.FromStringOrNil("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	product.Reference = reference.String()
	assert.EqualError(t, product.IsValid(), entity.RequiredDescriptionErr.Error())
}

func productCodeWithSpace(t *testing.T) {
	product := entity.NewProduct()
	product.PriceInCents = int64(10000)
	product.Code = "0001- DEF-UDSE-14587"
	product.Title = "exacqVision® VMS"
	product.Description = `The exacqVision® VMS (Video Management System) software
	installs on commercial off-the-shelf (COTS) servers running
	Windows or Linux operating systems to create an advanced
	security solution, providing recording of the latest, state-of-
	the-art IP video surveillance cameras.`
	reference := uuid.FromStringOrNil("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	product.Reference = reference.String()
	assert.EqualError(t, product.IsValid(), entity.InvalidCodeErr.Error())
}

func productTitleEmpty(t *testing.T) {
	product := entity.NewProduct()
	product.PriceInCents = int64(10000)
	product.Code = "0001-DEF-UDSE-14587"
	product.Description = `The exacqVision® VMS (Video Management System) software
	installs on commercial off-the-shelf (COTS) servers running
	Windows or Linux operating systems to create an advanced
	security solution, providing recording of the latest, state-of-
	the-art IP video surveillance cameras.`
	reference := uuid.FromStringOrNil("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	product.Reference = reference.String()
	assert.EqualError(t, product.IsValid(), entity.RequiredTitleErr.Error())
}

func productReferenceEmpty(t *testing.T) {
	product := entity.NewProduct()
	product.PriceInCents = int64(10000)
	product.Code = "0001-DEF-UDSE-14587"
	product.Title = "exacqVision® VMS"
	product.Description = `The exacqVision® VMS (Video Management System) software
	installs on commercial off-the-shelf (COTS) servers running
	Windows or Linux operating systems to create an advanced
	security solution, providing recording of the latest, state-of-
	the-art IP video surveillance cameras.`
	assert.EqualError(t, product.IsValid(), entity.RequiredReferenceErr.Error())
}

func productPriceInvalid(t *testing.T) {
	product := entity.NewProduct()
	product.Code = "0001-DEF-UDSE-14587"
	product.Title = "exacqVision® VMS"
	product.Description = `The exacqVision® VMS (Video Management System) software
	installs on commercial off-the-shelf (COTS) servers running
	Windows or Linux operating systems to create an advanced
	security solution, providing recording of the latest, state-of-
	the-art IP video surveillance cameras.`
	reference := uuid.FromStringOrNil("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	product.Reference = reference.String()
	assert.EqualError(t, product.IsValid(), entity.InvalidPriceErr.Error())
}
