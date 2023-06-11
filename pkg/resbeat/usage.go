package resbeat

import "time"

type CPUStats struct {
	UsageInNanos    uint64  `json:"usageInNanos"`
	LimitInCores    float64 `json:"limitInCors"`
	UsagePercentage float64 `json:"usagePercentage"`
}

type MemoryStats struct {
	UsagePercentage float64 `json:"usagePercentage"`
	LimitInBytes    uint64  `json:"limitInBytes"`
	UsageInBytes    uint64  `json:"usageInBytes"`
}

type Usage struct {
	CollectedAt time.Time    `json:"collectedAt"`
	CPU         *CPUStats    `json:"cpu,omitempty"`
	Memory      *MemoryStats `json:"memory"`
}
