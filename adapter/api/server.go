package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lbsti/eulabs-challenge/core/repository"
	"github.com/lbsti/eulabs-challenge/core/usecase"
)

type WebServer struct {
	productRepo repository.ProductRepository
	port        string
}

func NewWebServer(port string, productRepo repository.ProductRepository) *WebServer {
	return &WebServer{productRepo: productRepo, port: port}
}

func (ws WebServer) Run() {
	e := echo.New()
	productGroup := e.Group("/api")
	productGroup.POST("/v1/products", ws.handleProductCreate)
	productGroup.GET("/v1/products/:code", ws.handleProductGet)
	productGroup.DELETE("/v1/products/:code", ws.handleProductDelete)
	productGroup.PATCH("/v1/products", ws.handleProductUpdate)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", ws.port)))
}

func (ws WebServer) handleProductCreate(c echo.Context) error {
	productCreate := usecase.NewProductCreate(ws.productRepo)
	var inputDTO usecase.ProductInputDTO
	if err := c.Bind(&inputDTO); err != nil {
		return err
	}
	ctx := c.Request().Context()
	outputDTO, err := productCreate.Execute(ctx, inputDTO)
	if err != nil {
		code := Mapping(err).Code
		e := Mapping(err)
		return echo.NewHTTPError(code, e.ResultErr.Error())
	}
	return c.JSON(http.StatusCreated, outputDTO)
}

func (ws WebServer) handleProductGet(c echo.Context) error {
	code := c.Param("code")
	productGet := usecase.NewProductGet(ws.productRepo)
	ctx := c.Request().Context()
	outputDTO, err := productGet.Execute(ctx, code)
	if err != nil {
		code := Mapping(err).Code
		e := Mapping(err)
		return echo.NewHTTPError(code, e.ResultErr.Error())
	}
	return c.JSON(http.StatusOK, outputDTO)
}

func (ws WebServer) handleProductDelete(c echo.Context) error {
	code := c.Param("code")
	productDelete := usecase.NewProductDelete(ws.productRepo)
	ctx := c.Request().Context()
	_, err := productDelete.Execute(ctx, code)
	if err != nil {
		code := Mapping(err).Code
		e := Mapping(err)
		return echo.NewHTTPError(code, e.ResultErr.Error())
	}
	return c.NoContent(http.StatusOK)
}

func (ws WebServer) handleProductUpdate(c echo.Context) error {
	productUpdate := usecase.NewProductUpdate(ws.productRepo)

	var inputDTO usecase.ProductInputDTO
	if err := c.Bind(&inputDTO); err != nil {
		return err
	}
	ctx := c.Request().Context()
	if err := productUpdate.Execute(ctx, inputDTO); err != nil {
		code := Mapping(err).Code
		e := Mapping(err)
		return echo.NewHTTPError(code, e.ResultErr.Error())
	}
	return c.NoContent(http.StatusOK)
}
