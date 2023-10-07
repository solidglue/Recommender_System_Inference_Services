package common

type SeqExampleBuff struct { // 特征处理的输出
	Key  *string //  用户id或物品id
	Buff *[]byte // 输入模型的最小单元
}

type ExampleFeatures struct {
	UserExampleFeatures        *SeqExampleBuff
	UserContextExampleFeatures *SeqExampleBuff
	ItemSeqExampleFeatures     *[]SeqExampleBuff
}
