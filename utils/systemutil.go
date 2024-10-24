package utils

import (
	"fmt"
	"runtime"
)

type SystemRequirements struct {
	MinCPUCores int
}

func DefaultSystemRequirements() SystemRequirements {
	return SystemRequirements{
		MinCPUCores: 4,
	}
}

func getNumCPU() int {
	return runtime.NumCPU()
}

func (sr SystemRequirements) Validate() error {
	actualCPUs := getNumCPU()
	if actualCPUs < sr.MinCPUCores {
		return fmt.Errorf("insufficient CPU cores: got %d, want minimum %d", actualCPUs, sr.MinCPUCores)
	}
	return nil
}
