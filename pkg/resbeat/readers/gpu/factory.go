package gpu

func NewGPUReader() (*GPUReader, error) {
	var reader GPUReader

	if err := reader.Init(); err != nil {
		return nil, err
	}

	return &reader, nil
}
