package flags

var flagTensorflowInstance *flagTensorflow

type flagTensorflow struct {
	// tensorflow
	tfservingModelVersion *int64
	tfservingTimeoutMs    *int64
}

// singleton instance
func init() {
	flagTensorflowInstance = new(flagTensorflow)
}

func getFlagTensorflowInstance() *flagTensorflow {
	return flagTensorflowInstance
}

// tfserving_model_version
func (s *flagTensorflow) setTfservingModelVersion(tfservingModelVersion *int64) {
	s.tfservingModelVersion = tfservingModelVersion
}

func (s *flagTensorflow) GetTfservingModelVersion() *int64 {
	return s.tfservingModelVersion
}

// tfserving_model_version
func (s *flagTensorflow) setTfservingTimeoutMs(tfservingTimeoutMs *int64) {
	s.tfservingTimeoutMs = tfservingTimeoutMs
}

func (s *flagTensorflow) GetTfservingTimeoutMs() *int64 {
	return s.tfservingTimeoutMs
}
