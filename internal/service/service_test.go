package service

import (
	"fmt"
	. "reviewsch/internal/service/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockRepository is a mock implementation of Repository interface
type mockRepository struct {
	coupons map[string]*Coupon
	err     error
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		coupons: make(map[string]*Coupon),
	}
}

func (m *mockRepository) FindByCode(code string) (*Coupon, error) {
	if m.err != nil {
		return nil, m.err
	}
	coupon, exists := m.coupons[code]
	if !exists {
		return nil, fmt.Errorf("coupon not found")
	}
	return coupon, nil
}

func (m *mockRepository) Save(coupon Coupon) error {
	if m.err != nil {
		return m.err
	}
	m.coupons[coupon.Code] = &coupon
	return nil
}

func TestService_ApplyCoupon(t *testing.T) {
	tests := []struct {
		name         string
		basket       Basket
		code         string
		setupRepo    func(*mockRepository)
		expectedErr  string
		expectBasket *Basket
	}{
		{
			name: "successful coupon application",
			basket: Basket{
				Value: 100,
			},
			code: "TEST10",
			setupRepo: func(m *mockRepository) {
				m.coupons["TEST10"] = &Coupon{
					Code:     "TEST10",
					Discount: 10,
				}
			},
			expectBasket: &Basket{
				Value:                 100,
				AppliedDiscount:       10,
				ApplicationSuccessful: true,
				CouponCode:            "TEST10",
			},
		},
		{
			name: "empty coupon code",
			basket: Basket{
				Value: 100,
			},
			code:        "",
			setupRepo:   func(m *mockRepository) {},
			expectedErr: "empty coupon code",
		},
		{
			name: "invalid basket value",
			basket: Basket{
				Value: 0,
			},
			code: "TEST10",
			setupRepo: func(m *mockRepository) {
				m.coupons["TEST10"] = &Coupon{
					Code:     "TEST10",
					Discount: 10,
				}
			},
			expectedErr: "invalid basket value",
		},
		{
			name: "coupon not found",
			basket: Basket{
				Value: 100,
			},
			code:        "INVALID",
			setupRepo:   func(m *mockRepository) {},
			expectedErr: "coupon not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockRepository()
			tt.setupRepo(repo)

			service := New(repo)
			result, err := service.ApplyCoupon(tt.basket, tt.code)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, result)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectBasket, result)
		})
	}
}

func TestService_CreateCoupon(t *testing.T) {
	tests := []struct {
		name           string
		discount       int
		code           string
		minBasketValue float64
		setupRepo      func(*mockRepository)
		expectedErr    string
	}{
		{
			name:           "successful coupon creation",
			discount:       10,
			code:           "NEW10",
			minBasketValue: 50,
			setupRepo:      func(m *mockRepository) {},
		},
		{
			name:           "empty code",
			discount:       10,
			code:           "",
			minBasketValue: 50,
			setupRepo:      func(m *mockRepository) {},
			expectedErr:    "empty coupon code",
		},
		{
			name:           "repository error",
			discount:       10,
			code:           "ERROR10",
			minBasketValue: 50,
			setupRepo: func(m *mockRepository) {
				m.err = fmt.Errorf("database error")
			},
			expectedErr: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockRepository()
			tt.setupRepo(repo)

			service := New(repo)
			err := service.CreateCoupon(tt.discount, tt.code, tt.minBasketValue)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				return
			}

			assert.NoError(t, err)
			// Verify coupon was saved
			saved, err := repo.FindByCode(tt.code)
			assert.NoError(t, err)
			assert.Equal(t, tt.discount, saved.Discount)
			assert.Equal(t, tt.minBasketValue, saved.MinBasketValue)
		})
	}
}

func TestService_GetCoupons(t *testing.T) {
	tests := []struct {
		name        string
		codes       []string
		setupRepo   func(*mockRepository)
		expectedErr string
		expectCount int
	}{
		{
			name:  "successful multiple coupons retrieval",
			codes: []string{"CODE1", "CODE2"},
			setupRepo: func(m *mockRepository) {
				m.coupons["CODE1"] = &Coupon{Code: "CODE1", Discount: 10}
				m.coupons["CODE2"] = &Coupon{Code: "CODE2", Discount: 20}
			},
			expectCount: 2,
		},
		{
			name:  "one invalid code",
			codes: []string{"CODE1", "INVALID"},
			setupRepo: func(m *mockRepository) {
				m.coupons["CODE1"] = &Coupon{Code: "CODE1", Discount: 10}
			},
			expectedErr: "error finding coupon INVALID",
		},
		{
			name:        "empty codes list",
			codes:       []string{},
			setupRepo:   func(m *mockRepository) {},
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockRepository()
			tt.setupRepo(repo)

			service := New(repo)
			coupons, err := service.GetCoupons(tt.codes)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, coupons)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, coupons, tt.expectCount)
			for i, code := range tt.codes {
				assert.Equal(t, code, coupons[i].Code)
			}
		})
	}
}
