package entity

// Coupon represents a discount coupon
type Coupon struct {
	Code           string  `json:"code" binding:"required" example:"SUMMER2024"`
	Discount       int     `json:"discount" binding:"required" example:"10"`
	MinBasketValue float64 `json:"minBasketValue" binding:"required" example:"50.3"`
}
