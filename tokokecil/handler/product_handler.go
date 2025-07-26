package handler

import (
	"net/http"
	"tokokecil/model"
	"tokokecil/service"

	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	Service service.ProductService
}

func NewProductHandler(s service.ProductService) *ProductHandler {
	return &ProductHandler{Service: s}
}

func (h *ProductHandler) GetAllProducts(c echo.Context) error {
	products, err := h.Service.GetAllProducts()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) CreateProduct(c echo.Context) error {
	var p model.Product
	if err := c.Bind(&p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}
	product, err := h.Service.CreateProduct(p)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, product)
}
