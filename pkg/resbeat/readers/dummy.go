package readers

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
}

func NewDummyStatsReader() DummyStatsReader {
	return DummyStatsReader{
		MemoryUsageInBytes: uint64(125 * MB),
		MemoryLimitInBytes: uint64(1 * GB),
	}
}

func (r DummyStatsReader) GetMemoryUsageInBytes() (uint64, error) {
	return r.MemoryUsageInBytes, nil
}

func (r DummyStatsReader) GetMemoryLimitInBytes() (uint64, error) {
	return r.MemoryLimitInBytes, nil
}

func (r DummyStatsReader) GetCPUUsageInNanos() (uint64, error) {
	return 150_000_000, nil // TODO: make this number real
}

func (r DummyStatsReader) GetCPUUsageLimitInCores() (float64, error) {
	return 3.0, nil
}
