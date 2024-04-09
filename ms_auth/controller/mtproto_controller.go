package controller

import (
	"context"
	"fmt"
	"log"
	"ms_auth/infrastructure"
	"ms_auth/model"
	"ms_auth/pb"
	"ms_auth/utils"
)

type MTProtoController interface {
	Send(ctx context.Context, in *pb.Message) (*pb.Message, error)
}

type mtProtoController struct {
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

	log.Println("=> Decrypt message - Client send: ", string(data.Body))

	// Reply message
	dataServerReply := fmt.Sprintf("Server seen: %s", string(data.Body))
	newMessage := model.MessageSending{
		Salt:       data.Salt,
		SessionId:  data.SessionId,
		MessageId:  data.MessageId + 1,
		SeqNo:      data.SeqNo + 1,
		MessageLen: int32(len([]byte(dataServerReply))),
		Body:       []byte(dataServerReply),
	}

	reply, err := utils.SerializeMarshal(newMessage)
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
	return &mtProtoController{}
}
