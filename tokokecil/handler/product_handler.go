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

// GetAllProducts godoc
// @Summary      Get All Products
// @Description  Mendapatkan daftar seluruh produk (hanya untuk user yang sudah login/JWT protected)
// @Tags         Product
// @Produce      json
// @Success      200  {array}  dto.ProductResponse
// @Failure      401  {object} dto.ErrorResponse
// @Failure      500  {object} dto.ErrorResponse
// @Security     BearerAuth
// @Router       /products [get]
func (h *ProductHandler) GetAllProducts(c echo.Context) error {
	products, err := h.Service.GetAllProducts()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, products)
}

// @Summary      Create Product
// @Description  Menambahkan produk baru (butuh login/JWT)
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        data  body  dto.ProductRequest  true  "Product Data"
// @Success      201   {object} dto.ProductResponse
// @Failure      400   {object} dto.ErrorResponse
// @Security     BearerAuth
// @Router       /products [post]
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
