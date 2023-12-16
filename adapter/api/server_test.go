package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/lbsti/eulabs-challenge/core/entity"
	"github.com/lbsti/eulabs-challenge/core/usecase"
	"github.com/lbsti/eulabs-challenge/infra/repository"
	"github.com/stretchr/testify/assert"
)

func TestWebServer_handleProductCreate(t *testing.T) {
	t.Run("Should handle create product request with success", createProductSuccess)
	t.Run("Should results error if create product request is duplicated", createProductDuplicatedErr)
	t.Run("Should results error if create product request send invalid payload", createProductBindErr)
	t.Run("Should results error if create product request send empty code", createProductEmptyCodeErr)
}

func createProductSuccess(t *testing.T) {
	productInMemoryRepo := repository.NewProductRepositoryInMemory()
	ws := NewWebServer("8080", productInMemoryRepo)

	e := echo.New()
	var productInput usecase.ProductInputDTO
	productInput.Code = "XXCC"
	productInput.Description = "Description"
	productInput.PriceInCents = int64(2500)
	productInput.Reference = "XZsdf5tY-AA"
	productInput.Title = "Toy"

	body, err := json.Marshal(productInput)
	assert.NoError(t, err)
	rec := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	echoCtx := e.NewContext(req, rec)

	err2 := ws.handleProductCreate(echoCtx)
	assert.Nil(t, err2)

	var expectedResult usecase.ProductOutputDTO
	err3 := json.Unmarshal(rec.Body.Bytes(), &expectedResult)
	assert.Nil(t, err3)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, expectedResult.Reference, productInput.Reference)
	assert.Greater(t, expectedResult.ID, int64(0))
}

func createProductDuplicatedErr(t *testing.T) {
	productInMemoryRepo := repository.ProductRepositoryInMemorySpy{
		ExpectedError: entity.DuplicatedProductCodeErr,
	}
	ws := NewWebServer("8080", productInMemoryRepo)

	e := echo.New()
	var productInput usecase.ProductInputDTO
	productInput.Code = "XXCC"
	productInput.Description = "Description"
	productInput.PriceInCents = int64(2500)
	productInput.Reference = "XZsdf5tY-AA"
	productInput.Title = "Toy"

	body, err := json.Marshal(productInput)
	assert.NoError(t, err)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	grApi := e.Group("/api")
	grApi.POST("/v1/products", ws.handleProductCreate)
	e.ServeHTTP(rec, req)
	assert.Equal(t, "{\"message\":\"a product with this code already exists\"}\n", rec.Body.String())
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func createProductBindErr(t *testing.T) {
	productInMemoryRepo := repository.NewProductRepositoryInMemory()
	ws := NewWebServer("8080", productInMemoryRepo)
	e := echo.New()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader([]byte(`{`)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	grApi := e.Group("/api")
	grApi.POST("/v1/products", ws.handleProductCreate)
	e.ServeHTTP(rec, req)
	assert.Equal(t, "{\"message\":\"unexpected EOF\"}\n", rec.Body.String())
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func createProductEmptyCodeErr(t *testing.T) {
	productInMemoryRepo := repository.NewProductRepositoryInMemory()
	ws := NewWebServer("8080", productInMemoryRepo)

	e := echo.New()
	var productInput usecase.ProductInputDTO
	productInput.Description = "Description"
	productInput.PriceInCents = int64(2500)
	productInput.Reference = "XZsdf5tY-AA"
	productInput.Title = "Toy"

	body, err := json.Marshal(productInput)
	assert.NoError(t, err)
	rec := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	grApi := e.Group("/api")
	grApi.POST("/v1/products", ws.handleProductCreate)
	e.ServeHTTP(rec, req)
	assert.Equal(t, "{\"message\":\"code is invalid\"}\n", rec.Body.String())
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
