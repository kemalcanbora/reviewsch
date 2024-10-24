package router

import (
	"fmt"
	"net/http"
	. "reviewsch/internal/api/dto/entity"
	"reviewsch/internal/api/handler"

	"github.com/gin-gonic/gin"
)

// CouponHandler handles coupon-related operations
type CouponHandler struct {
	svc handler.Service
}

// NewCouponHandler creates a new CouponHandler instance
func NewCouponHandler(svc handler.Service) *CouponHandler {
	return &CouponHandler{
		svc: svc,
	}
}

// Apply godoc
// @Summary Apply a coupon to a basket
// @Description Apply a coupon to a basket
// @Tags Coupons
// @Produce json
// @Success 200 {object} entity.Basket
// @Router /v1/coupons/apply [get]
// @Security Bearer
// @Param Authorization header string true "Bearer JWT token"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse "Unauthorized"
func (h *CouponHandler) Apply(c *gin.Context) {
	apiReq := ApplicationRequest{}
	if err := c.ShouldBindJSON(&apiReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	basket, err := h.svc.ApplyCoupon(apiReq.Basket, apiReq.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, basket)
}

// Create godoc
// @Summary Create a new coupon
// @Description Create a new coupon
// @Tags Coupons
// @Produce json
// @Success 200 {object} SuccessResponse
// @Router /v1/coupons/create [post]
// @Security Bearer
// @Param Authorization header string true "Bearer JWT token"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse "Unauthorized"
func (h *CouponHandler) Create(c *gin.Context) {
	apiReq := Coupon{}
	fmt.Println("lol")
	if err := c.ShouldBindJSON(&apiReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format: " + err.Error(),
		})
		return
	}

	if err := h.svc.CreateCoupon(apiReq.Discount, apiReq.Code, apiReq.MinBasketValue); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Coupon created successfully",
	})
}

// Get godoc
// @Summary Get coupons by codes
// @Description Retrieve multiple coupons by their codes
// @Tags Coupons
// @Accept json
// @Produce json
// @Success 200 {objects} entity.Coupon
// @Router /v1/coupons [get]
// @Security Bearer
// @Param Authorization header string true "Bearer JWT token"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse "Unauthorized"
func (h *CouponHandler) Get(c *gin.Context) {
	apiReq := CouponRequest{}
	if err := c.ShouldBindJSON(&apiReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	coupons, err := h.svc.GetCoupons(apiReq.Codes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, coupons)
}

// ErrorResponse Define response structures for Swagger
type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}
