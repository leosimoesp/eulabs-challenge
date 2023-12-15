package usecase

import (
	"github.com/lbsti/eulabs-challenge/core/entity"
	"github.com/lbsti/eulabs-challenge/core/repository"
)

type ProductCreate struct {
	repository repository.ProductRepository
}

type ProductInputDTO struct {
	Title        string
	Description  string
	Code         string
	Reference    string
	PriceInCents int64
}

type ProductOutputDTO struct {
	CreatedAt string
	Reference string
	ID        int64
}

func NewProductCreate(productRepo repository.ProductRepository) *ProductCreate {
	return &ProductCreate{
		repository: productRepo,
	}
}

func (p *ProductCreate) Execute(input ProductInputDTO) (ProductOutputDTO, error) {
	if e := p.validate(input); e != nil {
		return ProductOutputDTO{}, e
	}

	productData, err := p.repository.Insert(repository.ProductRepositoryInput{
		Title:        input.Title,
		Description:  input.Description,
		Code:         input.Code,
		Reference:    input.Reference,
		PriceInCents: input.PriceInCents,
	})

	if err != nil {
		return ProductOutputDTO{}, err
	}

	return ProductOutputDTO{
		CreatedAt: productData.CreatedAt,
		Reference: productData.Reference,
		ID:        productData.ID,
	}, nil
}

func (p *ProductCreate) validate(input ProductInputDTO) error {
	product := entity.NewProduct()
	product.Title = input.Title
	product.Description = input.Description
	product.Code = input.Code
	product.Reference = input.Reference
	product.PriceInCents = input.PriceInCents

	if e := product.IsValid(); e != nil {
		return e
	}
	return nil
}
