package apis


type RecRequest struct {
	dataId      string //模型id，场景-模型            //必填
	groupId     string
	namespaceId string
	modelType   string   //召回，排序
	userId      string   //用户id                      //必填
	recallNum   int32    //返回条数                     //返回条数，粗排/召回。         注册传入
	itemList    []string //item列表                    //必填
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
	return "com.xxx.www.infer.RecRequest" // 如果与 Java 互通，需要与 Java 侧 User class全名对应,
}
