package resbeat

type CPUStats struct {
	LimitInCors     uint64
	UsagePercentage float32
}

type MemoryStats struct {
	UsagePercentage float32
	LimitInBytes    uint64
	UsageInBytes    uint64
}

type Usage struct {
	CPU    *CPUStats
	Memory *MemoryStats
}
