package apis

type RecResponse struct {
	code    int //请求状态
	message string
	data    []string //物品信息
}

// code
func (r *RecResponse) SetCode(code int) {
	r.code = code
}

func (r *RecResponse) GetCode() int {
	return r.code
}

// message
func (r *RecResponse) SetMessage(message string) {
	r.message = message
}

func (r *RecResponse) GetMessage() string {
	return r.message
}

// data
func (r *RecResponse) SetData(data []string) {
	r.data = data
}

func (r *RecResponse) GetData() []string {
	return r.data
}

func (rsp *RecResponse) JavaClassName() string {
	return "com.xxx.www.infer.RecResponse"
}
