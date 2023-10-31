package io

type RecRequest struct {
	dataId      string //nacos dataid
	groupId     string
	namespaceId string
	modelType   string //recall or rank
	userId      string
	recallNum   int32    //recall num
	itemList    []string //rank items
}

// dataId
func (r *RecRequest) SetDataId(dataId string) {
	r.dataId = dataId
}

func (r *RecRequest) GetDataId() string {
	return r.dataId
}

// groupId
func (r *RecRequest) SetGroupId(groupId string) {
	r.groupId = groupId
}

func (r *RecRequest) GetGroupId() string {
	return r.groupId
}

// namespaceId
func (r *RecRequest) SetNamespaceId(namespaceId string) {
	r.namespaceId = namespaceId
}

func (r *RecRequest) GetNamespaceId() string {
	return r.namespaceId
}

// modelType
func (r *RecRequest) SetModelType(modelType string) {
	r.modelType = modelType
}

func (r *RecRequest) GetModelType() string {
	return r.modelType
}

// userId
func (r *RecRequest) SetUserId(userId string) {
	r.userId = userId
}

func (r *RecRequest) GetUserId() string {
	return r.userId
}

// recallNum
func (r *RecRequest) SetRecallNum(recallNum int32) {
	r.recallNum = recallNum
}

func (r *RecRequest) GetRecallNum() int32 {
	return r.recallNum
}

// itemList
func (r *RecRequest) SetItemList(itemList []string) {
	r.itemList = itemList
}

func (r *RecRequest) GetItemList() []string {
	return r.itemList
}

func (req *RecRequest) JavaClassName() string {
	return "com.xxx.www.infer.RecRequest"
}
