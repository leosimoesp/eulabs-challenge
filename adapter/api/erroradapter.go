package api

import (
	"net/http"

	"github.com/lbsti/eulabs-challenge/core/entity"
)

type MappedError struct {
	ResultErr error
	Code      int
}

func Mapping(input error) MappedError {
	switch input {
	case entity.InvalidCodeErr,
		entity.RequiredReferenceErr,
		entity.InvalidPriceErr,
		entity.RequiredDescriptionErr,
		entity.RequiredTitleErr:
		return MappedError{
			ResultErr: input,
			Code:      http.StatusBadRequest,
		}
	case entity.DuplicatedProductCodeErr:
		return MappedError{
			ResultErr: input,
			Code:      http.StatusConflict,
		}
	}
	return MappedError{
		ResultErr: input,
		Code:      http.StatusInternalServerError,
	}
}
