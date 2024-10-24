package memdb

import (
	"reviewsch/internal/service/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	repo := New()
	assert.NotNil(t, repo)
	assert.NotNil(t, repo.entries)
}

func TestRepository_Save(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*Repository)
		coupon  entity.Coupon
		wantErr bool
	}{
		{
			name:  "save to nil map",
			setup: func(r *Repository) {},
			coupon: entity.Coupon{
				Code:     "TEST1",
				Discount: 10,
			},
			wantErr: false,
		},
		{
			name: "save to existing map",
			setup: func(r *Repository) {
				r.entries = map[string]entity.Coupon{
					"EXISTING": {Code: "EXISTING", Discount: 20},
				}
			},
			coupon: entity.Coupon{
				Code:     "TEST2",
				Discount: 15,
			},
			wantErr: false,
		},
		{
			name: "update existing coupon",
			setup: func(r *Repository) {
				r.entries = map[string]entity.Coupon{
					"TEST3": {Code: "TEST3", Discount: 20},
				}
			},
			coupon: entity.Coupon{
				Code:     "TEST3",
				Discount: 25,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := New()
			tt.setup(repo)

			err := repo.Save(tt.coupon)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// Initialize map if it was nil
			if repo.entries == nil {
				repo.entries = make(map[string]entity.Coupon)
			}

			// Verify the coupon was saved
			saved, exists := repo.entries[tt.coupon.Code]
			assert.True(t, exists, "coupon should be saved in the map")
			assert.Equal(t, tt.coupon, saved, "saved coupon should match input")
		})
	}
}

func TestRepository_FindByCode(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*Repository)
		code      string
		want      *entity.Coupon
		wantErr   bool
		errString string
	}{
		{
			name: "find existing coupon",
			setup: func(r *Repository) {
				r.entries = map[string]entity.Coupon{
					"TEST1": {Code: "TEST1", Discount: 10},
				}
			},
			code: "TEST1",
			want: &entity.Coupon{
				Code:     "TEST1",
				Discount: 10,
			},
			wantErr: false,
		},
		{
			name: "coupon not found",
			setup: func(r *Repository) {
				r.entries = map[string]entity.Coupon{
					"TEST1": {Code: "TEST1", Discount: 10},
				}
			},
			code:      "NONEXISTENT",
			want:      nil,
			wantErr:   true,
			errString: "coupon not found",
		},
		{
			name:      "nil map",
			setup:     func(r *Repository) {},
			code:      "TEST1",
			want:      nil,
			wantErr:   true,
			errString: "coupon not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := New()
			tt.setup(repo)

			got, err := repo.FindByCode(tt.code)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errString, err.Error())
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRepository_SaveAndFind(t *testing.T) {
	repo := New()
	repo.entries = make(map[string]entity.Coupon)

	// Test saving and then finding a coupon
	coupon := entity.Coupon{
		Code:     "TEST1",
		Discount: 10,
	}

	// Save the coupon
	err := repo.Save(coupon)
	assert.NoError(t, err)

	// Find the saved coupon
	found, err := repo.FindByCode(coupon.Code)
	assert.NoError(t, err)
	assert.Equal(t, &coupon, found)
}
