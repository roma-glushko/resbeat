package system

type ByteSize uint64

const (
	_           = iota // ignore first value by assigning to blank identifier
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
)

type DummyStatsReader struct {
	MemoryUsageInBytes uint64
	MemoryLimitInBytes uint64
	CPULimitInCores    float64
	CPUUsageInNano     uint64
}

func NewDummyStatsReader() DummyStatsReader {
	return DummyStatsReader{
		MemoryUsageInBytes: uint64(125 * MB),
		MemoryLimitInBytes: uint64(1 * GB),
		CPULimitInCores:    2.0,
		CPUUsageInNano:     uint64(149537069),
	}
}

func (r *DummyStatsReader) GetMemoryUsageInBytes() (uint64, error) {
	return r.MemoryUsageInBytes, nil
}

func (r *DummyStatsReader) GetMemoryLimitInBytes() (uint64, error) {
	return r.MemoryLimitInBytes, nil
}

func (r *DummyStatsReader) GetCPUUsageInNanos() (uint64, error) {
	return r.CPUUsageInNano, nil
}

func (r *DummyStatsReader) GetCPUUsageLimitInCores() (float64, error) {
	return r.CPULimitInCores, nil
}
