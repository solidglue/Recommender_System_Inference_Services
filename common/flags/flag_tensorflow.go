package flags

import "flag"

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
func (s *flagTensorflow) setTfservingModelVersion() {
	conf := flag.Int64("tfserving_model_version", 0, "")
	s.tfservingModelVersion = conf
}

func (s *flagTensorflow) GetTfservingModelVersion() *int64 {
	return s.tfservingModelVersion
}

// tfserving_model_version
func (s *flagTensorflow) setTfservingTimeoutMs() {
	conf := flag.Int64("tfserving_timeoutms", 100, "")
	s.tfservingTimeoutMs = conf
}

func (s *flagTensorflow) GetTfservingTimeoutMs() *int64 {
	return s.tfservingTimeoutMs
}
