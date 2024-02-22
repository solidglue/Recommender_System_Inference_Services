package io

import (
	"errors"
	"infer-microservices/internal/logs"
)

type RecRequest struct {
	userId      string
	dataId      string //nacos dataid
	groupId     string //nacos groupId
	namespaceId string //nacos namespaceId

}

// userId
func (r *RecRequest) SetUserId(userId string) {
	r.userId = userId
}

func (r *RecRequest) GetUserId() string {
	return r.userId
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

// JavaClassName
func (r *RecRequest) JavaClassName() string {
	return "com.loki.www.infer.RecRequest"
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

	return true
}
