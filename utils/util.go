package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func FloatRound(f float32, n int) float64 {
	format := "%." + strconv.Itoa(n) + "f"
	res, _ := strconv.ParseFloat(fmt.Sprintf(format, f), 64)
	return res
}

func ConvertStructToJson(param interface{}) string {
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}

func ConvertJsonToStruct(jsonStr string) map[string]interface{} {
	var tempMap map[string]interface{}
	var err error
	if jsonStr != "" {
		err = json.Unmarshal([]byte(jsonStr), &tempMap)
		if err != nil {
			panic(err)
		}
	} else {
		tempMap = make(map[string]interface{}, 0)
	}

	return tempMap
}
