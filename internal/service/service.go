package service

import (
	"fmt"
	. "reviewsch/internal/service/entity"

	"github.com/google/uuid"
)

type Repository interface {
	FindByCode(string) (*Coupon, error)
	Save(Coupon) error
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) ApplyCoupon(basket Basket, code string) (*Basket, error) {
	if code == "" {
		return nil, fmt.Errorf("empty coupon code")
	}

	coupon, err := s.repo.FindByCode(code)
	if err != nil {
		return nil, err
	}

	result := &basket
	if result.Value <= 0 {
		return nil, fmt.Errorf("invalid basket value")
	}

	result.AppliedDiscount = coupon.Discount
	result.ApplicationSuccessful = true
	result.CouponCode = code

	return result, nil
}

func (s *Service) CreateCoupon(discount int, code string, minBasketValue float64) error {
	if code == "" {
		return fmt.Errorf("empty coupon code")
	}

	coupon := Coupon{
		ID:             uuid.NewString(),
		Discount:       discount,
		Code:           code,
		MinBasketValue: minBasketValue,
	}

	return s.repo.Save(coupon)
}

func (s *Service) GetCoupons(codes []string) ([]Coupon, error) {
	coupons := make([]Coupon, 0, len(codes))

	for _, code := range codes {
		coupon, err := s.repo.FindByCode(code)
		if err != nil {
			return nil, fmt.Errorf("error finding coupon %s: %w", code, err)
		}
		coupons = append(coupons, *coupon)
	}

	return coupons, nil
}
