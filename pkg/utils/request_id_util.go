package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	io "infer-microservices/pkg/services/io"
	"time"
)

func CreateRequestId(in *io.RecRequest) string {
	//Multiple requests from a user to a model within 120 seconds are considered a single request
	timestamp := time.Now().Unix() / 120
	value := in.GetDataId() + in.GetGroupId() + in.GetUserId() + fmt.Sprintf("%d", timestamp)
	data := []byte(value)
	md5New := md5.New()
	md5New.Write(data)

	requestId := hex.EncodeToString(md5New.Sum(nil))
	return requestId
}

func CreateRequestId2(dataId string, groupId string, userId string) string {
	//Multiple requests from a user to a model within 120 seconds are considered a single request
	timestamp := time.Now().Unix() / 120
	value := dataId + groupId + userId + fmt.Sprintf("%d", timestamp)
	data := []byte(value)
	md5New := md5.New()
	md5New.Write(data)

	requestId := hex.EncodeToString(md5New.Sum(nil))
	return requestId
}
