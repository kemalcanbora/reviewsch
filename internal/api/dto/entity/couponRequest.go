package entity

// CouponRequest represents a request to get coupons
// @Description Request to retrieve multiple coupons
type CouponRequest struct {
	Codes []string `json:"codes" example:"['SUMMER2024', 'WINTER2024']"`
}
