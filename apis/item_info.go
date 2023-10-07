package apis

type ItemInfo struct {
	itemId string
	score  float32
}

// itemId
func (i *ItemInfo) SetItemId(itemId string) {
	i.itemId = itemId
}

func (i *ItemInfo) GetItemId() string {
	return i.itemId
}

// score
func (i *ItemInfo) SetScore(score float32) {
	i.score = score
}

func (i *ItemInfo) GetScore() float32 {
	return i.score
}
