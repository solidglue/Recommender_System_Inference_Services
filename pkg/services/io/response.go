package io

type RecResponse struct {
	code    int
	message string
	data    []string //ItemInfo string
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
	return "com.loki.www.infer.RecResponse"
}
