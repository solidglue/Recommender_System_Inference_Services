package io

import (
	"errors"
	"infer-microservices/pkg/logs"
	"strings"
)

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

func (r *RecRequest) JavaClassName() string {
	return "com.xxx.www.infer.RecRequest"
}

func (r *RecRequest) Check() bool {
	//check dataid
	if r.dataId == "" {
		err := errors.New("dataid can not be empty")
		logs.Error(err)
		return false
	}

	//check userid
	if r.userId == "" {
		err := errors.New("userid can not be empty")
		logs.Error(err)
		return false
	}

	if r.recallNum > 1000 {
		err := errors.New("recallNum should less than 2000 ")
		logs.Error(err)
		return false
	}

	//itemList
	if strings.ToLower(r.modelType) == "rank" && len(r.itemList) > 200 {
		err := errors.New("itemList's len should less than 300 ")
		logs.Error(err)
		return false
	}

	return true
}
