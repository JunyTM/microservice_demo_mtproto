package controller

import (
	"context"
	"ms_auth/infrastructure"
	"ms_auth/pb"
	"ms_auth/service"
	"ms_auth/utils"
)

type MTProtoController interface {
	Send(ctx context.Context, in *pb.Message) (*pb.Message, error)
}

type mtProtoController struct {
	basicQueryService service.BasicQueryService
	pb.UnimplementedEncryptedServiceServer
}

func (s *mtProtoController) Send(ctx context.Context, in *pb.Message) (*pb.Message, error) {
	authKey := infrastructure.GetAuthKey()

	// Caculate the AES key and IV by msgKey
	aesKey, aesIV, err := utils.ComputeAESKeyIV2([]byte(authKey), in.GetMessageKey())
	if err != nil {
		return nil, err
	}

	// Decrypt message
	plaintext, err := utils.DescriptAES_IGE(aesKey, aesIV, in.GetSerializeData())
	if err != nil {
		return nil, err
	}

	data, err := utils.SerializeUnMarshal(plaintext)
	if err != nil {
		return nil, err
	}

	// Handle the case business logic of message
	dataServerReply, err := s.basicQueryService.Authen(data)
	if err != nil {
		return nil, err
	}

	reply, err := utils.SerializeMarshal(*dataServerReply)
	if err != nil {
		return nil, err
	}

	// Caculate the new AES key and IV by Plaintext
	var newMsgKey []byte
	aesKey_, aesIV_, newMsgKey, err := utils.ComputeAESKeyIV(authKey, reply)
	if err != nil {
		return nil, err
	}

	// Encrypt message
	messageReply, err := utils.EnscriptAES_IGE(aesKey_, aesIV_, reply)
	if err != nil {
		return nil, err
	}

	return &pb.Message{
		AuthenId:      utils.GetAuthKeyId(authKey),
		MessageKey:    newMsgKey,
		SerializeData: messageReply,
	}, nil
}

func NewMTProtoController() *mtProtoController {
	return &mtProtoController{
		basicQueryService: service.NewBasicQueryService(),
	}
}
