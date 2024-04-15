package service

import (
	"bytes"
	"errors"
	"fmt"
	"ms_auth/infrastructure"
	"ms_auth/model"
	"ms_auth/utils"
	"reflect"
)

type BasicQueryService interface {
	Authen(data *model.MessageSending) (*model.MessageSending, error)
}

type basicQueryService struct {
	userService  *userService
	cacheService *cacheService
}

func (s *basicQueryService) Authen(dataMessageSending *model.MessageSending) (*model.MessageSending, error) {
	var msg *model.MessageSending
	var dataBuffer *bytes.Buffer
	var dataReponse model.AuthResponse

	dataObj, err := utils.BufferToStruct(dataMessageSending.Body)
	if err != nil {
		return nil, err
	}

	reflectObj := reflect.ValueOf(dataObj)
	switch reflectObj.Type() {

	// Handle login case
	case reflect.TypeOf(model.LoginPayload{}):
		payload := reflectObj.Interface().(model.LoginPayload)
		fmt.Println("login-payload: ", payload)
		data, err := s.userService.Login(payload.Email, payload.Password)
		if err != nil {
			dataReponse = model.AuthResponse{
				ServerMessage: err.Error(),
			}
		} else {
			dataReponse = model.AuthResponse{
				UserId:        data.ID,
				Email:         data.Email,
				ServerMessage: "200 OK",
			}
		}
		dataBuffer, err = utils.StructToBuffer(dataReponse, int32(3)) // ID Mapping = 3 ~ Model.AuthResponse
		if err != nil {
			return nil, err
		}

	// Handle register case
	case reflect.TypeOf(model.RegistPayload{}):
		payload := reflectObj.Interface().(model.RegistPayload)
		fmt.Println("regist-payload: ", payload)
		data, err := s.userService.CreateUser(payload.Name, payload.Email, payload.Password)
		if err != nil {
			dataReponse = model.AuthResponse{
				ServerMessage: err.Error(),
			}
		} else {
			dataReponse = model.AuthResponse{
				UserId:        data.ID,
				Email:         data.Email,
				ServerMessage: "200 OK",
			}
		}
		dataBuffer, err = utils.StructToBuffer(dataReponse, int32(3))	// ID Mapping = 3 ~ Model.AuthResponse
		if err != nil {
			return nil, err
		}

	default:
		return nil, errors.New("=>BasicQueryService - Authen: Object model not found")
	}

	msg = &model.MessageSending{
		Salt:       dataMessageSending.Salt,
		SessionId:  dataMessageSending.SessionId,
		MessageId:  dataMessageSending.MessageId,
		SeqNo:      dataMessageSending.SeqNo,
		MessageLen: int32(len(dataBuffer.Bytes())),
		Body:       dataBuffer.Bytes(),
	}
	return msg, nil
}

func NewBasicQueryService() BasicQueryService {
	return &basicQueryService{
		userService: &userService{
			db: infrastructure.GetDB(),
		},
		cacheService: &cacheService{},
	}
}
