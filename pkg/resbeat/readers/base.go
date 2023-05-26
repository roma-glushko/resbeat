package readers

// StatsReader represents components that reads resource stats from different resource controllers
type StatsReader interface {
	GetMemoryUsageInBytes() (uint64, error)
	GetMemoryLimitInBytes() (uint64, error)
	GetCPUUsageInNanos() (uint64, error)
	GetCPUQuotaInMicros() (uint64, error)
	GetCPUPeriodInMicros() (uint64, error)
}

//func NewStatsReader() (*StatsReader, error) {
//	return NewCGroupV1Reader()
//}

type ByteSize uint64

const (
	_           = iota // ignore first value by assigning to blank identifier
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
)

type DummyStatsReader struct{}

func (r DummyStatsReader) GetMemoryUsageInBytes() (uint64, error) {
	return uint64(125 * MB), nil
}

func (r DummyStatsReader) GetMemoryLimitInBytes() (uint64, error) {
	return uint64(1 * GB), nil
}

func (r DummyStatsReader) GetCPUUsageInNanos() (uint64, error) {
	return 1, nil // TODO: make this number real
}

func (r DummyStatsReader) GetCPUQuotaInMicros() (uint64, error) {
	return 3, nil // TODO: make this number real
}

func (r DummyStatsReader) GetCPUPeriodInMicros() (uint64, error) {
	return 4, nil // TODO: make this number real
}
