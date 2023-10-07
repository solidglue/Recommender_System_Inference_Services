package server

import (
	"errors"
	"infer-microservices/apis"
	"infer-microservices/utils/logs"
	"net/http"
	"reflect"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

func httpRequstParse(r *http.Request) (apis.RecRequest, error) {

	request := apis.RecRequest{}

	err := r.ParseForm()
	if err != nil {
		return request, err
	}

	method := r.Method
	if method != "POST" {
		return request, err
	}

	data := r.Form["data"]
	if len(data) == 0 {
		return request, err

	}

	requestMap := make(map[string]interface{}, 0)
	err = jsoniter.Unmarshal([]byte(data[0]), &requestMap)
	if err != nil {
		return request, err
	}

	request, err = inputCheck(requestMap)
	if err != nil {
		return request, err
	}

	return request, nil
}

func inputCheck(requestMap map[string]interface{}) (apis.RecRequest, error) {

	//TODO: 传来的参数不固定，且枚举太多，考虑反射(性能差，慎用) https://blog.csdn.net/DkSakura/article/details/116588382
	//TODO:反射应用点 - 如果api请求不规范，纠正一下？例如int类型传成了string类型，float类型与float64，
	//，而不是直接拒绝，尤其是http请求时，用户手动输入不规范场景。或者是不同系统传过来的类型不一样（32位，64位os）

	request := apis.RecRequest{}

	//dataId
	dataId, ok := requestMap["dataId"]
	if ok {
		request.SetDataId(dataId.(string))
	} else {
		return request, errors.New("dataId can not be empty")
	}

	//modelType
	modelType, ok := requestMap["modelType"]
	if ok {
		request.SetModelType(modelType.(string))
	} else {
		return request, errors.New("modelType can not be empty")
	}

	//userId
	userId, ok := requestMap["userId"]
	if ok {
		request.SetUserId(userId.(string))
	} else {
		return request, errors.New("userId can not be empty")
	}

	// //recallNum
	// recallNum, ok := requestMap["recallNum"].(int32)
	// if ok {
	// 	request.SetRecallNum(recallNum)
	// } else {
	// 	return request, errors.New("dataId can not be empty")
	// }

	// if recallNum > 1000 {
	// 	return request, errors.New("recallNum should less than 1000 ")
	// }

	//recallNum

	recallNumType := reflect.TypeOf(requestMap["recallNum"])
	//recallNumValue := reflect.ValueOf(requestMap["recallNum"])

	recallNum := int32(100)
	recallNumTypeKind := recallNumType.Kind()
	//var ok bool

	switch recallNumTypeKind {
	case reflect.String:
		recallNumStr, ok0 := requestMap["recallNum"].(string)
		if ok0 {
			recallNum64, err := strconv.ParseInt(recallNumStr, 10, 64)
			if err != nil {
				ok = false
			} else {
				recallNum = int32(recallNum64)
				ok = true
			}
		}

	case reflect.Float32, reflect.Float64, reflect.Int16, reflect.Int, reflect.Int64, reflect.Int8,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		recallNum, ok = requestMap["recallNum"].(int32)
	default:
		err := errors.New("unkown type, set recallnum to 100")
		logs.Error(err)
	}

	//recallNum, ok := requestMap["recallNum"].(int32)
	if ok {
		request.SetRecallNum(recallNum)
	} else {
		return request, errors.New("dataId can not be empty")
	}

	if recallNum > 1000 {
		return request, errors.New("recallNum should less than 1000 ")
	}

	//itemList
	itemList, ok := requestMap["itemList"].([]string)
	if ok {
		request.SetItemList(itemList)
	} else {
		return request, errors.New("itemList can not be empty")
	}

	if len(itemList) > 200 {
		return request, errors.New("itemList's len should less than 200 ")
	}

	return request, nil
}
