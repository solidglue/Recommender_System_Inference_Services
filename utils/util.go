package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// 截取小数位数
func FloatRound(f float32, n int) float64 {
	format := "%." + strconv.Itoa(n) + "f"
	res, _ := strconv.ParseFloat(fmt.Sprintf(format, f), 64)
	return res
}

func Struct2Json(param interface{}) string {
	//Struct转JSON
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}

func Json2Struct(str string) map[string]interface{} {
	//json转Struct
	var tempMap map[string]interface{}
	var err error

	if str != "" {
		err = json.Unmarshal([]byte(str), &tempMap)

		if err != nil {
			panic(err)
		}
	} else {

		tempMap = make(map[string]interface{}, 0)
	}

	return tempMap
}

func Json2Map(str string) map[string]interface{} {
	//json转Struct
	var tempMap map[string]interface{}
	var err error

	if str != "" {
		err = json.Unmarshal([]byte(str), &tempMap)

		if err != nil {
			panic(err)
		}
	} else {

		tempMap = make(map[string]interface{}, 0)
	}

	return tempMap
}
