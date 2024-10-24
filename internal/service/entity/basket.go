package entity

import (
	_ "github.com/gin-gonic/gin"
)

// Basket represents a shopping basket
// @Description Shopping basket with coupon application details
type Basket struct {
	Value                 float64 `json:"value" example:"100.50"`
	AppliedDiscount       int     `json:"appliedDiscount" example:"10"`
	ApplicationSuccessful bool    `json:"applicationSuccessful" example:"true"`
	CouponCode            string  `json:"couponCode" example:"SUMMER2024"`
}
