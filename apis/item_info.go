package apis

type itemInfo struct {
	itemId string
	score  float32
}

// itemId
func (i *itemInfo) SetItemId(itemId string) {
	i.itemId = itemId
}

func (i *itemInfo) GetItemId() string {
	return i.itemId
}

// score
func (i *itemInfo) SetScore(score float32) {
	i.score = score
}

func (i *itemInfo) GetScore() float32 {
	return i.score
}
