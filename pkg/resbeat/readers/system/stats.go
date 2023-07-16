package system

import "time"

type CPUStats struct {
	collectedAt             time.Time
	accumulatedUsageInNanos uint64
	UsageInNanos            uint64  `json:"usageInNanos"`
	LimitInCores            float64 `json:"limitInCors"`
	UsagePercentage         float64 `json:"usagePercentage"`
}

func (s *CPUStats) CollectedAt() time.Time {
	return s.collectedAt
}

func (s *CPUStats) AccumulatedUsageInNanos() uint64 {
	return s.accumulatedUsageInNanos
}

type MemoryStats struct {
	UsagePercentage float64 `json:"usagePercentage"`
	LimitInBytes    uint64  `json:"limitInBytes"`
	UsageInBytes    uint64  `json:"usageInBytes"`
}

type SystemStats struct {
	CPU    *CPUStats    `json:"cpu,omitempty"`
	Memory *MemoryStats `json:"memory"`
}
