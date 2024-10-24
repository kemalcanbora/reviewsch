package entity

import "reviewsch/internal/service/entity"

// ApplicationRequest Request/Response Models for Swagger documentation
// ApplicationRequest represents the request for applying a coupon
type ApplicationRequest struct {
	Basket entity.Basket `json:"basket"`
	Code   string        `json:"code" example:"SUMMER2024"`
}
