package service

import (
	"encoding/json"
	"errors"
	"ms_gmail/model"
	"ms_gmail/utils"
	"reflect"
)

type BasicQueryService interface {
	ToBuffer(payload model.RequestPayload) ([]byte, error)
}

type basicQueryService struct{}

func (s *basicQueryService) ToBuffer(payload model.RequestPayload) ([]byte, error) {
	objectModel := model.MapMessageObject[payload.ObjectModel]
	objectModelId := model.MapMessageObjectId[payload.ObjectModel]
	if objectModelId == 0 || objectModel == nil {
		return nil, errors.New("BasicQueryService: Object model not found")
	}

	// Convert map[string]interface{} to string
	insertedData := payload.Data
	jsonString, err := json.Marshal(insertedData)
	if err != nil {
		return nil, err
	}

	//convert string to struct
	newInstanceModel := reflect.New(reflect.TypeOf(objectModel)).Interface()
	err = json.Unmarshal(jsonString, newInstanceModel)
	if err != nil && insertedData != "" {
		return nil, err
	}

	messageBuffer, err := utils.StructToBuffer(newInstanceModel, objectModelId)
	if err != nil {
		return nil, err
	}

	return messageBuffer.Bytes(), nil
}

func NewBasicQueryService() BasicQueryService {
	return &basicQueryService{}
}
