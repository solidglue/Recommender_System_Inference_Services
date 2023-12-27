package flags

var flagKafkaInstance *flagKafka

type flagKafka struct {
	// tensorflow
	kafkaUrl   string
	kafkaTopic string
	kafkaGroup string
}

// singleton instance
func init() {
	flagKafkaInstance = new(flagKafka)
}

func getFlagKafkaInstance() *flagKafka {
	return flagKafkaInstance
}

// kafkaUrl
func (s *flagKafka) setKafkaUrl(kafkaUrl string) {
	s.kafkaUrl = kafkaUrl
}

func (s *flagKafka) GetKafkaUrl() string {
	return s.kafkaUrl
}

// kafkaTopic
func (s *flagKafka) setKafkaTopic(kafkaTopic string) {
	s.kafkaTopic = kafkaTopic
}

func (s *flagKafka) GetKafkaTopic() string {
	return s.kafkaTopic
}

// kafkaGroup
func (s *flagKafka) setKafkaGroup(kafkaGroup string) {
	s.kafkaGroup = kafkaGroup
}

func (s *flagKafka) GetKafkaGroup() string {
	return s.kafkaGroup
}
