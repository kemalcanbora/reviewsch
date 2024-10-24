package utils

import (
	"testing"
)

func TestDefaultSystemRequirements(t *testing.T) {
	sr := DefaultSystemRequirements()
	if sr.MinCPUCores != 4 {
		t.Errorf("DefaultSystemRequirements() = %v, want MinCPUCores = 4", sr)
	}
}

func TestSystemRequirements_Validate(t *testing.T) {
	tests := []struct {
		name       string
		sr         SystemRequirements
		wantError  bool
		errorMatch string
	}{
		{
			name:       "sufficient cores",
			sr:         SystemRequirements{MinCPUCores: 2},
			wantError:  false,
			errorMatch: "",
		},
		{
			name:       "exact cores",
			sr:         SystemRequirements{MinCPUCores: 4},
			wantError:  false,
			errorMatch: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sr.Validate()
			if (err != nil) != tt.wantError {
				t.Errorf("Validate() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if tt.wantError && err.Error() != tt.errorMatch {
				t.Errorf("Validate() error = %v, want %v", err, tt.errorMatch)
			}
		})
	}
}
