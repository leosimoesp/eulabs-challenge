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

func TestWebServer_handleProductGet(t *testing.T) {
	t.Run("Should handle get product request with success", getProductSuccess)
	t.Run("Should results error if product code does not exists", getProductNotFoundErr)
}

func getProductSuccess(t *testing.T) {
	productInMemoryRepo := repository.NewProductRepositoryInMemory()
	ws := NewWebServer("8080", productInMemoryRepo)

	e := echo.New()
	rec := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	echoCtx := e.NewContext(req, rec)
	echoCtx.SetPath("products/:code")
	echoCtx.SetParamNames("code")
	echoCtx.SetParamValues("XSZ-000741")

	grApi := e.Group("/api/v1")
	grApi.GET("/products", ws.handleProductGet)

	var productGetOutputDTO usecase.ProductGetOutputDTO
	err := ws.handleProductGet(echoCtx)
	assert.Nil(t, err)

	err2 := json.Unmarshal(rec.Body.Bytes(), &productGetOutputDTO)
	assert.Nil(t, err2)
	assert.Equal(t, "XSZ-000741", productGetOutputDTO.Code)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func getProductNotFoundErr(t *testing.T) {
	productInMemoryRepo := repository.ProductRepositoryInMemorySpy{
		ExpectedError: entity.ProductNotFoundErr,
	}
	ws := NewWebServer("8080", productInMemoryRepo)

	e := echo.New()
	rec := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	echoCtx := e.NewContext(req, rec)
	echoCtx.SetPath("products/:code")
	echoCtx.SetParamNames("code")
	echoCtx.SetParamValues("XSZ-000741")

	grApi := e.Group("/api/v1")
	grApi.GET("/products", ws.handleProductGet)
	e.ServeHTTP(rec, req)

	err := ws.handleProductGet(echoCtx)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestWebServer_handleProductDelete(t *testing.T) {
	t.Run("Should handle delete product request with success", deleteProductSuccess)
	t.Run("Should results error if product code does not exists", deleteProductNotFoundErr)
}

func deleteProductSuccess(t *testing.T) {
	productInMemoryRepo := repository.NewProductRepositoryInMemory()
	ws := NewWebServer("8080", productInMemoryRepo)

	e := echo.New()
	rec := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	echoCtx := e.NewContext(req, rec)
	echoCtx.SetPath("products/:code")
	echoCtx.SetParamNames("code")
	echoCtx.SetParamValues("XSZ-000741")

	grApi := e.Group("/api/v1")
	grApi.DELETE("/products", ws.handleProductDelete)

	err := ws.handleProductDelete(echoCtx)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func deleteProductNotFoundErr(t *testing.T) {
	productInMemoryRepo := repository.ProductRepositoryInMemorySpy{
		ExpectedError: entity.ProductNotFoundErr,
	}
	ws := NewWebServer("8080", productInMemoryRepo)

	e := echo.New()
	rec := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	echoCtx := e.NewContext(req, rec)
	echoCtx.SetPath("products/:code")
	echoCtx.SetParamNames("code")
	echoCtx.SetParamValues("XSZ-000741")

	grApi := e.Group("/api/v1")
	grApi.DELETE("/products", ws.handleProductDelete)
	e.ServeHTTP(rec, req)

	err := ws.handleProductDelete(echoCtx)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestWebServer_handleProductUpdate(t *testing.T) {
	t.Run("Should handle update product request with success", updateProductSuccess)
	t.Run("Should results error if update product request send invalid payload", updateProductBindErr)
	t.Run("Should results error if product code does not exists", updateProductNotFoundErr)
}

func updateProductSuccess(t *testing.T) {
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

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/products", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	echoCtx := e.NewContext(req, rec)

	err2 := ws.handleProductUpdate(echoCtx)
	assert.Nil(t, err2)

	assert.Equal(t, "", rec.Body.String())
	assert.Equal(t, http.StatusOK, rec.Code)
}

func updateProductBindErr(t *testing.T) {
	productInMemoryRepo := repository.NewProductRepositoryInMemory()
	ws := NewWebServer("8080", productInMemoryRepo)
	e := echo.New()

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/products",
		bytes.NewReader([]byte(`{`)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	grApi := e.Group("/api")
	grApi.PATCH("/v1/products", ws.handleProductUpdate)
	e.ServeHTTP(rec, req)
	assert.Equal(t, "{\"message\":\"unexpected EOF\"}\n", rec.Body.String())
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func updateProductNotFoundErr(t *testing.T) {
	productInMemoryRepo := repository.ProductRepositoryInMemorySpy{
		ExpectedError: entity.ProductNotFoundErr,
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

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	echoCtx := e.NewContext(req, rec)
	echoCtx.SetPath("products")

	grApi := e.Group("/api/v1")
	grApi.PATCH("/products", ws.handleProductUpdate)
	e.ServeHTTP(rec, req)

	err2 := ws.handleProductUpdate(echoCtx)
	assert.NotNil(t, err2)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
