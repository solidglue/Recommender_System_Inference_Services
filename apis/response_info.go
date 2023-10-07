package apis

type recResponse struct {
	code    int
	message string
	data    []string
}

// code
func (r *recResponse) SetCode(code int) {
	r.code = code
}

func (r *recResponse) GetCode() int {
	return r.code
}

// message
func (r *recResponse) SetMessage(message string) {
	r.message = message
}

func (r *recResponse) GetMessage() string {
	return r.message
}

// data
func (r *recResponse) SetData(data []string) {
	r.data = data
}

func (r *recResponse) GetData() []string {
	return r.data
}

func (rsp *recResponse) JavaClassName() string {
	return "com.xxx.www.infer.RecResponse"
}
