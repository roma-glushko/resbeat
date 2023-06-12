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

type SystemStats struct {
	CPU    *CPUStats    `json:"cpu,omitempty"`
	Memory *MemoryStats `json:"memory"`
}

type GPUStats struct {
	// TODO: implement
}

type DiskStats struct {
	// TODO: implement
}

type Usage struct {
	CollectedAt time.Time             `json:"collectedAt"`
	System      *SystemStats          `json:"system"`
	GPUs        *map[string]GPUStats  `json:"gpus,omitempty"`
	Disks       *map[string]DiskStats `json:"disks,omitempty"`
}
