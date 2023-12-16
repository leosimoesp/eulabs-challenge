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
