package system

type ByteSize uint64

const (
	_           = iota // ignore first value by assigning to blank identifier
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
)

type DummyStatsReader struct {
	memoryUsageInBytes uint64
	memoryLimitInBytes uint64
	cpuLimitInCores    float64
	cpuUsageInNano     uint64
}

func NewDummyStatsReader(memoryUsageInBytes, memoryLimitInBytes *uint64, cpuUsageInNano *uint64, limitInCores *float64) *DummyStatsReader {
	defMemUsageInBytes, defMemLimInBytes := uint64(125*MB), uint64(1*GB)
	defCPUUsageInNano, defLimInCores := uint64(149537069), 2.0

	if memoryUsageInBytes == nil {
		memoryUsageInBytes = &defMemUsageInBytes
	}

	if memoryLimitInBytes == nil {
		memoryLimitInBytes = &defMemLimInBytes
	}

	if cpuUsageInNano == nil {
		cpuUsageInNano = &defCPUUsageInNano
	}

	if limitInCores == nil {
		limitInCores = &defLimInCores
	}

	return &DummyStatsReader{
		memoryUsageInBytes: *memoryUsageInBytes,
		memoryLimitInBytes: *memoryLimitInBytes,
		cpuLimitInCores:    *limitInCores,
		cpuUsageInNano:     *cpuUsageInNano,
	}
}

func (r DummyStatsReader) MemoryUsageInBytes() (uint64, error) {
	return r.memoryUsageInBytes, nil
}

func (r *DummyStatsReader) SetMemoryUsageInBytes(usage uint64) {
	r.memoryUsageInBytes = usage
}

func (r DummyStatsReader) MemoryLimitInBytes() (uint64, error) {
	return r.memoryLimitInBytes, nil
}

func (r DummyStatsReader) CPUUsageInNanos() (uint64, error) {
	return r.cpuUsageInNano, nil
}

func (r *DummyStatsReader) SetCPUUsageInNanos(usage uint64) {
	r.cpuUsageInNano = usage
}

func (r DummyStatsReader) CPUUsageLimitInCores() (float64, error) {
	return r.cpuLimitInCores, nil
}
