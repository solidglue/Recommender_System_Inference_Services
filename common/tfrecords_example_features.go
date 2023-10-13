package common

type SeqExampleBuff struct {
	Key  *string
	Buff *[]byte
}

type ExampleFeatures struct {
	UserExampleFeatures        *SeqExampleBuff
	UserContextExampleFeatures *SeqExampleBuff
	ItemSeqExampleFeatures     *[]SeqExampleBuff
}
