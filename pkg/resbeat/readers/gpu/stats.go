package gpu

type AllGPUStats = map[string]GPUStats

type GPUStats struct {
	UsagePercentage    uint32
	MemoryUsedInBytes  uint64
	TotalMemoryInBytes uint64
}
