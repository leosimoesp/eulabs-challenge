package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lbsti/eulabs-challenge/internal/core/repository"
	"github.com/lbsti/eulabs-challenge/internal/core/usecase"
)

type WebServer struct {
	productRepo repository.ProductRepository
	port        string
}

func NewWebServer(port string, productRepo repository.ProductRepository) *WebServer {
	return &WebServer{productRepo: productRepo, port: port}
}

func (ws WebServer) Run() {
	echoInstance := echo.New()
	productGroup := echoInstance.Group("/api")
	productGroup.POST("/v1/products", ws.handleProductCreate)
	productGroup.GET("/v1/products/:code", ws.handleProductGet)
	productGroup.DELETE("/v1/products/:code", ws.handleProductDelete)
	productGroup.PATCH("/v1/products", ws.handleProductUpdate)
	echoInstance.Logger.Fatal(echoInstance.Start(fmt.Sprintf(":%s", ws.port)))
}

func (ws WebServer) handleProductCreate(echoCtx echo.Context) error {
	productCreate := usecase.NewProductCreate(ws.productRepo)
	var inputDTO usecase.ProductInputDTO
	if err := echoCtx.Bind(&inputDTO); err != nil {
		return err
	}
	ctx := echoCtx.Request().Context()
	outputDTO, err := productCreate.Execute(ctx, inputDTO)
	if err != nil {
		code := Mapping(err).Code
		wrappedErr := Mapping(err)
		return echo.NewHTTPError(code, wrappedErr.ResultErr.Error())
	}
	return echoCtx.JSON(http.StatusCreated, outputDTO)
}

func (ws WebServer) handleProductGet(echoCtx echo.Context) error {
	code := echoCtx.Param("code")
	productGet := usecase.NewProductGet(ws.productRepo)
	ctx := echoCtx.Request().Context()
	outputDTO, err := productGet.Execute(ctx, code)
	if err != nil {
		code := Mapping(err).Code
		wrappedErr := Mapping(err)
		return echo.NewHTTPError(code, wrappedErr.ResultErr.Error())
	}
	return echoCtx.JSON(http.StatusOK, outputDTO)
}

func (ws WebServer) handleProductDelete(echoCtx echo.Context) error {
	code := echoCtx.Param("code")
	productDelete := usecase.NewProductDelete(ws.productRepo)
	ctx := echoCtx.Request().Context()
	_, err := productDelete.Execute(ctx, code)
	if err != nil {
		code := Mapping(err).Code
		wrappedErr := Mapping(err)
		return echo.NewHTTPError(code, wrappedErr.ResultErr.Error())
	}
	return echoCtx.NoContent(http.StatusOK)
}

func (ws WebServer) handleProductUpdate(echoCtx echo.Context) error {
	productUpdate := usecase.NewProductUpdate(ws.productRepo)

	var inputDTO usecase.ProductInputDTO
	if err := echoCtx.Bind(&inputDTO); err != nil {
		return err
	}
	ctx := echoCtx.Request().Context()
	if err := productUpdate.Execute(ctx, inputDTO); err != nil {
		code := Mapping(err).Code
		wrappedErr := Mapping(err)
		return echo.NewHTTPError(code, wrappedErr.ResultErr.Error())
	}
	return echoCtx.NoContent(http.StatusOK)
}
