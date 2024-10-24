package entity

// Coupon represents a discount coupon
// @Description Discount coupon
type Coupon struct {
	ID             string
	Code           string
	Discount       int
	MinBasketValue float64
}
