package resbeat

type CPUStats struct {
	LimitInCors     uint64  `json:"limitInCors"`
	UsagePercentage float64 `json:"usagePercentage"`
}

type MemoryStats struct {
	UsagePercentage float64 `json:"usagePercentage"`
	LimitInBytes    uint64  `json:"limitInBytes"`
	UsageInBytes    uint64  `json:"usageInBytes"`
}

type Usage struct {
	CPU    *CPUStats    `json:"cpu,omitempty"`
	Memory *MemoryStats `json:"memory"`
}
