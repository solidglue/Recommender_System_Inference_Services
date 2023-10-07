package apis

type recRequest struct {
	dataId      string //nacos dataid
	groupId     string
	namespaceId string
	modelType   string //recall or rank
	userId      string
	recallNum   int32    //recall num
	itemList    []string //rank items
}

// dataId
func (r *recRequest) SetDataId(dataId string) {
	r.dataId = dataId
}

func (r *recRequest) GetDataId() string {
	return r.dataId
}

// groupId
func (r *recRequest) SetGroupId(groupId string) {
	r.groupId = groupId
}

func (r *recRequest) GetGroupId() string {
	return r.groupId
}

// namespaceId
func (r *recRequest) SetNamespaceId(namespaceId string) {
	r.namespaceId = namespaceId
}

func (r *recRequest) GetNamespaceId() string {
	return r.namespaceId
}

// modelType
func (r *recRequest) SetModelType(modelType string) {
	r.modelType = modelType
}

func (r *recRequest) GetModelType() string {
	return r.modelType
}

// userId
func (r *recRequest) SetUserId(userId string) {
	r.userId = userId
}

func (r *recRequest) GetUserId() string {
	return r.userId
}

// recallNum
func (r *recRequest) SetRecallNum(recallNum int32) {
	r.recallNum = recallNum
}

func (r *recRequest) GetRecallNum() int32 {
	return r.recallNum
}

// itemList
func (r *recRequest) SetItemList(itemList []string) {
	r.itemList = itemList
}

func (r *recRequest) GetItemList() []string {
	return r.itemList
}

func (req *recRequest) JavaClassName() string {
	return "com.xxx.www.infer.RecRequest"
}
